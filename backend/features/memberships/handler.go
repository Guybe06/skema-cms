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

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func toResponse(m *types.Membership) types.MemberResponse {
	return types.MemberResponse{
		ID:        m.ID,
		UserID:    m.UserID,
		Email:     m.Email,
		Role:      m.Role,
		Status:    m.Status,
		InvitedBy: m.InvitedBy,
		CreatedAt: m.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

// @Summary      Inviter un membre
// @Description  Envoie une invitation par email à rejoindre l'organisation. Réservé au propriétaire et aux administrateurs.
// @Tags         memberships
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string               true  "Slug de l'organisation"
// @Param        body  body      types.InviteRequest  true  "Email et rôle du futur membre"
// @Success      200   {object}  response.Body
// @Failure      400   {object}  response.Body
// @Failure      401   {object}  response.Body
// @Failure      403   {object}  response.Body
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
// @Description  Retourne la liste des membres et invitations en attente de l'organisation.
// @Tags         memberships
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug de l'organisation"
// @Success      200   {object}  response.Body{data=[]types.MemberResponse}
// @Failure      401   {object}  response.Body
// @Failure      403   {object}  response.Body
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
	response.OK(c, "Membres récupérés.", result)
}

// @Summary      Changer le rôle d'un membre
// @Description  Modifie le rôle d'un membre actif. Réservé au propriétaire.
// @Tags         memberships
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug    path      string                   true  "Slug de l'organisation"
// @Param        userID  path      string                   true  "ID de l'utilisateur"
// @Param        body    body      types.UpdateRoleRequest  true  "Nouveau rôle"
// @Success      200     {object}  response.Body
// @Failure      400     {object}  response.Body
// @Failure      401     {object}  response.Body
// @Failure      403     {object}  response.Body
// @Failure      404     {object}  response.Body
// @Router       /organizations/{slug}/members/{userID} [patch]
func (h *Handler) updateRole(c *gin.Context) {
	var req types.UpdateRoleRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	if err := h.svc.UpdateRole(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("userID"), req.Role); err != nil {
		switch err.Error() {
		case constants.ErrMemberNotFound:
			response.NotFound(c, err.Error())
		case constants.ErrNotAuthorized, constants.ErrCannotChangeOwner:
			response.Forbidden(c, err.Error())
		default:
			response.BadRequest(c, err.Error())
		}
		return
	}
	response.OK(c, constants.MsgRoleUpdated, nil)
}

// @Summary      Retirer un membre
// @Description  Retire un membre de l'organisation. Le propriétaire peut retirer n'importe qui, un membre peut se retirer lui-même.
// @Tags         memberships
// @Security     BearerAuth
// @Param        slug    path  string  true  "Slug de l'organisation"
// @Param        userID  path  string  true  "ID de l'utilisateur"
// @Success      204
// @Failure      401  {object}  response.Body
// @Failure      403  {object}  response.Body
// @Failure      404  {object}  response.Body
// @Router       /organizations/{slug}/members/{userID} [delete]
func (h *Handler) remove(c *gin.Context) {
	if err := h.svc.RemoveMember(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("userID")); err != nil {
		switch err.Error() {
		case constants.ErrMemberNotFound:
			response.NotFound(c, err.Error())
		case constants.ErrNotAuthorized, constants.ErrCannotRemoveOwner:
			response.Forbidden(c, err.Error())
		default:
			response.BadRequest(c, err.Error())
		}
		return
	}
	response.NoContent(c)
}

// @Summary      Accepter une invitation
// @Description  Lie le compte connecté à l'invitation et active le membership.
// @Tags         memberships
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      types.AcceptInviteRequest  true  "Token d'invitation"
// @Success      200   {object}  response.Body
// @Failure      400   {object}  response.Body  "Token invalide ou expiré"
// @Failure      401   {object}  response.Body
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
