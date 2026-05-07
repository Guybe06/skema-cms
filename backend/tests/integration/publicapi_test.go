package integration_test

import (
	"fmt"
	"net/http"
	"testing"
)

// TestPublicAPI teste l'API publique de bout en bout :
// création org → connexion → collection → champ → clé API → CRUD via clé.
func TestPublicAPI(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	// --- Org ---
	w := do(http.MethodPost, "/v1/organizations", map[string]string{"name": "PubOrg"}, token)
	assertStatus(t, w, http.StatusCreated)
	var orgResp struct{ Data struct{ Slug string `json:"slug"` } `json:"data"` }
	decode(t, w, &orgResp)
	slug := orgResp.Data.Slug

	// --- Connexion ---
	w = do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/connections", slug), map[string]any{
		"name": "Test DB", "driver": "postgres",
		"host": getenv("CMS_DB_HOST", "localhost"), "port": dbPort(),
		"database": testDBName, "user": getenv("CMS_DB_USER", "postgres"),
		"password": getenv("CMS_DB_PASSWORD", ""), "ssl_mode": "disable",
	}, token)
	assertStatus(t, w, http.StatusCreated)
	var connResp struct{ Data struct{ ID string `json:"id"` } `json:"data"` }
	decode(t, w, &connResp)
	connID := connResp.Data.ID

	// --- Collection ---
	w = do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/collections", slug), map[string]any{
		"name": "Produits", "table_name": "test_produits", "connection_id": connID,
	}, token)
	assertStatus(t, w, http.StatusCreated)

	// Ajouter un champ nom
	var collListResp struct{ Data []struct{ ID string `json:"id"` } `json:"data"` }
	wList := do(http.MethodGet, fmt.Sprintf("/v1/organizations/%s/collections", slug), nil, token)
	decode(t, wList, &collListResp)
	if len(collListResp.Data) == 0 {
		t.Fatal("aucune collection créée")
	}
	collID := collListResp.Data[0].ID

	w = do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/collections/%s/fields", slug, collID), map[string]any{
		"name": "nom", "column_name": "nom", "type": "text", "required": true,
	}, token)
	assertStatus(t, w, http.StatusCreated)

	// --- Clé API avec toutes les permissions ---
	w = do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/apikeys", slug), map[string]any{
		"name": "Clé test pub",
		"permissions": map[string]bool{
			"read": true, "create": true, "update": true, "delete": true,
		},
	}, token)
	assertStatus(t, w, http.StatusCreated)
	logResponse(t, "POST /apikeys (pub)", w)

	var keyResp struct {
		Data struct {
			RawKey string `json:"raw_key"`
		} `json:"data"`
	}
	decode(t, w, &keyResp)
	apiKey := keyResp.Data.RawKey
	if apiKey == "" {
		t.Fatal("raw_key manquant")
	}

	pubBase := fmt.Sprintf("/v1/pub/%s/test_produits", slug)
	var entryID string

	t.Run("créer une entrée via API publique", func(t *testing.T) {
		w := doWithKey(http.MethodPost, pubBase, map[string]any{
			"nom": "Produit Alpha",
		}, apiKey)
		assertStatus(t, w, http.StatusCreated)
		logResponse(t, "POST /pub/:org/:table", w)

		var resp struct {
			Data struct{ ID string `json:"id"` } `json:"data"`
		}
		decode(t, w, &resp)
		entryID = resp.Data.ID
	})

	t.Run("lister les entrées via API publique", func(t *testing.T) {
		w := doWithKey(http.MethodGet, pubBase, nil, apiKey)
		assertStatus(t, w, http.StatusOK)
		logResponse(t, "GET /pub/:org/:table", w)
	})

	t.Run("récupérer une entrée via API publique", func(t *testing.T) {
		if entryID == "" {
			t.Skip("entryID vide")
		}
		w := doWithKey(http.MethodGet, pubBase+"/"+entryID, nil, apiKey)
		assertStatus(t, w, http.StatusOK)
		logResponse(t, "GET /pub/:org/:table/:id", w)
	})

	t.Run("modifier une entrée via API publique", func(t *testing.T) {
		if entryID == "" {
			t.Skip("entryID vide")
		}
		w := doWithKey(http.MethodPatch, pubBase+"/"+entryID, map[string]any{
			"nom": "Produit Beta",
		}, apiKey)
		assertStatus(t, w, http.StatusOK)
		logResponse(t, "PATCH /pub/:org/:table/:id", w)
	})

	t.Run("supprimer une entrée via API publique", func(t *testing.T) {
		if entryID == "" {
			t.Skip("entryID vide")
		}
		w := doWithKey(http.MethodDelete, pubBase+"/"+entryID, nil, apiKey)
		assertStatus(t, w, http.StatusNoContent)
		logResponse(t, "DELETE /pub/:org/:table/:id", w)
	})

	t.Run("rejeté sans clé", func(t *testing.T) {
		w := doWithKey(http.MethodGet, pubBase, nil, "")
		assertStatus(t, w, http.StatusUnauthorized)
		logResponse(t, "GET /pub sans clé", w)
	})

	t.Run("rejeté avec clé invalide", func(t *testing.T) {
		w := doWithKey(http.MethodGet, pubBase, nil, "sk_live_invalide")
		assertStatus(t, w, http.StatusUnauthorized)
		logResponse(t, "GET /pub clé invalide", w)
	})
}
