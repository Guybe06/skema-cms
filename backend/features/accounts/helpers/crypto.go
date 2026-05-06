package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"skema-api/core/middleware/auth"
	"skema-api/features/accounts/constants"
)

var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-timing-protection"), constants.BcryptCost)

/*
 * HashPassword génère un hash bcrypt du mot de passe fourni.
 *
 * Attend  : le mot de passe en clair.
 * Retourne: le hash bcrypt ou une erreur.
 */

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), constants.BcryptCost)
	return string(hash), err
}

/*
 * CheckPassword compare un hash bcrypt avec un mot de passe en clair.
 *
 * Attend  : le hash stocké et le mot de passe soumis.
 * Retourne: true si le mot de passe correspond.
 */

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// HashDummy exécute un bcrypt fictif pour neutraliser les timing attacks.
func HashDummy() {
	_ = bcrypt.CompareHashAndPassword(dummyHash, []byte("dummy"))
}

/*
 * GenerateToken génère un token aléatoire sécurisé et son hash SHA-256.
 *
 * Attend  : aucun paramètre.
 * Retourne: le token brut (pour le client) et son hash (pour la base de données).
 */

func GenerateToken() (raw, hashed string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	raw = hex.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	hashed = hex.EncodeToString(sum[:])
	return raw, hashed, nil
}

// HashToken retourne le hash SHA-256 d'un token brut.
func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

/*
 * GenerateJWT génère un token JWT signé avec les claims utilisateur.
 *
 * Attend  : l'identifiant utilisateur, l'identifiant de session et le secret JWT.
 * Retourne: le token signé ou une erreur.
 */

func GenerateJWT(userID, sessionID, secret string) (string, error) {
	claims := auth.Claims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}
