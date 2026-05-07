package connections

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/connections/constants"
	"skema-api/features/connections/service"
	"skema-api/features/connections/types"
)

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

func toResponse(c *types.Connection) types.ConnectionResponse {
	return types.ConnectionResponse{
		ID: c.ID, OrganizationID: c.OrganizationID,
		Name: c.Name, Driver: c.Driver, Host: c.Host,
		Port: c.Port, Database: c.Database, User: c.User,
		SSLMode: c.SSLMode, CreatedAt: c.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

func handleErr(c *gin.Context, err error) {
	switch err.Error() {
	case constants.ErrNotAuthorized:
		response.Forbidden(c, err.Error())
	case constants.ErrConnectionNotFound:
		response.NotFound(c, err.Error())
	case constants.ErrConnectionFailed:
		response.BadRequest(c, err.Error())
	default:
		response.Internal(c, constants.MsgInternalError)
	}
}

// @Summary      Créer une connexion
// @Tags         connections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                          true  "Slug de l'organisation"
// @Param        body  body      types.CreateConnectionRequest   true  "Paramètres de connexion"
// @Success      201   {object}  response.Body{data=types.ConnectionResponse}
// @Failure      400,401,403  {object}  response.Body
// @Router       /organizations/{slug}/connections [post]
func (h *Handler) create(c *gin.Context) {
	var req types.CreateConnectionRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	conn, err := h.svc.Create(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, constants.MsgConnectionCreated, toResponse(conn))
}

// @Summary      Lister les connexions
// @Tags         connections
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug de l'organisation"
// @Success      200   {object}  response.Body{data=[]types.ConnectionResponse}
// @Failure      401,403  {object}  response.Body
// @Router       /organizations/{slug}/connections [get]
func (h *Handler) list(c *gin.Context) {
	conns, err := h.svc.List(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"))
	if err != nil {
		handleErr(c, err)
		return
	}
	result := make([]types.ConnectionResponse, 0, len(conns))
	for _, conn := range conns {
		result = append(result, toResponse(conn))
	}
	response.OK(c, constants.MsgConnectionsFound, result)
}
