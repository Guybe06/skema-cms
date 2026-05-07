package connections

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/response"
	"skema-api/core/validator"
	"skema-api/features/connections/constants"
	"skema-api/features/connections/types"
)

// @Summary      Obtenir une connexion
// @Tags         connections
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug de l'organisation"
// @Param        id    path      string  true  "ID de la connexion"
// @Success      200   {object}  response.Body{data=types.ConnectionResponse}
// @Failure      401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/connections/{id} [get]
func (h *Handler) get(c *gin.Context) {
	conn, err := h.svc.Get(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, constants.MsgConnectionFound, toResponse(conn))
}

// @Summary      Mettre à jour une connexion
// @Tags         connections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        slug  path      string                          true  "Slug de l'organisation"
// @Param        id    path      string                          true  "ID de la connexion"
// @Param        body  body      types.UpdateConnectionRequest   true  "Champs à mettre à jour"
// @Success      200   {object}  response.Body{data=types.ConnectionResponse}
// @Failure      400,401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/connections/{id} [patch]
func (h *Handler) update(c *gin.Context) {
	var req types.UpdateConnectionRequest
	if errs := validator.BindAndValidate(c, &req); errs != nil {
		response.ValidationError(c, response.MsgValidationError, errs)
		return
	}
	conn, err := h.svc.Update(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id"), req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, constants.MsgConnectionUpdated, toResponse(conn))
}

// @Summary      Supprimer une connexion
// @Tags         connections
// @Security     BearerAuth
// @Param        slug  path  string  true  "Slug de l'organisation"
// @Param        id    path  string  true  "ID de la connexion"
// @Success      204
// @Failure      401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/connections/{id} [delete]
func (h *Handler) delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.NoContent(c)
}

// @Summary      Tester une connexion
// @Tags         connections
// @Security     BearerAuth
// @Produce      json
// @Param        slug  path      string  true  "Slug de l'organisation"
// @Param        id    path      string  true  "ID de la connexion"
// @Success      200   {object}  response.Body
// @Failure      400,401,403,404  {object}  response.Body
// @Router       /organizations/{slug}/connections/{id}/test [post]
func (h *Handler) test(c *gin.Context) {
	if err := h.svc.TestConnection(c.Request.Context(), mwauth.GetUserID(c), c.Param("slug"), c.Param("id")); err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, constants.MsgConnectionTested, nil)
}
