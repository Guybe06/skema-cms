package integration_test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCollections(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	// créer une organisation
	w := do(http.MethodPost, "/v1/organizations", map[string]string{"name": "TestOrg"}, token)
	assertStatus(t, w, http.StatusCreated)
	logResponse(t, "POST /organizations", w)

	var orgResp struct {
		Data struct {
			Slug string `json:"slug"`
		} `json:"data"`
	}
	decode(t, w, &orgResp)
	slug := orgResp.Data.Slug

	// créer une connexion (vers la base de test elle-même)
	w = do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/connections", slug), map[string]any{
		"name":     "Test DB",
		"driver":   "postgres",
		"host":     getenv("CMS_DB_HOST", "localhost"),
		"port":     dbPort(),
		"database": testDBName,
		"user":     getenv("CMS_DB_USER", "postgres"),
		"password": getenv("CMS_DB_PASSWORD", ""),
		"ssl_mode": "disable",
	}, token)
	assertStatus(t, w, http.StatusCreated)
	logResponse(t, "POST /connections", w)

	var connResp struct {
		Data struct{ ID string `json:"id"` } `json:"data"`
	}
	decode(t, w, &connResp)
	connID := connResp.Data.ID

	t.Run("créer une collection", func(t *testing.T) {
		w := do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/collections", slug), map[string]any{
			"name":          "Articles",
			"table_name":    "test_articles",
			"connection_id": connID,
			"display_name":  "Articles de blog",
		}, token)
		assertStatus(t, w, http.StatusCreated)
		logResponse(t, "POST /collections", w)
	})

	t.Run("lister les collections", func(t *testing.T) {
		w := do(http.MethodGet, fmt.Sprintf("/v1/organizations/%s/collections", slug), nil, token)
		assertStatus(t, w, http.StatusOK)
		logResponse(t, "GET /collections", w)
	})

	var collID string
	t.Run("récupérer une collection par ID", func(t *testing.T) {
		wList := do(http.MethodGet, fmt.Sprintf("/v1/organizations/%s/collections", slug), nil, token)
		var listResp struct {
			Data []struct{ ID string `json:"id"` } `json:"data"`
		}
		decode(t, wList, &listResp)
		if len(listResp.Data) == 0 {
			t.Fatal("aucune collection trouvée")
		}
		collID = listResp.Data[0].ID

		w := do(http.MethodGet, fmt.Sprintf("/v1/organizations/%s/collections/%s", slug, collID), nil, token)
		assertStatus(t, w, http.StatusOK)
		logResponse(t, "GET /collections/:id", w)
	})

	t.Run("ajouter un champ texte", func(t *testing.T) {
		w := do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/collections/%s/fields", slug, collID), map[string]any{
			"name":        "titre",
			"column_name": "titre",
			"type":        "text",
			"required":    true,
		}, token)
		assertStatus(t, w, http.StatusCreated)
		logResponse(t, "POST /fields (titre)", w)
	})

	t.Run("ajouter un champ number", func(t *testing.T) {
		w := do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/collections/%s/fields", slug, collID), map[string]any{
			"name":        "vues",
			"column_name": "vues",
			"type":        "number",
		}, token)
		assertStatus(t, w, http.StatusCreated)
		logResponse(t, "POST /fields (vues)", w)
	})

	t.Run("supprimer un champ", func(t *testing.T) {
		// récupérer l'ID du champ vues
		wGet := do(http.MethodGet, fmt.Sprintf("/v1/organizations/%s/collections/%s", slug, collID), nil, token)
		var collResp struct {
			Data struct {
				Fields []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"fields"`
			} `json:"data"`
		}
		decode(t, wGet, &collResp)

		var fieldID string
		for _, f := range collResp.Data.Fields {
			if f.Name == "vues" {
				fieldID = f.ID
			}
		}
		if fieldID == "" {
			t.Skip("champ vues introuvable")
		}

		w := do(http.MethodDelete, fmt.Sprintf("/v1/organizations/%s/collections/%s/fields/%s", slug, collID, fieldID), nil, token)
		assertStatus(t, w, http.StatusNoContent)
		logResponse(t, "DELETE /fields/:fieldId", w)
	})

	t.Run("supprimer une collection", func(t *testing.T) {
		// créer une collection jetable
		wCreate := do(http.MethodPost, fmt.Sprintf("/v1/organizations/%s/collections", slug), map[string]any{
			"name":          "ASupprimer",
			"table_name":    "test_to_delete",
			"connection_id": connID,
		}, token)
		assertStatus(t, wCreate, http.StatusCreated)
		var cr struct{ Data struct{ ID string `json:"id"` } `json:"data"` }
		decode(t, wCreate, &cr)

		w := do(http.MethodDelete, fmt.Sprintf("/v1/organizations/%s/collections/%s", slug, cr.Data.ID), nil, token)
		assertStatus(t, w, http.StatusNoContent)
		logResponse(t, "DELETE /collections/:id", w)
	})
}
