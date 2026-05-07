package accounts

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/accounts/constants"
	"skema-api/features/accounts/types"
)

// @Summary      Vérifier l'adresse email
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      types.VerifyEmailRequest  true  "Token de vérification"
// @Success      200   {object}  response.Body
// @Failure      400  {object}  response.Body
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
// @Tags         accounts
// @Security     BearerAuth
// @Success      200   {object}  response.Body
// @Failure      400,401,429  {object}  response.Body
// @Router       /accounts/verify/resend [post]
func (h *Handler) resendVerification(c *gin.Context) {
	if err := h.svc.ResendVerification(c.Request.Context(), mwauth.GetUserID(c)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, constants.MsgVerificationSent, nil)
}
