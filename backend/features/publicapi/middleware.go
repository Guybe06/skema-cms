package publicapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"skema-api/core/response"
	apikeytypes "skema-api/features/apikeys/types"
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

// APIKeyAuth valide la clé API et charge l'organisation dans le contexte.
func APIKeyAuth(ks keySvc, os orgSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := extractRawKey(c)
		if raw == "" {
			response.Unauthorized(c, "Clé API manquante.")
			c.Abort()
			return
		}

		sum := sha256.Sum256([]byte(raw))
		hash := hex.EncodeToString(sum[:])

		k, err := ks.FindByHash(c.Request.Context(), hash)
		if err != nil || k == nil {
			response.Unauthorized(c, "Clé API invalide.")
			c.Abort()
			return
		}
		if k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt) {
			response.Unauthorized(c, "Clé API expirée.")
			c.Abort()
			return
		}

		orgSlug := c.Param("orgSlug")
		orgID, err := os.FindBySlugID(c.Request.Context(), orgSlug)
		if err != nil || orgID == "" {
			response.NotFound(c, "Organisation introuvable.")
			c.Abort()
			return
		}
		if k.OrganizationID != orgID {
			response.Forbidden(c, "Clé API non autorisée pour cette organisation.")
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

func getKey(c *gin.Context) *apikeytypes.APIKey {
	if v, ok := c.Get(keyCtxKey); ok {
		if k, ok := v.(*apikeytypes.APIKey); ok {
			return k
		}
	}
	return nil
}

func hasPermission(k *apikeytypes.APIKey, perm string) bool {
	var perms apikeytypes.Permissions
	if err := json.Unmarshal(k.Permissions, &perms); err != nil {
		return false
	}
	switch perm {
	case "read":
		return perms.Read
	case "create":
		return perms.Create
	case "update":
		return perms.Update
	case "delete":
		return perms.Delete
	}
	return false
}
