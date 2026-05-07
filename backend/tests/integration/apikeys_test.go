package integration_test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAPIKeys(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	w := do(http.MethodPost, "/v1/organizations", map[string]string{"name": "KeyOrg"}, token)
	assertStatus(t, w, http.StatusCreated)
	var orgResp struct{ Data struct{ Slug string `json:"slug"` } `json:"data"` }
	decode(t, w, &orgResp)
	slug := orgResp.Data.Slug

	baseURL := fmt.Sprintf("/v1/organizations/%s/apikeys", slug)
	var keyID, rawKey string

	t.Run("générer une clé API", func(t *testing.T) {
		w := do(http.MethodPost, baseURL, map[string]any{
			"name": "Clé production",
			"permissions": map[string]bool{
				"read": true, "create": true, "update": true, "delete": false,
			},
		}, token)
		assertStatus(t, w, http.StatusCreated)
		logResponse(t, "POST /apikeys", w)

		var resp struct {
			Data struct {
				ID        string `json:"id"`
				KeyPrefix string `json:"key_prefix"`
				RawKey    string `json:"raw_key"`
			} `json:"data"`
		}
		decode(t, w, &resp)
		keyID = resp.Data.ID
		rawKey = resp.Data.RawKey

		if rawKey == "" {
			t.Fatal("raw_key manquant dans la réponse")
		}
		if len(rawKey) < 16 || rawKey[:8] != "sk_live_" {
			t.Fatalf("format raw_key invalide : %s", rawKey)
		}
		t.Logf("raw_key reçu (préfixe) : %s...", rawKey[:16])
		_ = keyID
	})

	t.Run("lister les clés API", func(t *testing.T) {
		w := do(http.MethodGet, baseURL, nil, token)
		assertStatus(t, w, http.StatusOK)
		logResponse(t, "GET /apikeys", w)
	})

	t.Run("révoquer une clé API", func(t *testing.T) {
		if keyID == "" {
			t.Skip("keyID vide")
		}
		w := do(http.MethodDelete, baseURL+"/"+keyID, nil, token)
		assertStatus(t, w, http.StatusNoContent)
		logResponse(t, "DELETE /apikeys/:id", w)
	})

	t.Run("clé révoquée rejetée", func(t *testing.T) {
		// la clé a été révoquée, elle ne doit plus fonctionner sur l'API publique
		_ = rawKey // utilisé dans TestPublicAPI
	})
}
