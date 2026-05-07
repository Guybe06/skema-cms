package content

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/features/content/constants"
)

// @Summary      Modifier une entrée
// @Tags         content
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug     path  string  true  "Slug organisation"
// @Param        id       path  string  true  "ID collection"
// @Param        entryId  path  string  true  "ID entrée"
// @Param        body     body  object  true  "Champs à modifier"
// @Success      200  {object}  response.Body
// @Router       /organizations/{slug}/collections/{id}/content/{entryId} [patch]
func (h *Handler) update(c *gin.Context) {
	var data map[string]any
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, constants.MsgInvalidJSON)
		return
	}
	entry, err := h.svc.Update(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"), c.Param("entryId"), data)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, constants.MsgEntryUpdated, entry)
}

// @Summary      Supprimer une entrée
// @Tags         content
// @Security     BearerAuth
// @Param        slug     path  string  true  "Slug organisation"
// @Param        id       path  string  true  "ID collection"
// @Param        entryId  path  string  true  "ID entrée"
// @Success      204
// @Router       /organizations/{slug}/collections/{id}/content/{entryId} [delete]
func (h *Handler) delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"), c.Param("entryId")); err != nil {
		handleErr(c, err)
		return
	}
	response.NoContent(c)
}
