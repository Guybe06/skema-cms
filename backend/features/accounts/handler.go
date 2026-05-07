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

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

// @Summary      Créer un compte
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.RegisterRequest  true  "Informations d'inscription"
// @Success      201   {object}  response.Body{data=types.TokenResponse}
// @Failure      400,409,429  {object}  response.Body
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
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.LoginRequest  true  "Identifiants"
// @Success      200   {object}  response.Body{data=types.TokenResponse}
// @Failure      400,401,429  {object}  response.Body
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
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.RefreshRequest  true  "Refresh token"
// @Success      200   {object}  response.Body{data=types.TokenResponse}
// @Failure      400,401  {object}  response.Body
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
