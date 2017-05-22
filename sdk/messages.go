package sdk

import (
	"fmt"

	"golang.org/x/text/language"
)

type (
	lang language.Tag
	trad map[lang]string
)

//Supported API language
var (
	FR = lang(language.French)
	EN = lang(language.AmericanEnglish)
)

//Message list
var (
	MsgAppCreated                          = &Message{"MsgAppCreated", trad{FR: "L'application %s a été créée avec succès", EN: "Application %s successfully created"}, nil}
	MsgPipelineCreated                     = &Message{"MsgPipelineCreated", trad{FR: "Le pipeline %s a été créé avec succès", EN: "Pipeline %s successfully created"}, nil}
	MsgPipelineCreationAborted             = &Message{"MsgPipelineCreationAborted", trad{FR: "La création du pipeline %s a été abandonnée", EN: "Pipeline %s creation aborted"}, nil}
	MsgPipelineExists                      = &Message{"MsgPipelineExists", trad{FR: "Le pipeline %s existe déjà", EN: "Pipeline %s already exist"}, nil}
	MsgPipelineAttached                    = &Message{"MsgPipelineAttached", trad{FR: "Le pipeline %s a été attaché à l'application %s", EN: "Pipeline %s has been attached to application %s"}, nil}
	MsgPipelineTriggerCreated              = &Message{"MsgPipelineTriggerCreated", trad{FR: "Le trigger du pipeline %s de l'application %s vers le pipeline %s l'application %s a été créé avec succès", EN: "Trigger from pipeline %s of application %s to pipeline %s attached to application %s successfully created"}, nil}
	MsgAppGroupInheritPermission           = &Message{"MsgAppGroupInheritPermission", trad{FR: "Les permissions du projet sont appliquées sur l'application %s", EN: "Application %s inherits project permissions"}, nil}
	MsgAppGroupSetPermission               = &Message{"MsgAppGroupSetPermission", trad{FR: "Permission accordée au groupe %s sur l'application %s", EN: "Permission applied to group %s to application %s"}, nil}
	MsgAppVariablesCreated                 = &Message{"MsgAppVariablesCreated", trad{FR: "Les variables ont été ajoutées avec succès sur l'application %s", EN: "Application variable for %s are successfully created"}, nil}
	MsgHookCreated                         = &Message{"MsgHookCreated", trad{FR: "Hook créé sur le depôt %s vers le pipeline %s", EN: "Hook created on repository %s to pipeline %s"}, nil}
	MsgEnvironmentExists                   = &Message{"MsgEnvironmentExists", trad{FR: "L'environnement %s existe déjà", EN: "Environment %s already exist"}, nil}
	MsgEnvironmentCreated                  = &Message{"MsgEnvironmentCreated", trad{FR: "L'environnement %s a été créé avec succès", EN: "Environment %s successfully created"}, nil}
	MsgEnvironmentVariableUpdated          = &Message{"MsgEnvironmentVariableUpdated", trad{FR: "La variable %s de l'environnement %s a été mise à jour", EN: "Variable %s on environment %s has been updated"}, nil}
	MsgEnvironmentVariableCannotBeUpdated  = &Message{"MsgEnvironmentVariableCannotBeUpdated", trad{FR: "La variable %s de l'environnement %s n'a pu être mise à jour : %s", EN: "Variable %s on environment %s cannot be updated: %s"}, nil}
	MsgEnvironmentVariableCreated          = &Message{"MsgEnvironmentVariableCreated", trad{FR: "La variable %s de l'environnement %s a été ajoutée", EN: "Variable %s on environment %s has been added"}, nil}
	MsgEnvironmentVariableCannotBeCreated  = &Message{"MsgEnvironmentVariableCannotBeCreated", trad{FR: "La variable %s de l'environnement %s n'a pu être ajoutée : %s", EN: "Variable %s on environment %s cannot be added: %s"}, nil}
	MsgEnvironmentGroupUpdated             = &Message{"MsgEnvironmentGroupUpdated", trad{FR: "Le groupe %s de l'environnement %s a été mis à jour", EN: "Group %s on environment %s has been updated"}, nil}
	MsgEnvironmentGroupCannotBeUpdated     = &Message{"MsgEnvironmentGroupCannotBeUpdated", trad{FR: "Le groupe %s de l'environnement %s n'a pu être mis à jour : %s", EN: "Group %s on environment %s cannot be updated: %s"}, nil}
	MsgEnvironmentGroupCreated             = &Message{"MsgEnvironmentGroupCreated", trad{FR: "Le groupe %s de l'environnement %s a été ajouté", EN: "Group %s on environment %s has been added"}, nil}
	MsgEnvironmentGroupCannotBeCreated     = &Message{"MsgEnvironmentGroupCannotBeCreated", trad{FR: "Le groupe %s de l'environnement %s n'a pu être ajouté : %s", EN: "Group %s on environment %s cannot be added: %s"}, nil}
	MsgJobNotValidActionNotFound           = &Message{"MsgJobNotValidActionNotFound", trad{FR: "Erreur de validation du Job %s : L'action %s à l'étape %d n'a pas été trouvée", EN: "Job %s validation Failure: Unknown action %s on step #%d"}, nil}
	MsgJobNotValidInvalidActionParameter   = &Message{"MsgJobNotValidInvalidActionParameter", trad{FR: "Erreur de validation du Job %s : Le paramètre %s de l'étape %d - %s est invalide", EN: "Job %s validation Failure: Invalid parameter %s on step #%d %s"}, nil}
	MsgPipelineGroupUpdated                = &Message{"MsgPipelineGroupUpdated", trad{FR: "Les permissions du groupe %s sur le pipeline %s on été mises à jour", EN: "Permission for group %s on pipeline %s has been updated"}, nil}
	MsgPipelineGroupAdded                  = &Message{"MsgPipelineGroupAdded", trad{FR: "Les permissions du groupe %s sur le pipeline %s on été ajoutées", EN: "Permission for group %s on pipeline %s has been added"}, nil}
	MsgPipelineGroupDeleted                = &Message{"MsgPipelineGroupDeleted", trad{FR: "Les permissions du groupe %s sur le pipeline %s on été supprimées", EN: "Permission for group %s on pipeline %s has been deleted"}, nil}
	MsgPipelineStageUpdated                = &Message{"MsgPipelineStageUpdated", trad{FR: "Le stage %s a été mis à jour", EN: "Stage %s updated"}, nil}
	MsgPipelineStageAdded                  = &Message{"MsgPipelineStageAdded", trad{FR: "Le stage %s a été ajouté", EN: "Stage %s added"}, nil}
	MsgPipelineStageDeleted                = &Message{"MsgPipelineStageDeleted", trad{FR: "Le stage %s a été supprimé", EN: "Stage %s deleted"}, nil}
	MsgPipelineJobUpdated                  = &Message{"MsgPipelineJobUpdated", trad{FR: "Le job %s du stage %s a été mis à jour", EN: "Job %s in stage %s updated"}, nil}
	MsgPipelineJobAdded                    = &Message{"MsgPipelineJobAdded", trad{FR: "Le job %s du stage %s a été ajouté", EN: "Job %s in stage %s added"}, nil}
	MsgPipelineJobDeleted                  = &Message{"MsgPipelineJobDeleted", trad{FR: "Le job %s du stage %s a été supprimé", EN: "Job %s in stage %s deleted"}, nil}
	MsgSpawnInfoHatcheryStarts             = &Message{"MsgSpawnInfoHatcheryStarts", trad{FR: "La Hatchery %s (%s) a démarré le lancement du worker avec le model %s", EN: "Hatchery %s (%s) starts spawn worker with model %s"}, nil}
	MsgSpawnInfoHatcheryErrorSpawn         = &Message{"MsgSpawnInfoHatcheryErrorSpawn", trad{FR: "Une erreur est survenue lorsque la Hatchery %s (%s) a démarré un worker avec le model %s après %s, err:%s", EN: "Error while Hatchery %s (%s) spawn worker with model %s after %s, err:%s"}, nil}
	MsgSpawnInfoHatcheryStartsSuccessfully = &Message{"MsgSpawnInfoHatcheryStartsSuccessfully", trad{FR: "La Hatchery %s (%s) a démarré le worker %s avec succès en %s", EN: "Hatchery %s (%s) spawn worker %s successfully in %s"}, nil}
	MsgSpawnInfoWorkerEnd                  = &Message{"MsgSpawnInfoWorkerEnd", trad{FR: "Le worker %s a terminé et a passé %s à travailler sur les étapes", EN: "Worker %s finished working on this job and took %s to work on the steps"}, nil}
	MsgSpawnInfoJobTaken                   = &Message{"MsgSpawnInfoJobTaken", trad{FR: "Le job a été pris par le worker %s", EN: "Job was taken by worker %s"}, nil}
	MsgSpawnInfoWorkerForJob               = &Message{"MsgSpawnInfoWorkerForJob", trad{FR: "Ce worker %s a été créé pour lancer ce job", EN: "This worker %s was created to take this action"}, nil}
)

