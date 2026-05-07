package publicapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"skema-api/core/response"
	apikeytypes "skema-api/features/apikeys/types"
	pubconstants "skema-api/features/publicapi/constants"
)

const keyCtxKey = "pub_api_key"
const orgIDCtxKey = "pub_org_id"

type keySvc interface {
	FindByHash(ctx context.Context, hash string) (*apikeytypes.APIKey, error)
	TouchLastUsed(ctx context.Context, id string)
}

type orgSvc interface {
	FindBySlugID(ctx context.Context, slug string) (id string, err error)
}

/*
 * APIKeyAuth valide la clé API et charge l'organisation dans le contexte Gin.
 *
 * Attend  : un service de clés et un service d'organisations.
 * Retourne: un middleware Gin qui rejette les requêtes non autorisées.
 */

func APIKeyAuth(ks keySvc, os orgSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := extractRawKey(c)
		if raw == "" {
			response.Unauthorized(c, pubconstants.ErrMissingKey)
			c.Abort()
			return
		}

		sum := sha256.Sum256([]byte(raw))
		hash := hex.EncodeToString(sum[:])

		k, err := ks.FindByHash(c.Request.Context(), hash)
		if err != nil || k == nil {
			response.Unauthorized(c, pubconstants.ErrInvalidKey)
			c.Abort()
			return
		}
		if k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt) {
			response.Unauthorized(c, pubconstants.ErrExpiredKey)
			c.Abort()
			return
		}

		orgSlug := c.Param("orgSlug")
		orgID, err := os.FindBySlugID(c.Request.Context(), orgSlug)
		if err != nil || orgID == "" {
			response.NotFound(c, pubconstants.ErrOrgNotFound)
			c.Abort()
			return
		}
		if k.OrganizationID != orgID {
			response.Forbidden(c, pubconstants.ErrKeyUnauthorized)
			c.Abort()
			return
		}

		c.Set(keyCtxKey, k)
		c.Set(orgIDCtxKey, orgID)
		go ks.TouchLastUsed(c.Request.Context(), k.ID)
		c.Next()
	}
}

func extractRawKey(c *gin.Context) string {
	if auth := c.GetHeader("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return c.GetHeader("X-Api-Key")
}
