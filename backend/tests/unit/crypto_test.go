package unit_test

import (
	"testing"

	"skema-api/core/crypto"
)

// Vérifie que le chiffrement suivi du déchiffrement retourne la valeur d'origine.
func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := "ma-cle-secrete-32-caracteres-ici"
	original := "mot-de-passe-client-super-secret"

	encrypted, err := crypto.Encrypt(original, key)
	if err != nil {
		t.Fatalf("chiffrement échoué : %v", err)
	}
	if encrypted == original {
		t.Fatal("la valeur chiffrée ne doit pas être identique à l'originale")
	}

	decrypted, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("déchiffrement échoué : %v", err)
	}
	if decrypted != original {
		t.Fatalf("valeur attendue %q, obtenue %q", original, decrypted)
	}
}

// Vérifie que deux chiffrements du même texte produisent des résultats différents (nonce aléatoire).
func TestEncrypt_NonDeterministic(t *testing.T) {
	key := "ma-cle-secrete-32-caracteres-ici"
	value := "same-value"

	a, _ := crypto.Encrypt(value, key)
	b, _ := crypto.Encrypt(value, key)

	if a == b {
		t.Fatal("deux chiffrements du même texte doivent produire des résultats différents")
	}
}

// Vérifie que le déchiffrement avec une mauvaise clé retourne une erreur.
func TestDecrypt_WrongKey(t *testing.T) {
	encrypted, _ := crypto.Encrypt("secret", "bonne-cle-de-32-caracteres-aaaaa")
	_, err := crypto.Decrypt(encrypted, "mauvaise-cle-32-caracteres-bbbbb")
	if err == nil {
		t.Fatal("le déchiffrement avec une mauvaise clé doit retourner une erreur")
	}
}

// Vérifie que le déchiffrement d'une chaîne corrompue retourne une erreur.
func TestDecrypt_InvalidCiphertext(t *testing.T) {
	_, err := crypto.Decrypt("ceci-nest-pas-un-ciphertext-valide", "cle-de-32-caracteres-pour-test-aa")
	if err == nil {
		t.Fatal("le déchiffrement d'un ciphertext invalide doit retourner une erreur")
	}
}
