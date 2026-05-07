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

func toCollectionResp(c *types.Collection) types.CollectionResponse {
	r := types.CollectionResponse{
		ID: c.ID, ConnectionID: c.ConnectionID, OrganizationID: c.OrganizationID,
		Name: c.Name, TableName: c.TableName, DisplayName: c.DisplayName,
		Description: c.Description, CreatedAt: c.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
	for _, f := range c.Fields {
		r.Fields = append(r.Fields, toFieldResp(f))
	}
	return r
}

func toFieldResp(f *types.Field) *types.FieldResponse {
	return &types.FieldResponse{
		ID: f.ID, Name: f.Name, ColumnName: f.ColumnName, Type: f.Type,
		Required: f.Required, IsUnique: f.IsUnique,
		DefaultValue: f.DefaultValue, Options: f.Options, Position: f.Position,
	}
}

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
	response.OK(c, "Collections récupérées.", result)
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
	response.OK(c, "Collection récupérée.", toCollectionResp(col))
}

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

func handleErr(c *gin.Context, err error) {
	switch err.Error() {
	case constants.ErrNotAuthorized:
		response.Forbidden(c, err.Error())
	case constants.ErrCollectionNotFound, constants.ErrFieldNotFound:
		response.NotFound(c, err.Error())
	case constants.ErrTableNameTaken, constants.ErrColumnNameTaken, constants.ErrSchemaFailed:
		response.BadRequest(c, err.Error())
	default:
		response.Internal(c, "Une erreur est survenue.")
	}
}
