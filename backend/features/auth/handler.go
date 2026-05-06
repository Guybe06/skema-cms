package auth

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/auth/constants"
	"skema-api/features/auth/service"
	"skema-api/features/auth/types"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// @Summary      Inscription
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body types.RegisterRequest true "Informations d'inscription"
// @Success      201 {object} response.Body
// @Router       /auth/register [post]
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

// @Summary      Connexion
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body types.LoginRequest true "Identifiants"
// @Success      200 {object} response.Body
// @Router       /auth/login [post]
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

// @Summary      Déconnexion
// @Tags         auth
// @Security     BearerAuth
// @Success      204
// @Router       /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	sessionID, _ := c.Get(mwauth.ContextKeySessionID)
	_ = h.svc.Logout(c.Request.Context(), sessionID.(string))
	response.NoContent(c)
}

// @Summary      Renouveler le token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body types.RefreshRequest true "Refresh token"
// @Success      200 {object} response.Body
// @Router       /auth/refresh [post]
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

// @Summary      Vérifier l'email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body types.VerifyEmailRequest true "Token de vérification"
// @Success      200 {object} response.Body
// @Router       /auth/verify-email [post]
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
// @Tags         auth
// @Security     BearerAuth
// @Success      200 {object} response.Body
// @Router       /auth/resend-verification [post]
func (h *Handler) resendVerification(c *gin.Context) {
	if err := h.svc.ResendVerification(c.Request.Context(), mwauth.GetUserID(c)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, constants.MsgVerificationSent, nil)
}

// @Summary      Demander une réinitialisation de mot de passe
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body types.RequestResetRequest true "Adresse email"
// @Success      200 {object} response.Body
// @Router       /auth/request-reset [post]
func (h *Handler) requestReset(c *gin.Context) {
	var req types.RequestResetRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	_ = h.svc.RequestReset(c.Request.Context(), req.Email)
	response.OK(c, constants.MsgResetSent, nil)
}

// @Summary      Confirmer la réinitialisation de mot de passe
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body types.ConfirmResetRequest true "Token et nouveau mot de passe"
// @Success      200 {object} response.Body
// @Router       /auth/confirm-reset [post]
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
