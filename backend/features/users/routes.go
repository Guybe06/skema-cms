package users

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/features/users/service"
)

func RegisterRoutes(api *gin.RouterGroup, svc *service.Service, jwtSecret string) {
	h := NewHandler(svc)
	guard := mwauth.Middleware(jwtSecret)

	me := api.Group("/users/me", guard)
	{
		me.GET("", h.getMe)
		me.PATCH("", h.updateMe)
		me.POST("/password", h.changePassword)
	}
}
