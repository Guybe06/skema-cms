package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

/*
 * Encrypt chiffre une valeur texte avec AES-256-GCM.
 *
 * Attend  : la valeur à chiffrer et la clé secrète.
 * Retourne: la valeur chiffrée encodée en hexadécimal, ou une erreur.
 */

func Encrypt(value, key string) (string, error) {
	block, err := newCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)
	return hex.EncodeToString(ciphertext), nil
}

/*
 * Decrypt déchiffre une valeur chiffrée avec AES-256-GCM.
 *
 * Attend  : la valeur chiffrée en hex et la clé secrète utilisée au chiffrement.
 * Retourne: la valeur originale en clair, ou une erreur.
 */

func Decrypt(encrypted, key string) (string, error) {
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	block, err := newCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	plaintext, err := gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}

func newCipher(key string) (cipher.Block, error) {
	hash := sha256.Sum256([]byte(key))
	return aes.NewCipher(hash[:])
}
