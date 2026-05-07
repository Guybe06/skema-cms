package organizations

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/organizations/constants"
	"skema-api/features/organizations/types"
)

// @Summary      Mettre à jour une organisation
// @Tags         organizations
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                           true  "Slug de l'organisation"
// @Param        body  body      types.UpdateOrganizationRequest  true  "Nouveau nom"
// @Success      200   {object}  response.Body{data=types.OrganizationResponse}
// @Failure      400,401,403,404  {object}  response.Body
// @Router       /organizations/{slug} [patch]
func (h *Handler) update(c *gin.Context) {
	var req types.UpdateOrganizationRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	org, err := h.svc.Update(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), req.Name)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, constants.MsgOrgUpdated, toResponse(org))
}

// @Summary      Supprimer une organisation
// @Tags         organizations
// @Security     BearerAuth
// @Param        slug  path  string  true  "Slug de l'organisation"
// @Success      204
// @Failure      401,403,404  {object}  response.Body
// @Router       /organizations/{slug} [delete]
func (h *Handler) delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug")); err != nil {
		handleErr(c, err)
		return
	}
	response.NoContent(c)
}

// @Summary      Transférer la propriété
// @Tags         organizations
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                           true  "Slug de l'organisation"
// @Param        body  body      types.TransferOwnershipRequest   true  "ID du nouveau propriétaire"
// @Success      200   {object}  response.Body
// @Failure      400,401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/transfer [post]
func (h *Handler) transfer(c *gin.Context) {
	var req types.TransferOwnershipRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	if err := h.svc.TransferOwnership(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), req.NewOwnerID); err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, constants.MsgOwnerTransferred, nil)
}
