package organizations

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/organizations/constants"
	"skema-api/features/organizations/service"
	"skema-api/features/organizations/types"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func toResponse(org *types.Organization) types.OrganizationResponse {
	return types.OrganizationResponse{
		ID:        org.ID,
		Name:      org.Name,
		Slug:      org.Slug,
		OwnerID:   org.OwnerID,
		CreatedAt: org.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

// @Summary      Créer une organisation
// @Description  Crée une nouvelle organisation dont l'utilisateur connecté devient propriétaire. Le slug est généré automatiquement depuis le nom.
// @Tags         organizations
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      types.CreateOrganizationRequest  true  "Nom de l'organisation"
// @Success      201   {object}  response.Body{data=types.OrganizationResponse}
// @Failure      400   {object}  response.Body
// @Failure      401   {object}  response.Body
// @Router       /organizations [post]
func (h *Handler) create(c *gin.Context) {
	var req types.CreateOrganizationRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	org, err := h.svc.Create(c.Request.Context(), mwauth.GetUserID(c), req.Name)
	if err != nil {
		response.Internal(c, "Une erreur est survenue.")
		return
	}
	response.Created(c, constants.MsgOrgCreated, toResponse(org))
}

// @Summary      Lister mes organisations
// @Description  Retourne la liste des organisations dont l'utilisateur connecté est propriétaire.
// @Tags         organizations
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Body{data=[]types.OrganizationResponse}
// @Failure      401  {object}  response.Body
// @Router       /organizations [get]
func (h *Handler) list(c *gin.Context) {
	orgs, err := h.svc.ListByOwner(c.Request.Context(), mwauth.GetUserID(c))
	if err != nil {
		response.Internal(c, "Une erreur est survenue.")
		return
	}
	result := make([]types.OrganizationResponse, 0, len(orgs))
	for _, org := range orgs {
		result = append(result, toResponse(org))
	}
	response.OK(c, "Organisations récupérées.", result)
}

// @Summary      Obtenir une organisation
// @Description  Retourne les détails d'une organisation par son slug.
// @Tags         organizations
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug de l'organisation"  example:"acme-corp"
// @Success      200   {object}  response.Body{data=types.OrganizationResponse}
// @Failure      401   {object}  response.Body
// @Failure      404   {object}  response.Body
// @Router       /organizations/{slug} [get]
func (h *Handler) get(c *gin.Context) {
	org, err := h.svc.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		response.NotFound(c, constants.ErrOrgNotFound)
		return
	}
	response.OK(c, "Organisation récupérée.", toResponse(org))
}

// @Summary      Mettre à jour une organisation
// @Description  Met à jour le nom (et le slug) d'une organisation. Réservé au propriétaire.
// @Tags         organizations
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                           true  "Slug de l'organisation"
// @Param        body  body      types.UpdateOrganizationRequest  true  "Nouveau nom"
// @Success      200   {object}  response.Body{data=types.OrganizationResponse}
// @Failure      400   {object}  response.Body
// @Failure      401   {object}  response.Body
// @Failure      403   {object}  response.Body  "Pas propriétaire"
// @Failure      404   {object}  response.Body
// @Router       /organizations/{slug} [patch]
func (h *Handler) update(c *gin.Context) {
	var req types.UpdateOrganizationRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	org, err := h.svc.Update(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), req.Name)
	if err != nil {
		switch err.Error() {
		case constants.ErrOrgNotFound:
			response.NotFound(c, err.Error())
		case constants.ErrNotOwner:
			response.Forbidden(c, err.Error())
		default:
			response.Internal(c, "Une erreur est survenue.")
		}
		return
	}
	response.OK(c, constants.MsgOrgUpdated, toResponse(org))
}

// @Summary      Supprimer une organisation
// @Description  Supprime définitivement une organisation et toutes ses données. Réservé au propriétaire.
// @Tags         organizations
// @Security     BearerAuth
// @Param        slug  path  string  true  "Slug de l'organisation"
// @Success      204
// @Failure      401  {object}  response.Body
// @Failure      403  {object}  response.Body
// @Failure      404  {object}  response.Body
// @Router       /organizations/{slug} [delete]
func (h *Handler) delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug")); err != nil {
		switch err.Error() {
		case constants.ErrOrgNotFound:
			response.NotFound(c, err.Error())
		case constants.ErrNotOwner:
			response.Forbidden(c, err.Error())
		default:
			response.Internal(c, "Une erreur est survenue.")
		}
		return
	}
	response.NoContent(c)
}

// @Summary      Transférer la propriété
// @Description  Transfère la propriété de l'organisation à un autre membre actif. Réservé au propriétaire.
// @Tags         organizations
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                           true  "Slug de l'organisation"
// @Param        body  body      types.TransferOwnershipRequest   true  "ID du nouveau propriétaire"
// @Success      200   {object}  response.Body
// @Failure      400   {object}  response.Body  "Nouveau propriétaire non membre"
// @Failure      401   {object}  response.Body
// @Failure      403   {object}  response.Body
// @Failure      404   {object}  response.Body
// @Router       /organizations/{slug}/transfer [post]
func (h *Handler) transfer(c *gin.Context) {
	var req types.TransferOwnershipRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	if err := h.svc.TransferOwnership(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), req.NewOwnerID); err != nil {
		switch err.Error() {
		case constants.ErrOrgNotFound:
			response.NotFound(c, err.Error())
		case constants.ErrNotOwner:
			response.Forbidden(c, err.Error())
		case constants.ErrNewOwnerNotMember:
			response.BadRequest(c, err.Error())
		default:
			response.Internal(c, "Une erreur est survenue.")
		}
		return
	}
	response.OK(c, constants.MsgOwnerTransferred, nil)
}
