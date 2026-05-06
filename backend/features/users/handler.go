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

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func toProfileResponse(u *types.User) types.ProfileResponse {
	return types.ProfileResponse{
		ID:            u.ID,
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

// @Summary      Obtenir le profil
// @Description  Retourne le profil de l'utilisateur connecté.
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
	response.OK(c, "Profil récupéré.", toProfileResponse(u))
}

// @Summary      Mettre à jour le profil
// @Description  Met à jour le prénom et/ou le nom de l'utilisateur connecté.
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      types.UpdateProfileRequest  true  "Champs à modifier"
// @Success      200   {object}  response.Body{data=types.ProfileResponse}
// @Failure      400   {object}  response.Body
// @Failure      401   {object}  response.Body
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
		response.Internal(c, "Une erreur est survenue.")
		return
	}
	response.OK(c, constants.MsgProfileUpdated, toProfileResponse(updated))
}

// @Summary      Changer le mot de passe
// @Description  Vérifie l'ancien mot de passe et applique le nouveau.
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      types.ChangePasswordRequest  true  "Ancien et nouveau mot de passe"
// @Success      200   {object}  response.Body
// @Failure      400   {object}  response.Body  "Mot de passe actuel incorrect"
// @Failure      401   {object}  response.Body
// @Router       /users/me/password [post]
func (h *Handler) changePassword(c *gin.Context) {
	var req types.ChangePasswordRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	if err := h.svc.ChangePassword(c.Request.Context(), mwauth.GetUserID(c), req.CurrentPassword, req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, constants.MsgPasswordChanged, nil)
}
