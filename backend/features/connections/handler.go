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

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func toResponse(c *types.Connection) types.ConnectionResponse {
	return types.ConnectionResponse{
		ID:             c.ID,
		OrganizationID: c.OrganizationID,
		Name:           c.Name,
		Driver:         c.Driver,
		Host:           c.Host,
		Port:           c.Port,
		Database:       c.Database,
		User:           c.User,
		SSLMode:        c.SSLMode,
		CreatedAt:      c.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

// @Summary      Créer une connexion
// @Description  Crée une connexion base de données dans une organisation. Réservé aux membres actifs.
// @Tags         connections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                          true  "Slug de l'organisation"
// @Param        body  body      types.CreateConnectionRequest   true  "Paramètres de connexion"
// @Success      201   {object}  response.Body{data=types.ConnectionResponse}
// @Failure      400   {object}  response.Body
// @Failure      401   {object}  response.Body
// @Failure      403   {object}  response.Body
// @Router       /organizations/{slug}/connections [post]
func (h *Handler) create(c *gin.Context) {
	var req types.CreateConnectionRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	conn, err := h.svc.Create(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), req)
	if err != nil {
		switch err.Error() {
		case constants.ErrNotAuthorized:
			response.Forbidden(c, err.Error())
		default:
			response.Internal(c, "Une erreur est survenue.")
		}
		return
	}
	response.Created(c, constants.MsgConnectionCreated, toResponse(conn))
}

// @Summary      Lister les connexions
// @Description  Retourne la liste des connexions d'une organisation.
// @Tags         connections
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug de l'organisation"
// @Success      200   {object}  response.Body{data=[]types.ConnectionResponse}
// @Failure      401   {object}  response.Body
// @Failure      403   {object}  response.Body
// @Router       /organizations/{slug}/connections [get]
func (h *Handler) list(c *gin.Context) {
	conns, err := h.svc.List(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"))
	if err != nil {
		switch err.Error() {
		case constants.ErrNotAuthorized:
			response.Forbidden(c, err.Error())
		default:
			response.Internal(c, "Une erreur est survenue.")
		}
		return
	}
	result := make([]types.ConnectionResponse, 0, len(conns))
	for _, conn := range conns {
		result = append(result, toResponse(conn))
	}
	response.OK(c, "Connexions récupérées.", result)
}

// @Summary      Obtenir une connexion
// @Description  Retourne les détails d'une connexion par son ID.
// @Tags         connections
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug de l'organisation"
// @Param        id    path      string  true  "ID de la connexion"
// @Success      200   {object}  response.Body{data=types.ConnectionResponse}
// @Failure      401   {object}  response.Body
// @Failure      403   {object}  response.Body
// @Failure      404   {object}  response.Body
// @Router       /organizations/{slug}/connections/{id} [get]
func (h *Handler) get(c *gin.Context) {
	conn, err := h.svc.Get(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"))
	if err != nil {
		switch err.Error() {
		case constants.ErrNotAuthorized:
			response.Forbidden(c, err.Error())
		case constants.ErrConnectionNotFound:
			response.NotFound(c, err.Error())
		default:
			response.Internal(c, "Une erreur est survenue.")
		}
		return
	}
	response.OK(c, "Connexion récupérée.", toResponse(conn))
}

// @Summary      Mettre à jour une connexion
// @Description  Met à jour les paramètres d'une connexion.
// @Tags         connections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                          true  "Slug de l'organisation"
// @Param        id    path      string                          true  "ID de la connexion"
// @Param        body  body      types.UpdateConnectionRequest   true  "Champs à mettre à jour"
// @Success      200   {object}  response.Body{data=types.ConnectionResponse}
// @Failure      400   {object}  response.Body
// @Failure      401   {object}  response.Body
// @Failure      403   {object}  response.Body
// @Failure      404   {object}  response.Body
// @Router       /organizations/{slug}/connections/{id} [patch]
func (h *Handler) update(c *gin.Context) {
	var req types.UpdateConnectionRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	conn, err := h.svc.Update(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"), req)
	if err != nil {
		switch err.Error() {
		case constants.ErrNotAuthorized:
			response.Forbidden(c, err.Error())
		case constants.ErrConnectionNotFound:
			response.NotFound(c, err.Error())
		default:
			response.Internal(c, "Une erreur est survenue.")
		}
		return
	}
	response.OK(c, constants.MsgConnectionUpdated, toResponse(conn))
}

// @Summary      Supprimer une connexion
// @Description  Supprime définitivement une connexion.
// @Tags         connections
// @Security     BearerAuth
// @Param        slug  path  string  true  "Slug de l'organisation"
// @Param        id    path  string  true  "ID de la connexion"
// @Success      204
// @Failure      401  {object}  response.Body
// @Failure      403  {object}  response.Body
// @Failure      404  {object}  response.Body
// @Router       /organizations/{slug}/connections/{id} [delete]
func (h *Handler) delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id")); err != nil {
		switch err.Error() {
		case constants.ErrNotAuthorized:
			response.Forbidden(c, err.Error())
		case constants.ErrConnectionNotFound:
			response.NotFound(c, err.Error())
		default:
			response.Internal(c, "Une erreur est survenue.")
		}
		return
	}
	response.NoContent(c)
}

// @Summary      Tester une connexion
// @Description  Tente d'établir une connexion à la base distante pour vérifier les credentials.
// @Tags         connections
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug de l'organisation"
// @Param        id    path      string  true  "ID de la connexion"
// @Success      200   {object}  response.Body
// @Failure      400   {object}  response.Body  "Connexion échouée"
// @Failure      401   {object}  response.Body
// @Failure      403   {object}  response.Body
// @Failure      404   {object}  response.Body
// @Router       /organizations/{slug}/connections/{id}/test [post]
func (h *Handler) test(c *gin.Context) {
	if err := h.svc.TestConnection(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id")); err != nil {
		switch err.Error() {
		case constants.ErrNotAuthorized:
			response.Forbidden(c, err.Error())
		case constants.ErrConnectionNotFound:
			response.NotFound(c, err.Error())
		case constants.ErrConnectionFailed:
			response.BadRequest(c, err.Error())
		default:
			response.Internal(c, "Une erreur est survenue.")
		}
		return
	}
	response.OK(c, constants.MsgConnectionTested, nil)
}
