package crypto

import "errors"

var (
	ErrInvalidCiphertext = errors.New("texte chiffré invalide")
	ErrDecryptionFailed  = errors.New("échec du déchiffrement")
	ErrInvalidKey        = errors.New("clé de chiffrement invalide")
)
