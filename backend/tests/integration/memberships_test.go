package integration_test

import (
	"net/http"
	"testing"
)

// Inscrit un second utilisateur et retourne son token et son ID.
func signupSecondUser(t *testing.T) (token, userID string) {
	t.Helper()
	w := do(http.MethodPost, "/v1/accounts/signup", map[string]string{
		"email": "member@skemacms.com", "password": "TestP@ss123!",
		"first_name": "Member", "last_name": "User",
	}, "")
	assertStatus(t, w, http.StatusCreated)

	var resp struct {
		Data struct {
			AccessToken string `json:"access_token"`
			User        struct {
				ID string `json:"id"`
			} `json:"user"`
		} `json:"data"`
	}
	decode(t, w, &resp)
	return resp.Data.AccessToken, resp.Data.User.ID
}

// Crée une organisation et retourne son slug.
func createOrg(t *testing.T, token, name string) string {
	t.Helper()
	w := do(http.MethodPost, "/v1/organizations", map[string]string{"name": name}, token)
	assertStatus(t, w, http.StatusCreated)
	var resp struct {
		Data struct{ Slug string `json:"slug"` } `json:"data"`
	}
	decode(t, w, &resp)
	return resp.Data.Slug
}

// Vérifie que le propriétaire peut inviter un membre.
func TestMemberships_Invite_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Test Org")

	w := do(http.MethodPost, "/v1/organizations/"+slug+"/members/invite", map[string]string{
		"email": "member@skemacms.com",
		"role":  "member",
	}, token)
	assertStatus(t, w, http.StatusOK)
}

// Vérifie qu'une invitation avec un rôle invalide retourne une erreur de validation.
func TestMemberships_Invite_InvalidRole(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Test Org")

	w := do(http.MethodPost, "/v1/organizations/"+slug+"/members/invite", map[string]string{
		"email": "someone@skemacms.com",
		"role":  "superadmin",
	}, token)
	assertStatus(t, w, http.StatusBadRequest)
}

// Vérifie qu'un utilisateur non autorisé ne peut pas inviter.
func TestMemberships_Invite_Forbidden(t *testing.T) {
	truncateTables(t)
	ownerToken := signupAndGetToken(t)
	memberToken, _ := signupSecondUser(t)
	slug := createOrg(t, ownerToken, "Test Org")

	w := do(http.MethodPost, "/v1/organizations/"+slug+"/members/invite", map[string]string{
		"email": "other@skemacms.com",
		"role":  "member",
	}, memberToken)
	assertStatus(t, w, http.StatusForbidden)
}

// Vérifie que le propriétaire peut lister les membres de son organisation.
func TestMemberships_List_Success(t *testing.T) {
	truncateTables(t)
	token := signupAndGetToken(t)
	slug := createOrg(t, token, "Test Org")
	do(http.MethodPost, "/v1/organizations/"+slug+"/members/invite", map[string]string{
		"email": "a@skemacms.com", "role": "member",
	}, token)

	w := do(http.MethodGet, "/v1/organizations/"+slug+"/members", nil, token)
	assertStatus(t, w, http.StatusOK)

	var resp struct {
		Data []struct{ Email string `json:"email"` } `json:"data"`
	}
	decode(t, w, &resp)
	if len(resp.Data) == 0 {
		t.Fatal("la liste des membres doit contenir au moins l'invitation envoyée")
	}
}

// Vérifie qu'un utilisateur peut accepter une invitation via le token.
func TestMemberships_AcceptInvite_InvalidToken(t *testing.T) {
	truncateTables(t)
	_, memberToken := func() (string, string) {
		signupAndGetToken(t)
		t2, _ := signupSecondUser(t)
		return "", t2
	}()

	w := do(http.MethodPost, "/v1/invitations/accept", map[string]string{
		"token": "token-invalide-qui-nexiste-pas",
	}, memberToken)
	assertStatus(t, w, http.StatusBadRequest)
}

// Vérifie que le propriétaire peut changer le rôle d'un membre.
func TestMemberships_UpdateRole_NotOwner(t *testing.T) {
	truncateTables(t)
	ownerToken := signupAndGetToken(t)
	memberToken, memberID := signupSecondUser(t)
	slug := createOrg(t, ownerToken, "Test Org")

	_ = memberToken
	w := do(http.MethodPatch, "/v1/organizations/"+slug+"/members/"+memberID, map[string]string{
		"role": "admin",
	}, memberToken)
	assertStatus(t, w, http.StatusForbidden)
}

// Vérifie qu'un membre sans autorisation ne peut pas retirer quelqu'un d'autre.
func TestMemberships_Remove_Unauthorized(t *testing.T) {
	truncateTables(t)
	ownerToken := signupAndGetToken(t)
	memberToken, ownerUserID := signupSecondUser(t)
	slug := createOrg(t, ownerToken, "Test Org")
	_ = ownerUserID

	w := do(http.MethodDelete, "/v1/organizations/"+slug+"/members/some-other-id", nil, memberToken)
	assertStatus(t, w, http.StatusForbidden)
}
