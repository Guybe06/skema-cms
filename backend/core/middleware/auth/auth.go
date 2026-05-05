package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"skema-api/core/response"
)

type Claims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

/*
 * Middleware valide le token JWT présent dans le header Authorization.
 *
 * Attend  : la clé secrète JWT pour vérifier la signature.
 * Retourne: un HandlerFunc Gin qui appelle Abort() avec 401 si le token est invalide.
 */

func Middleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(HeaderAuthorization)
		if header == "" || !strings.HasPrefix(header, HeaderBearerPrefix) {
			response.Unauthorized(c, MsgMissingToken)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, HeaderBearerPrefix)
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, MsgInvalidToken)
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeySessionID, claims.SessionID)
		c.Next()
	}
}

/*
 * GetUserID extrait l'identifiant de l'utilisateur authentifié depuis le contexte Gin.
 *
 * Attend  : un contexte Gin avec le middleware Auth déjà exécuté.
 * Retourne: l'identifiant utilisateur ou une chaîne vide.
 */

func GetUserID(c *gin.Context) string {
	id, _ := c.Get(ContextKeyUserID)
	userID, _ := id.(string)
	return userID
}
