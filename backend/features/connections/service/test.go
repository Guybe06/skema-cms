package service

import (
	"context"
	"errors"

	"skema-api/core/conduit/adapters/factory"
	"skema-api/core/crypto"
	"skema-api/features/connections/constants"
)

func (s *Service) TestConnection(ctx context.Context, requesterID, orgSlug, id string) error {
	c, err := s.Get(ctx, requesterID, orgSlug, id)
	if err != nil {
		return err
	}

	rec, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New(constants.ErrConnectionFailed)
	}
	password, err := crypto.Decrypt(rec.PasswordEncrypted, s.encryptionKey)
	if err != nil {
		return errors.New(constants.ErrConnectionFailed)
	}

	dsn := buildDSN(c.Driver, c.Host, c.Port, c.Database, c.User, password, c.SSLMode)
	conn, err := factory.New(ctx, c.Driver, dsn)
	if err != nil {
		return errors.New(constants.ErrConnectionFailed)
	}
	defer conn.Close()

	if err := conn.Ping(ctx); err != nil {
		return errors.New(constants.ErrConnectionFailed)
	}
	return nil
}
