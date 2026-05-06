package integration_test

import (
	"net/http"
	"testing"
)

// Vérifie que l'utilisateur connecté peut récupérer son profil.
func TestUsers_GetMe_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	w := do(http.MethodGet, "/v1/users/me", nil, token)
	assertStatus(t, w, http.StatusOK)

	var resp struct {
		Data struct {
			Email     string `json:"email"`
			FirstName string `json:"first_name"`
		} `json:"data"`
	}
	decode(t, w, &resp)
	if resp.Data.Email != "test@skemacms.com" {
		t.Fatalf("email attendu test@skemacms.com, obtenu %q", resp.Data.Email)
	}
	if resp.Data.FirstName != "Test" {
		t.Fatalf("prénom attendu Test, obtenu %q", resp.Data.FirstName)
	}
}

// Vérifie que le profil est inaccessible sans token d'authentification.
func TestUsers_GetMe_Unauthorized(t *testing.T) {
	truncateTables(t)
	w := do(http.MethodGet, "/v1/users/me", nil, "")
	assertStatus(t, w, http.StatusUnauthorized)
}

// Vérifie que l'utilisateur peut mettre à jour son prénom et son nom.
func TestUsers_UpdateMe_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	w := do(http.MethodPatch, "/v1/users/me", map[string]string{
		"first_name": "Jean",
		"last_name":  "Dupont",
	}, token)
	assertStatus(t, w, http.StatusOK)

	var resp struct {
		Data struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"data"`
	}
	decode(t, w, &resp)
	if resp.Data.FirstName != "Jean" {
		t.Fatalf("prénom attendu Jean, obtenu %q", resp.Data.FirstName)
	}
	if resp.Data.LastName != "Dupont" {
		t.Fatalf("nom attendu Dupont, obtenu %q", resp.Data.LastName)
	}
}

// Vérifie que l'utilisateur peut changer son mot de passe avec le bon mot de passe actuel.
func TestUsers_ChangePassword_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	w := do(http.MethodPost, "/v1/users/me/password", map[string]string{
		"current_password": "TestP@ss123!",
		"new_password":     "N3wP@ss456!",
	}, token)
	assertStatus(t, w, http.StatusOK)
}

// Vérifie qu'un mauvais mot de passe actuel empêche le changement.
func TestUsers_ChangePassword_WrongCurrent(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	w := do(http.MethodPost, "/v1/users/me/password", map[string]string{
		"current_password": "mauvais-mot-de-passe",
		"new_password":     "N3wP@ss456!",
	}, token)
	assertStatus(t, w, http.StatusBadRequest)
}

// Vérifie qu'on ne peut pas définir le même mot de passe que l'actuel.
func TestUsers_ChangePassword_SamePassword(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	w := do(http.MethodPost, "/v1/users/me/password", map[string]string{
		"current_password": "TestP@ss123!",
		"new_password":     "TestP@ss123!",
	}, token)
	assertStatus(t, w, http.StatusBadRequest)
}
