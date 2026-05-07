package memberships

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/memberships/constants"
	"skema-api/features/memberships/service"
	"skema-api/features/memberships/types"
)

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

func toResponse(m *types.Membership) types.MemberResponse {
	return types.MemberResponse{
		ID: m.ID, UserID: m.UserID, Email: m.Email,
		Role: m.Role, Status: m.Status, InvitedBy: m.InvitedBy,
		CreatedAt: m.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

// @Summary      Inviter un membre
// @Tags         memberships
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string               true  "Slug de l'organisation"
// @Param        body  body      types.InviteRequest  true  "Email et rôle du futur membre"
// @Success      200   {object}  response.Body
// @Failure      400,401,403  {object}  response.Body
// @Router       /organizations/{slug}/members/invite [post]
func (h *Handler) invite(c *gin.Context) {
	var req types.InviteRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	if err := h.svc.Invite(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), req.Email, req.Role); err != nil {
		switch err.Error() {
		case constants.ErrAlreadyMember:
			response.Conflict(c, err.Error())
		case constants.ErrNotAuthorized:
			response.Forbidden(c, err.Error())
		default:
			response.BadRequest(c, err.Error())
		}
		return
	}
	response.OK(c, constants.MsgInviteSent, nil)
}

// @Summary      Lister les membres
// @Tags         memberships
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug de l'organisation"
// @Success      200   {object}  response.Body{data=[]types.MemberResponse}
// @Failure      401,403  {object}  response.Body
// @Router       /organizations/{slug}/members [get]
func (h *Handler) list(c *gin.Context) {
	members, err := h.svc.ListMembers(c.Request.Context(), c.Param("slug"), mwauth.GetUserID(c))
	if err != nil {
		response.Forbidden(c, err.Error())
		return
	}
	result := make([]types.MemberResponse, 0, len(members))
	for _, m := range members {
		result = append(result, toResponse(m))
	}
	response.OK(c, constants.MsgMembersFound, result)
}

// @Summary      Accepter une invitation
// @Tags         memberships
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      types.AcceptInviteRequest  true  "Token d'invitation"
// @Success      200   {object}  response.Body
// @Failure      400,401  {object}  response.Body
// @Router       /invitations/accept [post]
func (h *Handler) acceptInvite(c *gin.Context) {
	var req types.AcceptInviteRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	if err := h.svc.AcceptInvite(c.Request.Context(), mwauth.GetUserID(c), req.Token); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, constants.MsgInviteAccepted, nil)
}
