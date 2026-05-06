package constants

import "time"

const (
	TokenTypeEmailVerification = "email_verification"
	TokenTypePasswordReset     = "password_reset"

	AccessTokenExpiry  = 1 * time.Hour
	RefreshTokenExpiry = 30 * 24 * time.Hour
	VerifyTokenExpiry  = 24 * time.Hour
	ResetTokenExpiry   = 1 * time.Hour

	BcryptCost = 12

	ErrEmailTaken        = "Cette adresse email est déjà utilisée."
	ErrInvalidCredentials = "Email ou mot de passe incorrect."
	ErrEmailNotVerified  = "Veuillez vérifier votre adresse email avant de vous connecter."
	ErrSessionNotFound   = "Session introuvable ou expirée."
	ErrTokenInvalid      = "Lien invalide ou expiré."
	ErrTokenAlreadySent  = "Un email a déjà été envoyé récemment."

	MsgRegistered        = "Compte créé avec succès. Vérifiez votre email."
	MsgLoggedIn          = "Connexion réussie."
	MsgLoggedOut         = "Déconnexion réussie."
	MsgTokenRefreshed    = "Token renouvelé avec succès."
	MsgEmailVerified     = "Adresse email vérifiée avec succès."
	MsgVerificationSent  = "Email de vérification envoyé."
	MsgResetSent         = "Email de réinitialisation envoyé si le compte existe."
	MsgPasswordReset     = "Mot de passe réinitialisé avec succès."

	SubjectVerification = "Vérifiez votre adresse email - Skema"
	SubjectReset        = "Réinitialisation de votre mot de passe - Skema"

	CacheKeyResend = "auth:resend:%s"
	CacheKeyReset  = "auth:reset:%s"
	CacheTTLResend = 2 * time.Minute
)
