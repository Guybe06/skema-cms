package cors

import (
	"strings"

	"github.com/gin-gonic/gin"
)

/*
 * Middleware configure les politiques Cross-Origin Resource Sharing.
 *
 * Attend  : les origines autorisées (virgule-séparées) et l'environnement courant.
 * Retourne: un HandlerFunc Gin qui gère les headers CORS et les requêtes OPTIONS.
 */

func Middleware(allowedOrigins, env string) gin.HandlerFunc {
	origins := parseOrigins(allowedOrigins, env)

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if isAllowed(origin, origins) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Vary", "Origin")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func parseOrigins(raw, env string) []string {
	if env != "production" {
		return []string{"*"}
	}
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func isAllowed(origin string, allowed []string) bool {
	for _, o := range allowed {
		if o == "*" || o == origin {
			return true
		}
	}
	return false
}
