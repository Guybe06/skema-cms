package auth

const (
	HeaderAuthorization = "Authorization"
	HeaderBearerPrefix  = "Bearer "
	ContextKeyUserID    = "user_id"
	ContextKeySessionID = "session_id"

	MsgMissingToken = "Token d'authentification manquant."
	MsgInvalidToken = "Token d'authentification invalide ou expiré."
)
