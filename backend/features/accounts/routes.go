package accounts

import (
	"github.com/gin-gonic/gin"
	mwauth "skema-api/core/middleware/auth"
	"skema-api/core/middleware/limiter"
	"skema-api/features/accounts/service"
)

func RegisterRoutes(api *gin.RouterGroup, svc *service.Service, jwtSecret string) {
	h := NewHandler(svc)
	guard := mwauth.Middleware(jwtSecret)

	accounts := api.Group("/accounts")
	{
		accounts.POST("/signup", limiter.Strict(), h.register)
		accounts.POST("/signin", limiter.Strict(), h.login)
		accounts.POST("/refresh", h.refresh)
		accounts.POST("/verify", h.verifyEmail)
		accounts.POST("/verify/resend", limiter.Strict(), guard, h.resendVerification)

		accounts.POST("/password/reset", limiter.Strict(), h.requestReset)
		accounts.POST("/password/reset/confirm", h.confirmReset)

		protected := accounts.Group("", guard)
		{
			protected.POST("/signout", h.logout)
		}
	}
}
