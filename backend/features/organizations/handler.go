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

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

func toResponse(org *types.Organization) types.OrganizationResponse {
	return types.OrganizationResponse{
		ID: org.ID, Name: org.Name, Slug: org.Slug,
		OwnerID:   org.OwnerID,
		CreatedAt: org.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

func handleErr(c *gin.Context, err error) {
	switch err.Error() {
	case constants.ErrOrgNotFound:
		response.NotFound(c, err.Error())
	case constants.ErrNotOwner:
		response.Forbidden(c, err.Error())
	case constants.ErrNewOwnerNotMember:
		response.BadRequest(c, err.Error())
	default:
		response.Internal(c, constants.MsgInternalError)
	}
}

// @Summary      Créer une organisation
// @Tags         organizations
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      types.CreateOrganizationRequest  true  "Nom de l'organisation"
// @Success      201   {object}  response.Body{data=types.OrganizationResponse}
// @Failure      400,401  {object}  response.Body
// @Router       /organizations [post]
func (h *Handler) create(c *gin.Context) {
	var req types.CreateOrganizationRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	org, err := h.svc.Create(c.Request.Context(), mwauth.GetUserID(c), req.Name)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, constants.MsgOrgCreated, toResponse(org))
}

// @Summary      Lister mes organisations
// @Tags         organizations
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Body{data=[]types.OrganizationResponse}
// @Failure      401  {object}  response.Body
// @Router       /organizations [get]
func (h *Handler) list(c *gin.Context) {
	orgs, err := h.svc.ListByOwner(c.Request.Context(), mwauth.GetUserID(c))
	if err != nil {
		handleErr(c, err)
		return
	}
	result := make([]types.OrganizationResponse, 0, len(orgs))
	for _, org := range orgs {
		result = append(result, toResponse(org))
	}
	response.OK(c, constants.MsgOrgsFound, result)
}

// @Summary      Obtenir une organisation
// @Tags         organizations
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug de l'organisation"
// @Success      200   {object}  response.Body{data=types.OrganizationResponse}
// @Failure      401,404  {object}  response.Body
// @Router       /organizations/{slug} [get]
func (h *Handler) get(c *gin.Context) {
	org, err := h.svc.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		response.NotFound(c, constants.ErrOrgNotFound)
		return
	}
	response.OK(c, constants.MsgOrgFound, toResponse(org))
}
