package users

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/users/constants"
	"skema-api/features/users/service"
	"skema-api/features/users/types"
)

type Handler struct{ svc *service.Service }

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

func toProfileResponse(u *types.User) types.ProfileResponse {
	return types.ProfileResponse{
		ID: u.ID, Email: u.Email, FirstName: u.FirstName,
		LastName: u.LastName, EmailVerified: u.EmailVerified,
		CreatedAt: u.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

// @Summary      Obtenir le profil
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Body{data=types.ProfileResponse}
// @Failure      401  {object}  response.Body
// @Router       /users/me [get]
func (h *Handler) getMe(c *gin.Context) {
	u, err := h.svc.GetProfile(c.Request.Context(), mwauth.GetUserID(c))
	if err != nil {
		response.NotFound(c, constants.ErrUserNotFound)
		return
	}
	response.OK(c, constants.MsgProfileFound, toProfileResponse(u))
}

// @Summary      Mettre à jour le profil
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      types.UpdateProfileRequest  true  "Champs à modifier"
// @Success      200   {object}  response.Body{data=types.ProfileResponse}
// @Failure      400,401  {object}  response.Body
// @Router       /users/me [patch]
func (h *Handler) updateMe(c *gin.Context) {
	var req types.UpdateProfileRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	u, err := h.svc.GetProfile(c.Request.Context(), mwauth.GetUserID(c))
	if err != nil {
		response.NotFound(c, constants.ErrUserNotFound)
		return
	}
	firstName := u.FirstName
	lastName := u.LastName
	if req.FirstName != "" {
		firstName = req.FirstName
	}
	if req.LastName != "" {
		lastName = req.LastName
	}
	updated, err := h.svc.UpdateProfile(c.Request.Context(), u.ID, firstName, lastName)
	if err != nil {
		response.Internal(c, constants.MsgInternalError)
		return
	}
	response.OK(c, constants.MsgProfileUpdated, toProfileResponse(updated))
}
