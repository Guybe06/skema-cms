package unit_test

import (
	"strings"
	"testing"

	"skema-api/features/accounts/helpers"
)

// Vérifie que le hash d'un mot de passe est différent de la valeur d'origine.
func TestHashPassword_ProducesHash(t *testing.T) {
	hash, err := helpers.HashPassword("monMotDePasse123!")
	if err != nil {
		t.Fatalf("hashage échoué : %v", err)
	}
	if hash == "monMotDePasse123!" {
		t.Fatal("le hash ne doit pas être identique au mot de passe")
	}
	if !strings.HasPrefix(hash, "$2a$") {
		t.Fatal("le hash doit être au format bcrypt")
	}
}

// Vérifie que CheckPassword accepte le bon mot de passe et rejette un faux.
func TestCheckPassword(t *testing.T) {
	hash, _ := helpers.HashPassword("correct-password")

	if !helpers.CheckPassword(hash, "correct-password") {
		t.Fatal("CheckPassword doit retourner true pour le bon mot de passe")
	}
	if helpers.CheckPassword(hash, "wrong-password") {
		t.Fatal("CheckPassword doit retourner false pour un mauvais mot de passe")
	}
}

// Vérifie que GenerateToken retourne un token brut et un hash différents.
func TestGenerateToken_RawAndHashDiffer(t *testing.T) {
	raw, hashed, err := helpers.GenerateToken()
	if err != nil {
		t.Fatalf("génération du token échouée : %v", err)
	}
	if raw == "" || hashed == "" {
		t.Fatal("le token brut et son hash ne doivent pas être vides")
	}
	if raw == hashed {
		t.Fatal("le token brut et son hash doivent être différents")
	}
}

// Vérifie que HashToken produit toujours le même hash pour le même token.
func TestHashToken_Deterministic(t *testing.T) {
	raw, _, _ := helpers.GenerateToken()
	h1 := helpers.HashToken(raw)
	h2 := helpers.HashToken(raw)
	if h1 != h2 {
		t.Fatal("HashToken doit être déterministe pour un même token")
	}
}

// Vérifie que deux tokens générés sont toujours différents.
func TestGenerateToken_Unique(t *testing.T) {
	raw1, _, _ := helpers.GenerateToken()
	raw2, _, _ := helpers.GenerateToken()
	if raw1 == raw2 {
		t.Fatal("deux tokens générés doivent être différents")
	}
}

// Vérifie que GenerateJWT retourne un token JWT non vide sans erreur.
func TestGenerateJWT(t *testing.T) {
	token, err := helpers.GenerateJWT("user-id-123", "session-id-456", "secret-jwt-pour-les-tests")
	if err != nil {
		t.Fatalf("génération du JWT échouée : %v", err)
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("un JWT doit avoir 3 parties séparées par des points, obtenu : %d", len(parts))
	}
}
