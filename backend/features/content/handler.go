package content

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/features/content/constants"
	"skema-api/features/content/service"
)

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

// @Summary      Lister les entrées d'une collection
// @Tags         content
// @Security     BearerAuth
// @Produce      json
// @Param        slug      path   string  true   "Slug organisation"
// @Param        id        path   string  true   "ID collection"
// @Param        page      query  int     false  "Page (défaut 1)"
// @Param        per_page  query  int     false  "Taille page (défaut 20)"
// @Param        sort      query  string  false  "Colonne de tri"
// @Param        order     query  string  false  "asc ou desc"
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
// @Param        slug  path      string  true  "Slug organisation"
// @Param        id    path      string  true  "ID collection"
// @Param        body  body      object  true  "Données de l'entrée"
// @Success      201   {object}  response.Body
// @Router       /organizations/{slug}/collections/{id}/content [post]
func (h *Handler) create(c *gin.Context) {
	var data map[string]any
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, constants.MsgInvalidJSON)
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
	response.OK(c, constants.MsgEntryFound, entry)
}
