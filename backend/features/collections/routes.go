package collections

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/features/collections/service"
)

func RegisterRoutes(api *gin.RouterGroup, svc *service.Service, jwtSecret string) {
	h := NewHandler(svc)
	guard := mwauth.Middleware(jwtSecret)

	cols := api.Group("/organizations/:slug/collections", guard)
	{
		cols.POST("", h.create)
		cols.GET("", h.list)
		cols.GET("/:id", h.get)
		cols.PATCH("/:id", h.update)
		cols.DELETE("/:id", h.delete)
		cols.POST("/:id/fields", h.addField)
		cols.DELETE("/:id/fields/:fieldId", h.removeField)
	}
}
