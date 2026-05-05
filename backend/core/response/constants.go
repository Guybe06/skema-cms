package response

const (
	StatusOK          = 200
	StatusCreated     = 201
	StatusNoContent   = 204
	StatusBadRequest  = 400
	StatusUnauthorized = 401
	StatusForbidden   = 403
	StatusNotFound    = 404
	StatusConflict        = 409
	StatusTooManyRequests = 429
	StatusInternal        = 500

	MsgSuccess         = "Opération réussie."
	MsgCreated         = "Ressource créée avec succès."
	MsgDeleted         = "Ressource supprimée avec succès."
	MsgInternalError   = "Une erreur interne est survenue."
	MsgNotFound        = "Ressource introuvable."
	MsgUnauthorized    = "Authentification requise."
	MsgForbidden       = "Accès refusé."
	MsgConflict        = "Cette ressource existe déjà."
	MsgValidationError   = "Les données envoyées sont invalides."
	MsgTooManyRequests   = "Trop de requêtes. Veuillez réessayer plus tard."
	MsgRequestTooLarge   = "La taille de la requête dépasse la limite autorisée."
	MsgRequestTimeout    = "La requête a expiré."
)
