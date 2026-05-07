package content

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/features/content/service"
)

func RegisterRoutes(api *gin.RouterGroup, svc *service.Service, jwtSecret string) {
	h := NewHandler(svc)
	guard := mwauth.Middleware(jwtSecret)

	entries := api.Group("/organizations/:slug/collections/:id/content", guard)
	{
		entries.GET("", h.list)
		entries.POST("", h.create)
		entries.GET("/:entryId", h.get)
		entries.PATCH("/:entryId", h.update)
		entries.DELETE("/:entryId", h.delete)
	}
}
