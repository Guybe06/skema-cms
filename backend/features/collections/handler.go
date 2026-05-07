package collections

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/collections/constants"
	"skema-api/features/collections/service"
	"skema-api/features/collections/types"
)

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

// @Summary      Créer une collection
// @Tags         collections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                          true  "Slug organisation"
// @Param        body  body      types.CreateCollectionRequest   true  "Données collection"
// @Success      201   {object}  response.Body{data=types.CollectionResponse}
// @Failure      400,401,403  {object}  response.Body
// @Router       /organizations/{slug}/collections [post]
func (h *Handler) create(c *gin.Context) {
	var req types.CreateCollectionRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	col, err := h.svc.Create(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, constants.MsgCollectionCreated, toCollectionResp(col))
}

// @Summary      Lister les collections
// @Tags         collections
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug organisation"
// @Success      200   {object}  response.Body{data=[]types.CollectionResponse}
// @Failure      401,403  {object}  response.Body
// @Router       /organizations/{slug}/collections [get]
func (h *Handler) list(c *gin.Context) {
	cols, err := h.svc.List(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"))
	if err != nil {
		handleErr(c, err)
		return
	}
	result := make([]types.CollectionResponse, 0, len(cols))
	for _, col := range cols {
		result = append(result, toCollectionResp(col))
	}
	response.OK(c, constants.MsgCollectionsFound, result)
}

// @Summary      Obtenir une collection (avec ses champs)
// @Tags         collections
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug organisation"
// @Param        id    path      string  true  "ID collection"
// @Success      200   {object}  response.Body{data=types.CollectionResponse}
// @Failure      401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/collections/{id} [get]
func (h *Handler) get(c *gin.Context) {
	col, err := h.svc.Get(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, constants.MsgCollectionFound, toCollectionResp(col))
}
