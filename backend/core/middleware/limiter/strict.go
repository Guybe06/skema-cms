package limiter

import (
	"github.com/gin-gonic/gin"
	"skema-api/core/response"
)

/*
 * Strict applique une limite de 10 requêtes/minute par IP.
 *
 * Attend  : aucun paramètre, à appliquer sur les routes d'authentification.
 * Retourne: un HandlerFunc Gin qui rejette les requêtes excessives avec 429.
 */

func Strict() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !authStore.get(c.ClientIP()).Allow() {
			response.TooManyRequests(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
