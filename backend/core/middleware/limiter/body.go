package limiter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"skema-api/core/response"
)

/*
 * Body rejette les requêtes dont le corps dépasse la taille maximale.
 *
 * Attend  : la taille maximale en mégaoctets (anti-DDoS par corps volumineux).
 * Retourne: un HandlerFunc Gin qui répond 400 si la limite est dépassée.
 */

func Body(maxMB int64) gin.HandlerFunc {
	maxBytes := maxMB * 1024 * 1024

	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			response.RequestTooLarge(c)
			c.Abort()
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
