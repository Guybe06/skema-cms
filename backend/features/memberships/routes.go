package memberships

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/features/memberships/service"
)

func RegisterRoutes(api *gin.RouterGroup, svc *service.Service, jwtSecret string) {
	h := NewHandler(svc)
	guard := mwauth.Middleware(jwtSecret)

	orgs := api.Group("/organizations/:slug/members", guard)
	{
		orgs.POST("/invite", h.invite)
		orgs.GET("", h.list)
		orgs.PATCH("/:userID", h.updateRole)
		orgs.DELETE("/:userID", h.remove)
	}

	api.POST("/invitations/accept", guard, h.acceptInvite)
}
