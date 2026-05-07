package users

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/users/constants"
	"skema-api/features/users/types"
)

// @Summary      Changer le mot de passe
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      types.ChangePasswordRequest  true  "Ancien et nouveau mot de passe"
// @Success      200   {object}  response.Body
// @Failure      400,401  {object}  response.Body
// @Router       /users/me/password [post]
func (h *Handler) changePassword(c *gin.Context) {
	var req types.ChangePasswordRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	if err := h.svc.ChangePassword(c.Request.Context(), mwauth.GetUserID(c), req.CurrentPassword, req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, constants.MsgPasswordChanged, nil)
}
