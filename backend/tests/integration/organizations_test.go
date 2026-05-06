package integration_test

import (
	"net/http"
	"testing"
)

// Vérifie qu'un utilisateur connecté peut créer une organisation.
func TestOrgs_Create_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	w := do(http.MethodPost, "/v1/organizations", map[string]string{
		"name": "Acme Corp",
	}, token)
	assertStatus(t, w, http.StatusCreated)

	var resp struct {
		Data struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"data"`
	}
	decode(t, w, &resp)
	if resp.Data.Name != "Acme Corp" {
		t.Fatalf("nom attendu Acme Corp, obtenu %q", resp.Data.Name)
	}
	if resp.Data.Slug != "acme-corp" {
		t.Fatalf("slug attendu acme-corp, obtenu %q", resp.Data.Slug)
	}
}

// Vérifie que deux organisations avec le même nom reçoivent des slugs différents.
func TestOrgs_Create_SlugUnique(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	do(http.MethodPost, "/v1/organizations", map[string]string{"name": "Mon Org"}, token)
	w := do(http.MethodPost, "/v1/organizations", map[string]string{"name": "Mon Org"}, token)
	assertStatus(t, w, http.StatusCreated)

	var resp struct {
		Data struct{ Slug string `json:"slug"` } `json:"data"`
	}
	decode(t, w, &resp)
	if resp.Data.Slug == "mon-org" {
		t.Fatal("le second slug doit être différent du premier (ex: mon-org-2)")
	}
}

// Vérifie que la création d'une organisation requiert une authentification.
func TestOrgs_Create_Unauthorized(t *testing.T) {
	truncateTables(t)
	w := do(http.MethodPost, "/v1/organizations", map[string]string{"name": "Acme"}, "")
	assertStatus(t, w, http.StatusUnauthorized)
}

// Vérifie que l'utilisateur peut lister ses organisations.
func TestOrgs_List_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	do(http.MethodPost, "/v1/organizations", map[string]string{"name": "Org A"}, token)
	do(http.MethodPost, "/v1/organizations", map[string]string{"name": "Org B"}, token)

	w := do(http.MethodGet, "/v1/organizations", nil, token)
	assertStatus(t, w, http.StatusOK)

	var resp struct {
		Data []struct{ Name string `json:"name"` } `json:"data"`
	}
	decode(t, w, &resp)
	if len(resp.Data) != 2 {
		t.Fatalf("2 organisations attendues, obtenu %d", len(resp.Data))
	}
}

// Vérifie qu'on peut récupérer une organisation par son slug.
func TestOrgs_Get_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	do(http.MethodPost, "/v1/organizations", map[string]string{"name": "Skema Inc"}, token)

	w := do(http.MethodGet, "/v1/organizations/skema-inc", nil, token)
	assertStatus(t, w, http.StatusOK)
}

// Vérifie qu'un slug inexistant retourne une erreur 404.
func TestOrgs_Get_NotFound(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	w := do(http.MethodGet, "/v1/organizations/inexistante", nil, token)
	assertStatus(t, w, http.StatusNotFound)
}

// Vérifie que le propriétaire peut mettre à jour le nom de son organisation.
func TestOrgs_Update_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	do(http.MethodPost, "/v1/organizations", map[string]string{"name": "Ancien Nom"}, token)

	w := do(http.MethodPatch, "/v1/organizations/ancien-nom", map[string]string{
		"name": "Nouveau Nom",
	}, token)
	assertStatus(t, w, http.StatusOK)

	var resp struct {
		Data struct{ Slug string `json:"slug"` } `json:"data"`
	}
	decode(t, w, &resp)
	if resp.Data.Slug != "nouveau-nom" {
		t.Fatalf("slug attendu nouveau-nom, obtenu %q", resp.Data.Slug)
	}
}

// Vérifie que le propriétaire peut supprimer son organisation.
func TestOrgs_Delete_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	do(http.MethodPost, "/v1/organizations", map[string]string{"name": "A Supprimer"}, token)

	w := do(http.MethodDelete, "/v1/organizations/a-supprimer", nil, token)
	assertStatus(t, w, http.StatusNoContent)

	w2 := do(http.MethodGet, "/v1/organizations/a-supprimer", nil, token)
	assertStatus(t, w2, http.StatusNotFound)
}
