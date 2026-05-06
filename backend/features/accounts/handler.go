package accounts

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/accounts/constants"
	"skema-api/features/accounts/service"
	"skema-api/features/accounts/types"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// @Summary      Créer un compte
// @Description  Crée un nouveau compte et retourne les tokens d'accès. Un email de vérification est envoyé.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.RegisterRequest  true  "Informations d'inscription"
// @Success      201   {object}  response.Body{data=types.TokenResponse}
// @Failure      400   {object}  response.Body
// @Failure      409   {object}  response.Body  "Email déjà utilisé"
// @Failure      429   {object}  response.Body  "Trop de requêtes"
// @Router       /accounts/signup [post]
func (h *Handler) register(c *gin.Context) {
	var req types.RegisterRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	tokens, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		response.Conflict(c, err.Error())
		return
	}
	response.Created(c, constants.MsgRegistered, tokens)
}

// @Summary      Se connecter
// @Description  Authentifie un utilisateur et retourne les tokens d'accès.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.LoginRequest  true  "Identifiants"
// @Success      200   {object}  response.Body{data=types.TokenResponse}
// @Failure      400   {object}  response.Body
// @Failure      401   {object}  response.Body  "Identifiants invalides"
// @Failure      429   {object}  response.Body  "Trop de requêtes"
// @Router       /accounts/signin [post]
func (h *Handler) login(c *gin.Context) {
	var req types.LoginRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	tokens, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.OK(c, constants.MsgLoggedIn, tokens)
}

// @Summary      Se déconnecter
// @Description  Révoque la session active. Nécessite un token Bearer valide.
// @Tags         accounts
// @Security     BearerAuth
// @Success      204
// @Failure      401  {object}  response.Body
// @Router       /accounts/signout [post]
func (h *Handler) logout(c *gin.Context) {
	sessionID, _ := c.Get(mwauth.ContextKeySessionID)
	_ = h.svc.Logout(c.Request.Context(), sessionID.(string))
	response.NoContent(c)
}

// @Summary      Renouveler les tokens
// @Description  Échange un refresh token valide contre une nouvelle paire access/refresh. L'ancien token est révoqué.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.RefreshRequest  true  "Refresh token"
// @Success      200   {object}  response.Body{data=types.TokenResponse}
// @Failure      400   {object}  response.Body
// @Failure      401   {object}  response.Body  "Token invalide ou expiré"
// @Router       /accounts/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	var req types.RefreshRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	tokens, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.OK(c, constants.MsgTokenRefreshed, tokens)
}

// @Summary      Vérifier l'adresse email
// @Description  Valide le token envoyé par email et marque le compte comme vérifié.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.VerifyEmailRequest  true  "Token de vérification"
// @Success      200   {object}  response.Body
// @Failure      400   {object}  response.Body  "Token invalide ou expiré"
// @Router       /accounts/verify [post]
func (h *Handler) verifyEmail(c *gin.Context) {
	var req types.VerifyEmailRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	if err := h.svc.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, constants.MsgEmailVerified, nil)
}

// @Summary      Renvoyer l'email de vérification
// @Description  Renvoi un nouvel email de vérification. Limité à une demande toutes les 2 minutes.
// @Tags         accounts
// @Security     BearerAuth
// @Success      200   {object}  response.Body
// @Failure      400   {object}  response.Body  "Email déjà vérifié ou délai non écoulé"
// @Failure      401   {object}  response.Body
// @Failure      429   {object}  response.Body
// @Router       /accounts/verify/resend [post]
func (h *Handler) resendVerification(c *gin.Context) {
	if err := h.svc.ResendVerification(c.Request.Context(), mwauth.GetUserID(c)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, constants.MsgVerificationSent, nil)
}

// @Summary      Demander une réinitialisation de mot de passe
// @Description  Envoie un email de réinitialisation si l'adresse existe. Répond toujours 200 pour ne pas divulguer les emails.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.RequestResetRequest  true  "Adresse email"
// @Success      200   {object}  response.Body
// @Failure      400   {object}  response.Body
// @Failure      429   {object}  response.Body
// @Router       /accounts/password/reset [post]
func (h *Handler) requestReset(c *gin.Context) {
	var req types.RequestResetRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	_ = h.svc.RequestReset(c.Request.Context(), req.Email)
	response.OK(c, constants.MsgResetSent, nil)
}

// @Summary      Confirmer le nouveau mot de passe
// @Description  Réinitialise le mot de passe à partir du token reçu par email. Toutes les sessions actives sont révoquées.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.ConfirmResetRequest  true  "Token et nouveau mot de passe"
// @Success      200   {object}  response.Body
// @Failure      400   {object}  response.Body  "Token invalide ou expiré"
// @Router       /accounts/password/reset/confirm [post]
func (h *Handler) confirmReset(c *gin.Context) {
	var req types.ConfirmResetRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	if err := h.svc.ConfirmReset(c.Request.Context(), req.Token, req.Password); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, constants.MsgPasswordReset, nil)
}
