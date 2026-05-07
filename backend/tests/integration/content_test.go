package integration_test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestContent(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	// Org
	w := do(http.MethodPost, "/v1/organizations", map[string]string{"name": "ContentOrg"}, token)
	assertStatus(t, w, http.StatusCreated)
	var orgResp struct{ Data struct{ Slug string `json:"slug"` } `json:"data"` }
	decode(t, w, &orgResp)
	slug := orgResp.Data.Slug

	// Connexion
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

	// Collection avec un champ titre
	w = do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/collections", slug), map[string]any{
		"name": "Posts", "table_name": "test_posts", "connection_id": connID,
	}, token)
	assertStatus(t, w, http.StatusCreated)
	var collResp struct{ Data struct{ ID string `json:"id"` } `json:"data"` }
	decode(t, w, &collResp)
	collID := collResp.Data.ID

	// Ajouter un champ titre
	w = do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/collections/%s/fields", slug, collID), map[string]any{
		"name": "titre", "column_name": "titre", "type": "text", "required": true,
	}, token)
	assertStatus(t, w, http.StatusCreated)

	baseURL := fmt.Sprintf("/v1/organizations/%s/collections/%s/content", slug, collID)
	var entryID string

	t.Run("créer une entrée", func(t *testing.T) {
		w := do(http.MethodPost, baseURL, map[string]any{
			"titre": "Mon premier article",
		}, token)
		assertStatus(t, w, http.StatusCreated)
		logResponse(t, "POST /content", w)

		var resp struct {
			Data struct{ ID string `json:"id"` } `json:"data"`
		}
		decode(t, w, &resp)
		entryID = resp.Data.ID
	})

	t.Run("lister les entrées", func(t *testing.T) {
		w := do(http.MethodGet, baseURL, nil, token)
		assertStatus(t, w, http.StatusOK)
		logResponse(t, "GET /content", w)
	})

	t.Run("récupérer une entrée par ID", func(t *testing.T) {
		if entryID == "" {
			t.Skip("entryID vide")
		}
		w := do(http.MethodGet, baseURL+"/"+entryID, nil, token)
		assertStatus(t, w, http.StatusOK)
		logResponse(t, "GET /content/:id", w)
	})

	t.Run("modifier une entrée", func(t *testing.T) {
		if entryID == "" {
			t.Skip("entryID vide")
		}
		w := do(http.MethodPatch, baseURL+"/"+entryID, map[string]any{
			"titre": "Titre modifié",
		}, token)
		assertStatus(t, w, http.StatusOK)
		logResponse(t, "PATCH /content/:id", w)
	})

	t.Run("supprimer une entrée", func(t *testing.T) {
		if entryID == "" {
			t.Skip("entryID vide")
		}
		w := do(http.MethodDelete, baseURL+"/"+entryID, nil, token)
		assertStatus(t, w, http.StatusNoContent)
		logResponse(t, "DELETE /content/:id", w)
	})

	t.Run("404 sur entrée inexistante", func(t *testing.T) {
		w := do(http.MethodGet, baseURL+"/00000000-0000-0000-0000-000000000000", nil, token)
		assertStatus(t, w, http.StatusNotFound)
		logResponse(t, "GET /content/:id 404", w)
	})
}
