package timeout

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"skema-api/core/response"
)

/*
 * Middleware annule les requêtes qui dépassent la durée maximale autorisée.
 *
 * Attend  : la durée maximale en secondes (anti-Slowloris).
 * Retourne: un HandlerFunc Gin qui répond 408 si le délai est dépassé.
 */

func Middleware(seconds int) gin.HandlerFunc {
	d := time.Duration(seconds) * time.Second

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{}, 1)
		go func() {
			c.Next()
			done <- struct{}{}
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			response.RequestTimeout(c)
			c.Abort()
		}
	}
}
