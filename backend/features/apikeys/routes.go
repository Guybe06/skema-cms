package apikeys

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/features/apikeys/service"
)

func RegisterRoutes(api *gin.RouterGroup, svc *service.Service, jwtSecret string) {
	h := NewHandler(svc)
	guard := mwauth.Middleware(jwtSecret)

	keys := api.Group("/organizations/:slug/apikeys", guard)
	{
		keys.POST("", h.generate)
		keys.GET("", h.list)
		keys.DELETE("/:id", h.revoke)
	}
}