// Messages contains all sdk Messages
var Messages = map[string]*Message{
	MsgAppCreated.ID:                          MsgAppCreated,
	MsgPipelineCreated.ID:                     MsgPipelineCreated,
	MsgPipelineCreationAborted.ID:             MsgPipelineCreationAborted,
	MsgPipelineExists.ID:                      MsgPipelineExists,
	MsgPipelineAttached.ID:                    MsgPipelineAttached,
	MsgPipelineTriggerCreated.ID:              MsgPipelineTriggerCreated,
	MsgAppGroupInheritPermission.ID:           MsgAppGroupInheritPermission,
	MsgAppGroupSetPermission.ID:               MsgAppGroupSetPermission,
	MsgAppVariablesCreated.ID:                 MsgAppVariablesCreated,
	MsgHookCreated.ID:                         MsgHookCreated,
	MsgEnvironmentExists.ID:                   MsgEnvironmentExists,
	MsgEnvironmentCreated.ID:                  MsgEnvironmentCreated,
	MsgEnvironmentVariableUpdated.ID:          MsgEnvironmentVariableUpdated,
	MsgEnvironmentVariableCannotBeUpdated.ID:  MsgEnvironmentVariableCannotBeUpdated,
	MsgEnvironmentVariableCreated.ID:          MsgEnvironmentVariableCreated,
	MsgEnvironmentVariableCannotBeCreated.ID:  MsgEnvironmentVariableCannotBeCreated,
	MsgEnvironmentGroupUpdated.ID:             MsgEnvironmentGroupUpdated,
	MsgEnvironmentGroupCannotBeUpdated.ID:     MsgEnvironmentGroupCannotBeUpdated,
	MsgEnvironmentGroupCreated.ID:             MsgEnvironmentGroupCreated,
	MsgEnvironmentGroupCannotBeCreated.ID:     MsgEnvironmentGroupCannotBeCreated,
	MsgJobNotValidActionNotFound.ID:           MsgJobNotValidActionNotFound,
	MsgJobNotValidInvalidActionParameter.ID:   MsgJobNotValidInvalidActionParameter,
	MsgPipelineGroupUpdated.ID:                MsgPipelineGroupUpdated,
	MsgPipelineGroupAdded.ID:                  MsgPipelineGroupAdded,
	MsgPipelineGroupDeleted.ID:                MsgPipelineGroupDeleted,
	MsgPipelineStageUpdated.ID:                MsgPipelineStageUpdated,
	MsgPipelineStageAdded.ID:                  MsgPipelineStageAdded,
	MsgPipelineStageDeleted.ID:                MsgPipelineStageDeleted,
	MsgPipelineJobUpdated.ID:                  MsgPipelineJobUpdated,
	MsgPipelineJobAdded.ID:                    MsgPipelineJobAdded,
	MsgPipelineJobDeleted.ID:                  MsgPipelineJobDeleted,
	MsgSpawnInfoHatcheryStarts.ID:             MsgSpawnInfoHatcheryStarts,
	MsgSpawnInfoHatcheryErrorSpawn.ID:         MsgSpawnInfoHatcheryErrorSpawn,
	MsgSpawnInfoHatcheryStartsSuccessfully.ID: MsgSpawnInfoHatcheryStartsSuccessfully,
	MsgSpawnInfoWorkerEnd.ID:                  MsgSpawnInfoWorkerEnd,
	MsgSpawnInfoJobTaken.ID:                   MsgSpawnInfoJobTaken,
	MsgSpawnInfoWorkerForJob.ID:               MsgSpawnInfoWorkerForJob,
}

