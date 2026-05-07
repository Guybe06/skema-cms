package accounts

import (
	"github.com/gin-gonic/gin"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/accounts/constants"
	"skema-api/features/accounts/types"
)

// @Summary      Demander une réinitialisation de mot de passe
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.RequestResetRequest  true  "Adresse email"
// @Success      200   {object}  response.Body
// @Failure      400,429  {object}  response.Body
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
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.ConfirmResetRequest  true  "Token et nouveau mot de passe"
// @Success      200   {object}  response.Body
// @Failure      400  {object}  response.Body
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
