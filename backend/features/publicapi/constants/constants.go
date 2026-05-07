package constants

const (
	ErrMissingKey      = "Clé API manquante."
	ErrInvalidKey      = "Clé API invalide."
	ErrExpiredKey      = "Clé API expirée."
	ErrOrgNotFound     = "Organisation introuvable."
	ErrKeyUnauthorized = "Clé API non autorisée pour cette organisation."
	ErrCollNotFound    = "collection introuvable"
	ErrConnFailed      = "connexion impossible"
	ErrEntryNotFound   = "Entrée introuvable."
	ErrInvalidJSON     = "Corps JSON invalide."
	ErrQueryFailed     = "Erreur lors de la requête."
	ErrScanFailed      = "Erreur lors du scan."
	ErrCreateFailed    = "Erreur lors de la création."
	ErrFetchFailed     = "Erreur lors de la récupération."

	MsgPermRead   = "Permission lecture manquante."
	MsgPermCreate = "Permission création manquante."
	MsgPermUpdate = "Permission modification manquante."
	MsgPermDelete = "Permission suppression manquante."

	MsgEntriesFound = "Entrées récupérées."
	MsgEntryFound   = "Entrée récupérée."
	MsgEntryCreated = "Entrée créée."
	MsgEntryUpdated = "Entrée modifiée."
)
