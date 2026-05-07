package collections

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/collections/constants"
	"skema-api/features/collections/types"
)

// @Summary      Ajouter un champ
// @Tags         collections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string               true  "Slug organisation"
// @Param        id    path      string               true  "ID collection"
// @Param        body  body      types.AddFieldRequest true  "Champ à ajouter"
// @Success      201   {object}  response.Body{data=types.FieldResponse}
// @Failure      400,401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/collections/{id}/fields [post]
func (h *Handler) addField(c *gin.Context) {
	var req types.AddFieldRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	f, err := h.svc.AddField(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"), req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, constants.MsgFieldAdded, toFieldResp(f))
}

// @Summary      Supprimer un champ
// @Tags         collections
// @Security     BearerAuth
// @Param        slug     path  string  true  "Slug organisation"
// @Param        id       path  string  true  "ID collection"
// @Param        fieldId  path  string  true  "ID champ"
// @Success      204
// @Failure      401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/collections/{id}/fields/{fieldId} [delete]
func (h *Handler) removeField(c *gin.Context) {
	if err := h.svc.RemoveField(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"), c.Param("fieldId")); err != nil {
		handleErr(c, err)
		return
	}
	response.NoContent(c)
}
