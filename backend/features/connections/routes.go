package connections

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/features/connections/service"
)

func RegisterRoutes(api *gin.RouterGroup, svc *service.Service, jwtSecret string) {
	h := NewHandler(svc)
	guard := mwauth.Middleware(jwtSecret)

	conns := api.Group("/organizations/:slug/connections", guard)
	{
		conns.POST("", h.create)
		conns.GET("", h.list)
		conns.GET("/:id", h.get)
		conns.PATCH("/:id", h.update)
		conns.DELETE("/:id", h.delete)
		conns.POST("/:id/test", h.test)
	}
}
