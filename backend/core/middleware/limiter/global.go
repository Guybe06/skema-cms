package limiter

import (
	"github.com/gin-gonic/gin"
	"skema-api/core/response"
)

/*
 * Global applique une limite de 100 requêtes/minute par IP.
 *
 * Attend  : aucun paramètre, s'applique via router.Use().
 * Retourne: un HandlerFunc Gin qui rejette les requêtes excessives avec 429.
 */

func Global() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !globalStore.get(c.ClientIP()).Allow() {
			response.TooManyRequests(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
