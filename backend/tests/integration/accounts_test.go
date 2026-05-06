package integration_test

import (
	"net/http"
	"testing"
)

// Vérifie qu'un nouvel utilisateur peut s'inscrire et reçoit un token d'accès.
func TestAccounts_Signup_Success(t *testing.T) {
	truncateTables(t)
	w := do(http.MethodPost, "/v1/accounts/signup", map[string]string{
		"email":      "jane@skemacms.com",
		"password":   "S3cur3P@ss!",
		"first_name": "Jane",
		"last_name":  "Doe",
	}, "")
	assertStatus(t, w, http.StatusCreated)

	var resp struct {
		Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}
	decode(t, w, &resp)
	if resp.Data.AccessToken == "" {
		t.Fatal("un access_token doit être retourné après l'inscription")
	}
	if resp.Data.RefreshToken == "" {
		t.Fatal("un refresh_token doit être retourné après l'inscription")
	}
}

// Vérifie qu'on ne peut pas créer deux comptes avec le même email.
func TestAccounts_Signup_DuplicateEmail(t *testing.T) {
	truncateTables(t)
	body := map[string]string{
		"email": "dup@skemacms.com", "password": "S3cur3P@ss!",
		"first_name": "A", "last_name": "B",
	}
	do(http.MethodPost, "/v1/accounts/signup", body, "")
	w := do(http.MethodPost, "/v1/accounts/signup", body, "")
	assertStatus(t, w, http.StatusConflict)
}

// Vérifie que les champs obligatoires sont validés à l'inscription.
func TestAccounts_Signup_ValidationErrors(t *testing.T) {
	truncateTables(t)
	w := do(http.MethodPost, "/v1/accounts/signup", map[string]string{
		"email": "pas-un-email", "password": "court",
	}, "")
	assertStatus(t, w, http.StatusBadRequest)
}

// Vérifie qu'un utilisateur peut se connecter avec ses identifiants.
func TestAccounts_Signin_Success(t *testing.T) {
	truncateTables(t)
	signupAndGetToken(t)

	w := do(http.MethodPost, "/v1/accounts/signin", map[string]string{
		"email":    "test@skemacms.com",
		"password": "TestP@ss123!",
	}, "")
	assertStatus(t, w, http.StatusOK)
}

// Vérifie qu'un mauvais mot de passe retourne une erreur 401.
func TestAccounts_Signin_WrongPassword(t *testing.T) {
	truncateTables(t)
	signupAndGetToken(t)

	w := do(http.MethodPost, "/v1/accounts/signin", map[string]string{
		"email":    "test@skemacms.com",
		"password": "mauvais-mot-de-passe",
	}, "")
	assertStatus(t, w, http.StatusUnauthorized)
}

// Vérifie qu'un email inconnu retourne une erreur 401 sans divulguer son existence.
func TestAccounts_Signin_UnknownEmail(t *testing.T) {
	truncateTables(t)
	w := do(http.MethodPost, "/v1/accounts/signin", map[string]string{
		"email":    "inconnu@skemacms.com",
		"password": "n-importe-quoi",
	}, "")
	assertStatus(t, w, http.StatusUnauthorized)
}

// Vérifie qu'un utilisateur connecté peut se déconnecter.
func TestAccounts_Signout_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)

	w := do(http.MethodPost, "/v1/accounts/signout", nil, token)
	assertStatus(t, w, http.StatusNoContent)
}

// Vérifie que la déconnexion sans token retourne une erreur 401.
func TestAccounts_Signout_Unauthorized(t *testing.T) {
	truncateTables(t)
	w := do(http.MethodPost, "/v1/accounts/signout", nil, "")
	assertStatus(t, w, http.StatusUnauthorized)
}

// Vérifie que le refresh token permet d'obtenir une nouvelle paire de tokens.
func TestAccounts_Refresh_Success(t *testing.T) {
	truncateTables(t)

	var signup struct {
		Data struct {
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}
	w := do(http.MethodPost, "/v1/accounts/signup", map[string]string{
		"email": "ref@skemacms.com", "password": "S3cur3P@ss!",
		"first_name": "Ref", "last_name": "User",
	}, "")
	decode(t, w, &signup)

	w2 := do(http.MethodPost, "/v1/accounts/refresh", map[string]string{
		"refresh_token": signup.Data.RefreshToken,
	}, "")
	assertStatus(t, w2, http.StatusOK)

	var refreshed struct {
		Data struct{ AccessToken string `json:"access_token"` } `json:"data"`
	}
	decode(t, w2, &refreshed)
	if refreshed.Data.AccessToken == "" {
		t.Fatal("un nouvel access_token doit être retourné après le refresh")
	}
}

// Vérifie qu'un refresh token invalide retourne une erreur 401.
func TestAccounts_Refresh_InvalidToken(t *testing.T) {
	truncateTables(t)
	w := do(http.MethodPost, "/v1/accounts/refresh", map[string]string{
		"refresh_token": "token-completement-invalide",
	}, "")
	assertStatus(t, w, http.StatusUnauthorized)
}
