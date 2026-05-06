package constants

const (
	ErrOrgNotFound      = "organisation introuvable"
	ErrSlugTaken        = "ce nom d'organisation est déjà utilisé"
	ErrNotOwner         = "seul le propriétaire peut effectuer cette action"
	ErrNewOwnerNotMember = "le nouveau propriétaire doit être membre de l'organisation"

	MsgOrgCreated     = "Organisation créée."
	MsgOrgUpdated     = "Organisation mise à jour."
	MsgOrgDeleted     = "Organisation supprimée."
	MsgOwnerTransferred = "Propriété transférée."
)
