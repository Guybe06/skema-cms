package collections

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/collections/constants"
	"skema-api/features/collections/types"
)

// @Summary      Modifier une collection
// @Tags         collections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                          true  "Slug organisation"
// @Param        id    path      string                          true  "ID collection"
// @Param        body  body      types.UpdateCollectionRequest   true  "Données"
// @Success      200   {object}  response.Body{data=types.CollectionResponse}
// @Failure      400,401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/collections/{id} [patch]
func (h *Handler) update(c *gin.Context) {
	var req types.UpdateCollectionRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	col, err := h.svc.Update(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"), req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, constants.MsgCollectionUpdated, toCollectionResp(col))
}

// @Summary      Supprimer une collection
// @Tags         collections
// @Security     BearerAuth
// @Param        slug  path  string  true  "Slug organisation"
// @Param        id    path  string  true  "ID collection"
// @Success      204
// @Failure      401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/collections/{id} [delete]
func (h *Handler) delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.NoContent(c)
}
