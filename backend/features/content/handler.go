package content

import (
	"strconv"

	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/features/content/constants"
	"skema-api/features/content/service"
)

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

func parsePagination(c *gin.Context) (page, perPage int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ = strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return
}

// @Summary      Lister les entrées d'une collection
// @Tags         content
// @Security     BearerAuth
// @Produce      json
// @Param        slug         path   string  true   "Slug organisation"
// @Param        id           path   string  true   "ID collection"
// @Param        page         query  int     false  "Page (défaut 1)"
// @Param        per_page     query  int     false  "Taille page (défaut 20)"
// @Param        sort         query  string  false  "Colonne de tri"
// @Param        order        query  string  false  "asc ou desc"
// @Success      200  {object}  response.ListBody
// @Router       /organizations/{slug}/collections/{id}/content [get]
func (h *Handler) list(c *gin.Context) {
	page, perPage := parsePagination(c)
	entries, total, err := h.svc.List(c.Request.Context(),
		mwauth.GetUserID(c), c.Param("slug"), c.Param("id"),
		page, perPage, c.Query("sort"), c.Query("order"))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.List(c, constants.MsgEntriesFound, entries, total, page, perPage)
}

// @Summary      Créer une entrée
// @Tags         content
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string          true  "Slug organisation"
// @Param        id    path      string          true  "ID collection"
// @Param        body  body      object          true  "Données de l'entrée"
// @Success      201   {object}  response.Body
// @Router       /organizations/{slug}/collections/{id}/content [post]
func (h *Handler) create(c *gin.Context) {
	var data map[string]any
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "Corps JSON invalide.")
		return
	}
	entry, err := h.svc.Create(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"), data)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, constants.MsgEntryCreated, entry)
}

// @Summary      Obtenir une entrée
// @Tags         content
// @Security     BearerAuth
// @Produce      json
// @Param        slug     path  string  true  "Slug organisation"
// @Param        id       path  string  true  "ID collection"
// @Param        entryId  path  string  true  "ID entrée"
// @Success      200  {object}  response.Body
// @Router       /organizations/{slug}/collections/{id}/content/{entryId} [get]
func (h *Handler) get(c *gin.Context) {
	entry, err := h.svc.Get(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"), c.Param("entryId"))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Entrée récupérée.", entry)
}

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
		response.BadRequest(c, "Corps JSON invalide.")
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

func handleErr(c *gin.Context, err error) {
	switch err.Error() {
	case constants.ErrNotAuthorized:
		response.Forbidden(c, err.Error())
	case constants.ErrEntryNotFound:
		response.NotFound(c, err.Error())
	default:
		response.Internal(c, "Une erreur est survenue.")
	}
}
