package organizations

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/features/organizations/service"
)

func RegisterRoutes(api *gin.RouterGroup, svc *service.Service, jwtSecret string) {
	h := NewHandler(svc)
	guard := mwauth.Middleware(jwtSecret)

	orgs := api.Group("/organizations", guard)
	{
		orgs.POST("", h.create)
		orgs.GET("", h.list)
		orgs.GET("/:slug", h.get)
		orgs.PATCH("/:slug", h.update)
		orgs.DELETE("/:slug", h.delete)
		orgs.POST("/:slug/transfer", h.transfer)
	}
}
