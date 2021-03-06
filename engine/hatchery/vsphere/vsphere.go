package vsphere

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rockbears/log"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/ovh/cds/engine/api"
	"github.com/ovh/cds/engine/service"
	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/cdsclient"
	"github.com/ovh/cds/sdk/hatchery"
)

var (
	ipsInfos = struct {
		mu  sync.RWMutex
		ips map[string]ipInfos
	}{
		mu:  sync.RWMutex{},
		ips: map[string]ipInfos{},
	}
)

// New instanciates a new Hatchery vsphere
func New() *HatcheryVSphere {
	s := new(HatcheryVSphere)
	s.GoRoutines = sdk.NewGoRoutines()
	s.Router = &api.Router{
		Mux: mux.NewRouter(),
	}
	return s
}

var _ hatchery.InterfaceWithModels = new(HatcheryVSphere)

// Init cdsclient config.
func (h *HatcheryVSphere) Init(config interface{}) (cdsclient.ServiceConfig, error) {
	var cfg cdsclient.ServiceConfig
	sConfig, ok := config.(HatcheryConfiguration)
	if !ok {
		return cfg, sdk.WithStack(fmt.Errorf("invalid vsphere hatchery configuration"))
	}

	cfg.Host = sConfig.API.HTTP.URL
	cfg.Token = sConfig.API.Token
	cfg.InsecureSkipVerifyTLS = sConfig.API.HTTP.Insecure
	cfg.RequestSecondsTimeout = sConfig.API.RequestTimeout
	return cfg, nil
}

func (h *HatcheryVSphere) GetLogger() *logrus.Logger {
	return h.ServiceLogger
}

// ApplyConfiguration apply an object of type HatcheryConfiguration after checking it
func (h *HatcheryVSphere) ApplyConfiguration(cfg interface{}) error {
	if err := h.CheckConfiguration(cfg); err != nil {
		return err
	}

	var ok bool
	h.Config, ok = cfg.(HatcheryConfiguration)
	if !ok {
		return fmt.Errorf("Invalid configuration")
	}

	h.Common.Common.ServiceName = h.Config.Name
	h.Common.Common.ServiceType = sdk.TypeHatchery
	h.HTTPURL = h.Config.URL
	h.MaxHeartbeatFailures = h.Config.API.MaxHeartbeatFailures
	var err error
	h.Common.Common.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(h.Config.RSAPrivateKey))
	if err != nil {
		return fmt.Errorf("unable to parse RSA private Key: %v", err)
	}

	return nil
}

// Status returns sdk.MonitoringStatus, implements interface service.Service
func (h *HatcheryVSphere) Status(ctx context.Context) *sdk.MonitoringStatus {
	m := h.NewMonitoringStatus()
	m.AddLine(sdk.MonitoringStatusLine{Component: "Workers", Value: fmt.Sprintf("%d/%d", len(h.WorkersStarted(ctx)), h.Config.Provision.MaxWorker), Status: sdk.MonitoringStatusOK})
	return m
}

// CheckConfiguration checks the validity of the configuration object
func (h *HatcheryVSphere) CheckConfiguration(cfg interface{}) error {
	hconfig, ok := cfg.(HatcheryConfiguration)
	if !ok {
		return fmt.Errorf("Invalid hatchery vsphere configuration")
	}

	if err := hconfig.Check(); err != nil {
		return fmt.Errorf("Invalid hatchery vsphere configuration: %v", err)
	}

	if hconfig.VSphereUser == "" {
		return fmt.Errorf("vsphere-user is mandatory")
	}

	if hconfig.VSphereEndpoint == "" {
		return fmt.Errorf("vsphere-endpoint is mandatory")
	}

	if hconfig.VSpherePassword == "" {
		return fmt.Errorf("vsphere-password is mandatory")
	}

	if hconfig.VSphereDatacenterString == "" {
		return fmt.Errorf("vsphere-datacenter is mandatory")
	}

	if hconfig.IPRange != "" {
		ips, err := sdk.IPinRanges(context.Background(), hconfig.IPRange)
		if err != nil {
			return fmt.Errorf("flag or environment variable openstack-ip-range error: %v", err)
		}
		for _, ip := range ips {
			ipsInfos.ips[ip] = ipInfos{}
		}
	}
	return nil
}

// CanSpawn return wether or not hatchery can spawn model
// requirements are not supported
func (h *HatcheryVSphere) CanSpawn(ctx context.Context, model *sdk.Model, jobID int64, requirements []sdk.Requirement) bool {
	for _, r := range requirements {
		if r.Type == sdk.ServiceRequirement || r.Type == sdk.MemoryRequirement || r.Type == sdk.HostnameRequirement {
			return false
		}
	}
	return true
}

// Start inits client and routines for hatchery
func (h *HatcheryVSphere) Start(ctx context.Context) error {
	return hatchery.Create(ctx, h)
}

// Serve start the hatchery server
func (h *HatcheryVSphere) Serve(ctx context.Context) error {
	return h.CommonServe(ctx, h)
}

//Configuration returns Hatchery CommonConfiguration
func (h *HatcheryVSphere) Configuration() service.HatcheryCommonConfiguration {
	return h.Config.HatcheryCommonConfiguration
}

