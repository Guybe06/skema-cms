package auth

import (
	"github.com/gin-gonic/gin"
	"skema-api/core/middleware/auth"
	"skema-api/core/middleware/limiter"
)

/*
 * RegisterRoutes enregistre toutes les routes du module d'authentification.
 *
 * Attend  : le groupe de routes /v1, le service auth et le secret JWT.
 * Retourne: rien.
 */

func RegisterRoutes(api *gin.RouterGroup, svc *Service, jwtSecret string) {
	h := NewHandler(svc)
	authMiddleware := auth.Middleware(jwtSecret)

	g := api.Group("/auth")
	{
		g.POST("/register", limiter.Strict(), h.register)
		g.POST("/login", limiter.Strict(), h.login)
		g.POST("/refresh", h.refresh)
		g.POST("/verify-email", h.verifyEmail)
		g.POST("/request-reset", limiter.Strict(), h.requestReset)
		g.POST("/confirm-reset", h.confirmReset)

		protected := g.Group("", authMiddleware)
		{
			protected.POST("/logout", h.logout)
			protected.POST("/resend-verification", limiter.Strict(), h.resendVerification)
		}
	}
}