//Message represent a struc format translated messages
type Message struct {
	ID     string
	Format trad
	Args   []interface{}
}

//NewMessage instanciantes a new message
func NewMessage(m *Message, args ...interface{}) Message {
	return Message{
		Format: m.Format,
		Args:   args,
	}
}

// SupportedLanguages on API errors
var (
	SupportedLanguages = []language.Tag{
		language.AmericanEnglish,
		language.French,
	}
	matcher = language.NewMatcher(SupportedLanguages)
)

//String returns formated string for the specified language
func (m *Message) String(al string) string {
	acceptedLanguages, _, err := language.ParseAcceptLanguage(al)
	if err != nil {
		return fmt.Sprintf(m.Format[EN], m.Args...)
	}

	t, _, _ := matcher.Match(acceptedLanguages...)
	switch t {
	case language.French, language.AmericanEnglish:
		return fmt.Sprintf(m.Format[lang(t)], m.Args...)
	default:
		return fmt.Sprintf(m.Format[EN], m.Args...)
	}
}

// Errors implement error interface and is a set of error
type Errors []Message

func (e *Errors) Error() string {
	return e.TranslatedError(language.AmericanEnglish.String())
}

// TranslatedError returns translation for all errors
func (e *Errors) TranslatedError(al string) string {
	var s string
	for i, err := range *e {
		if i != 0 {
			s += "\n"
		}
		s += err.String(language.AmericanEnglish.String())
	}
	return s
}
