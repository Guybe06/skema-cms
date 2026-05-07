package memberships

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/memberships/constants"
	"skema-api/features/memberships/types"
)

// @Summary      Changer le rôle d'un membre
// @Tags         memberships
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug    path      string                   true  "Slug de l'organisation"
// @Param        userID  path      string                   true  "ID de l'utilisateur"
// @Param        body    body      types.UpdateRoleRequest  true  "Nouveau rôle"
// @Success      200     {object}  response.Body
// @Failure      400,401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/members/{userID} [patch]
func (h *Handler) updateRole(c *gin.Context) {
	var req types.UpdateRoleRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	if err := h.svc.UpdateRole(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("userID"), req.Role); err != nil {
		switch err.Error() {
		case constants.ErrMemberNotFound:
			response.NotFound(c, err.Error())
		case constants.ErrNotAuthorized, constants.ErrCannotChangeOwner:
			response.Forbidden(c, err.Error())
		default:
			response.BadRequest(c, err.Error())
		}
		return
	}
	response.OK(c, constants.MsgRoleUpdated, nil)
}

// @Summary      Retirer un membre
// @Tags         memberships
// @Security     BearerAuth
// @Param        slug    path  string  true  "Slug de l'organisation"
// @Param        userID  path  string  true  "ID de l'utilisateur"
// @Success      204
// @Failure      401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/members/{userID} [delete]
func (h *Handler) remove(c *gin.Context) {
	if err := h.svc.RemoveMember(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("userID")); err != nil {
		switch err.Error() {
		case constants.ErrMemberNotFound:
			response.NotFound(c, err.Error())
		case constants.ErrNotAuthorized, constants.ErrCannotRemoveOwner:
			response.Forbidden(c, err.Error())
		default:
			response.BadRequest(c, err.Error())
		}
		return
	}
	response.NoContent(c)
}
