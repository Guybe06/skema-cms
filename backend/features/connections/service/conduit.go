package service

import (
	"context"
	"errors"
	"strconv"

	"skema-api/core/conduit"
	"skema-api/core/conduit/adapters/factory"
	"skema-api/core/crypto"
)

// OpenConduit déchiffre les credentials et ouvre une connexion Conduit vers la base client.
func (s *Service) OpenConduit(ctx context.Context, connectionID string) (conduit.Conduit, error) {
	rec, err := s.repo.FindByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, errors.New("connexion introuvable")
	}

	host, err := crypto.Decrypt(rec.HostEncrypted, s.encryptionKey)
	if err != nil {
		return nil, err
	}
	portStr, err := crypto.Decrypt(rec.PortEncrypted, s.encryptionKey)
	if err != nil {
		return nil, err
	}
	port, _ := strconv.Atoi(portStr)
	database, err := crypto.Decrypt(rec.DatabaseEncrypted, s.encryptionKey)
	if err != nil {
		return nil, err
	}
	user, err := crypto.Decrypt(rec.UserEncrypted, s.encryptionKey)
	if err != nil {
		return nil, err
	}
	password, err := crypto.Decrypt(rec.PasswordEncrypted, s.encryptionKey)
	if err != nil {
		return nil, err
	}

	dsn := buildDSN(rec.Driver, host, port, database, user, password, rec.SSLMode)
	return factory.New(ctx, rec.Driver, dsn)
}
