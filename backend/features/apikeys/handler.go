package apikeys

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/apikeys/constants"
	"skema-api/features/apikeys/service"
	"skema-api/features/apikeys/types"
)

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

func toResp(k *types.APIKey) types.APIKeyResponse {
	return types.APIKeyResponse{
		ID:                 k.ID,
		OrganizationID:     k.OrganizationID,
		Name:               k.Name,
		KeyPrefix:          k.KeyPrefix,
		Permissions:        k.Permissions,
		AllowedCollections: k.AllowedCollections,
		CreatedAt:          k.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

// @Summary      Générer une clé API
// @Description  Génère une nouvelle clé API. La clé brute est retournée une seule fois.
// @Tags         api-keys
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                       true  "Slug organisation"
// @Param        body  body      types.CreateAPIKeyRequest    true  "Paramètres"
// @Success      201   {object}  response.Body{data=types.APIKeyCreatedResponse}
// @Failure      400,401,403  {object}  response.Body
// @Router       /organizations/{slug}/apikeys [post]
func (h *Handler) generate(c *gin.Context) {
	var req types.CreateAPIKeyRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	raw, k, err := h.svc.Generate(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, constants.MsgKeyCreated, types.APIKeyCreatedResponse{
		APIKeyResponse: toResp(k),
		RawKey:         raw,
	})
}

// @Summary      Lister les clés API
// @Tags         api-keys
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug organisation"
// @Success      200   {object}  response.Body{data=[]types.APIKeyResponse}
// @Failure      401,403  {object}  response.Body
// @Router       /organizations/{slug}/apikeys [get]
func (h *Handler) list(c *gin.Context) {
	keys, err := h.svc.List(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"))
	if err != nil {
		handleErr(c, err)
		return
	}
	result := make([]types.APIKeyResponse, 0, len(keys))
	for _, k := range keys {
		result = append(result, toResp(k))
	}
	response.OK(c, "Clés API récupérées.", result)
}

// @Summary      Révoquer une clé API
// @Tags         api-keys
// @Security     BearerAuth
// @Param        slug  path  string  true  "Slug organisation"
// @Param        id    path  string  true  "ID clé"
// @Success      204
// @Failure      401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/apikeys/{id} [delete]
func (h *Handler) revoke(c *gin.Context) {
	if err := h.svc.Revoke(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.NoContent(c)
}

func handleErr(c *gin.Context, err error) {
	switch err.Error() {
	case constants.ErrNotAuthorized:
		response.Forbidden(c, err.Error())
	case constants.ErrKeyNotFound:
		response.NotFound(c, err.Error())
	default:
		response.Internal(c, "Une erreur est survenue.")
	}
}
