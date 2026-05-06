package integration_test

import (
	"net/http"
	"testing"
)

// Crée une connexion dans l'organisation et retourne son ID.
func createConnection(t *testing.T, token, slug string) string {
	t.Helper()
	w := do(http.MethodPost, "/v1/organizations/"+slug+"/connections", map[string]any{
		"name":     "Test DB",
		"driver":   "postgres",
		"host":     "localhost",
		"port":     5432,
		"database": "testdb",
		"user":     "postgres",
		"password": "secret",
		"ssl_mode": "disable",
	}, token)
	assertStatus(t, w, http.StatusCreated)
	var resp struct {
		Data struct{ ID string `json:"id"` } `json:"data"`
	}
	decode(t, w, &resp)
	return resp.Data.ID
}

// Vérifie que le propriétaire peut créer une connexion.
func TestConnections_Create_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Conn Org")

	w := do(http.MethodPost, "/v1/organizations/"+slug+"/connections", map[string]any{
		"name":     "Prod DB",
		"driver":   "postgres",
		"host":     "db.example.com",
		"port":     5432,
		"database": "myapp",
		"user":     "admin",
		"password": "pass123",
		"ssl_mode": "disable",
	}, token)
	assertStatus(t, w, http.StatusCreated)

	var resp struct {
		Data struct {
			Name   string `json:"name"`
			Driver string `json:"driver"`
		} `json:"data"`
	}
	decode(t, w, &resp)
	if resp.Data.Name != "Prod DB" {
		t.Fatalf("nom attendu 'Prod DB', obtenu '%s'", resp.Data.Name)
	}
}

// Vérifie qu'un driver invalide retourne une erreur de validation.
func TestConnections_Create_InvalidDriver(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Conn Org")

	w := do(http.MethodPost, "/v1/organizations/"+slug+"/connections", map[string]any{
		"name":     "Bad DB",
		"driver":   "oracle",
		"host":     "db.example.com",
		"port":     1521,
		"database": "myapp",
		"user":     "admin",
		"password": "pass",
	}, token)
	assertStatus(t, w, http.StatusBadRequest)
}

// Vérifie qu'un utilisateur non membre ne peut pas créer une connexion.
func TestConnections_Create_Forbidden(t *testing.T) {
	truncateTables(t)
	ownerToken := signupAndGetToken(t)
	memberToken, _ := signupSecondUser(t)
	slug := createOrg(t, ownerToken, "Conn Org")

	w := do(http.MethodPost, "/v1/organizations/"+slug+"/connections", map[string]any{
		"name":     "Hacked DB",
		"driver":   "postgres",
		"host":     "evil.com",
		"port":     5432,
		"database": "stolen",
		"user":     "hacker",
		"password": "pwned",
	}, memberToken)
	assertStatus(t, w, http.StatusForbidden)
}

// Vérifie que le propriétaire peut lister les connexions.
func TestConnections_List_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Conn Org")
	createConnection(t, token, slug)

	w := do(http.MethodGet, "/v1/organizations/"+slug+"/connections", nil, token)
	assertStatus(t, w, http.StatusOK)

	var resp struct {
		Data []struct{ ID string `json:"id"` } `json:"data"`
	}
	decode(t, w, &resp)
	if len(resp.Data) == 0 {
		t.Fatal("la liste doit contenir au moins une connexion")
	}
}

// Vérifie que le propriétaire peut récupérer une connexion par ID.
func TestConnections_Get_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Conn Org")
	id := createConnection(t, token, slug)

	w := do(http.MethodGet, "/v1/organizations/"+slug+"/connections/"+id, nil, token)
	assertStatus(t, w, http.StatusOK)
}

// Vérifie qu'un ID inexistant retourne 404.
func TestConnections_Get_NotFound(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Conn Org")

	w := do(http.MethodGet, "/v1/organizations/"+slug+"/connections/00000000-0000-0000-0000-000000000000", nil, token)
	assertStatus(t, w, http.StatusNotFound)
}

// Vérifie que le propriétaire peut mettre à jour une connexion.
func TestConnections_Update_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Conn Org")
	id := createConnection(t, token, slug)

	w := do(http.MethodPatch, "/v1/organizations/"+slug+"/connections/"+id, map[string]any{
		"name": "Updated DB",
	}, token)
	assertStatus(t, w, http.StatusOK)

	var resp struct {
		Data struct{ Name string `json:"name"` } `json:"data"`
	}
	decode(t, w, &resp)
	if resp.Data.Name != "Updated DB" {
		t.Fatalf("nom attendu 'Updated DB', obtenu '%s'", resp.Data.Name)
	}
}

// Vérifie que le propriétaire peut supprimer une connexion.
func TestConnections_Delete_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Conn Org")
	id := createConnection(t, token, slug)

	w := do(http.MethodDelete, "/v1/organizations/"+slug+"/connections/"+id, nil, token)
	assertStatus(t, w, http.StatusNoContent)

	// Vérifie que la connexion n'est plus accessible.
	w2 := do(http.MethodGet, "/v1/organizations/"+slug+"/connections/"+id, nil, token)
	assertStatus(t, w2, http.StatusNotFound)
}

// Vérifie que le test de connexion échoue avec des credentials invalides.
func TestConnections_Test_Fail(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Conn Org")
	id := createConnection(t, token, slug)

	w := do(http.MethodPost, "/v1/organizations/"+slug+"/connections/"+id+"/test", nil, token)
	// La connexion à localhost:5432/testdb avec user=postgres échoue en test (DB inexistante).
	assertStatus(t, w, http.StatusBadRequest)
}
