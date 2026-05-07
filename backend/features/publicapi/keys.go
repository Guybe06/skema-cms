package publicapi

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	apikeytypes "skema-api/features/apikeys/types"
)

func getKey(c *gin.Context) *apikeytypes.APIKey {
	if v, ok := c.Get(keyCtxKey); ok {
		if k, ok := v.(*apikeytypes.APIKey); ok {
			return k
		}
	}
	return nil
}

/*
 * hasPermission vérifie qu'une clé API dispose de la permission demandée.
 *
 * Attend  : une clé API et le nom de la permission (read, create, update, delete).
 * Retourne: true si la permission est accordée.
 */

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