// NeedRegistration return true if worker model need regsitration
func (h *HatcheryVSphere) NeedRegistration(ctx context.Context, m *sdk.Model) bool {
	model, errG := h.getModelByName(ctx, m.Name)
	if errG != nil || model.Config == nil || model.Config.Annotation == "" {
		return true
	}

	var annot annotation
	if err := json.Unmarshal([]byte(model.Config.Annotation), &annot); err != nil {
		return true
	}

	return !annot.ToDelete && (m.NeedRegistration || fmt.Sprintf("%d", m.UserLastModified.Unix()) != annot.WorkerModelLastModified)
}

// WorkerModelsEnabled returns Worker model enabled
func (h *HatcheryVSphere) WorkerModelsEnabled() ([]sdk.Model, error) {
	return h.CDSClient().WorkerModelEnabledList()
}

// WorkerModelSecretList returns secret for given model.
func (h *HatcheryVSphere) WorkerModelSecretList(m sdk.Model) (sdk.WorkerModelSecrets, error) {
	return h.CDSClient().WorkerModelSecretList(m.Group.Name, m.Name)
}

// WorkersStartedByModel returns the number of instances of given model started but
// not necessarily register on CDS yet
func (h *HatcheryVSphere) WorkersStartedByModel(ctx context.Context, model *sdk.Model) int {
	var x int
	for _, s := range h.getServers(ctx) {
		if strings.Contains(strings.ToLower(s.Name), strings.ToLower(model.Name)) {
			x++
		}
	}
	log.Debug(ctx, "WorkersStartedByModel> %s : %d", model.Name, x)

	return x
}

// WorkersStarted returns the number of instances started but
// not necessarily register on CDS yet
func (h *HatcheryVSphere) WorkersStarted(ctx context.Context) []string {
	srvs := h.getServers(ctx)
	res := make([]string, len(srvs))
	for i, s := range srvs {
		if strings.Contains(strings.ToLower(s.Name), "worker") {
			res[i] = s.Name
		}
	}
	return res
}

// ModelType returns type of hatchery
func (*HatcheryVSphere) ModelType() string {
	return sdk.VSphere
}

func (h *HatcheryVSphere) main(ctx context.Context) {
	serverListTick := time.NewTicker(10 * time.Second).C
	cdnConfTick := time.NewTicker(10 * time.Second).C
	killAwolServersTick := time.NewTicker(20 * time.Second).C
	killDisabledWorkersTick := time.NewTicker(60 * time.Second).C

	for {
		select {
		case <-serverListTick:
			h.updateServerList(ctx)
		case <-killAwolServersTick:
			h.killAwolServers(ctx)
		case <-killDisabledWorkersTick:
			h.killDisabledWorkers(ctx)
		case <-cdnConfTick:
			if err := h.RefreshServiceLogger(ctx); err != nil {
				log.Error(ctx, "Hatchery> vsphere> Cannot get cdn configuration : %v", err)
			}
		case <-ctx.Done():
			if ctx.Err() != nil {
				log.Error(ctx, "Hatchery> vsphereq> Exiting routines")
			}
			return
		}
	}
}

func (h *HatcheryVSphere) updateServerList(ctx context.Context) {
	var out string
	var total int
	status := map[string]int{}

	for _, s := range h.getServers(ctx) {
		out += fmt.Sprintf("- [%s] %s ", s.Summary.Config.Name, s.Summary.Runtime.PowerState)
		if _, ok := status[string(s.Summary.Runtime.PowerState)]; !ok {
			status[string(s.Summary.Runtime.PowerState)] = 0
		}
		status[string(s.Summary.Runtime.PowerState)]++
		total++
	}
	var st string
	for k, s := range status {
		st += fmt.Sprintf("%d %s ", s, k)
	}
	log.Debug(context.TODO(), "Got %d servers %s", total, st)
	if total > 0 {
		log.Debug(context.TODO(), out)
	}
}

// killDisabledWorkers kill workers which are disabled
func (h *HatcheryVSphere) killDisabledWorkers(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	workerPoolDisabled, err := hatchery.WorkerPool(ctx, h, sdk.StatusDisabled)
	if err != nil {
		log.Error(ctx, "killDisabledWorkers> Pool> Error: %v", err)
		return
	}

	srvs := h.getServers(ctx)
	for _, w := range workerPoolDisabled {
		for _, s := range srvs {
			if s.Name == w.Name {
				log.Info(ctx, "killDisabledWorkers> killDisabledWorkers %v", s.Name)
				_ = h.deleteServer(s)
				break
			}
		}
	}
}

// killAwolServers kill unused servers
func (h *HatcheryVSphere) killAwolServers(ctx context.Context) {
	srvs := h.getServers(ctx)

	for _, s := range srvs {
		var annot annotation
		if s.Config == nil || s.Config.Annotation == "" {
			continue
		}
		if err := json.Unmarshal([]byte(s.Config.Annotation), &annot); err != nil {
			continue
		}

		if annot.ToDelete || (s.Summary.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOn && (!annot.Model || annot.RegisterOnly)) {
			if err := h.deleteServer(s); err != nil {
				log.Warn(context.Background(), "killAwolServers> cannot delete server %s", s.Name)
			}
		}
	}
}
