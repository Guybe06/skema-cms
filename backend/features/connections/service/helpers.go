package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"skema-api/core/crypto"
	"skema-api/features/connections/constants"
	"skema-api/features/connections/types"
)

func (s *Service) checkAccess(ctx context.Context, orgID, ownerID, requesterID string) error {
	if ownerID == requesterID {
		return nil
	}
	ok, err := s.orgsRepo.IsMember(ctx, orgID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(constants.ErrNotAuthorized)
	}
	return nil
}

func (s *Service) encrypt(host string, port int, database, user, password string) (*types.EncryptedRecord, error) {
	hostEnc, err := crypto.Encrypt(host, s.encryptionKey)
	if err != nil {
		return nil, err
	}
	portEnc, err := crypto.Encrypt(strconv.Itoa(port), s.encryptionKey)
	if err != nil {
		return nil, err
	}
	dbEnc, err := crypto.Encrypt(database, s.encryptionKey)
	if err != nil {
		return nil, err
	}
	userEnc, err := crypto.Encrypt(user, s.encryptionKey)
	if err != nil {
		return nil, err
	}
	passEnc, err := crypto.Encrypt(password, s.encryptionKey)
	if err != nil {
		return nil, err
	}
	return &types.EncryptedRecord{
		HostEncrypted:     hostEnc,
		PortEncrypted:     portEnc,
		DatabaseEncrypted: dbEnc,
		UserEncrypted:     userEnc,
		PasswordEncrypted: passEnc,
	}, nil
}

func (s *Service) decryptRecord(rec *types.EncryptedRecord) (*types.Connection, error) {
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
	return s.toConn(rec, host, port, database, user), nil
}

func (s *Service) toConn(rec *types.EncryptedRecord, host string, port int, database, user string) *types.Connection {
	return &types.Connection{
		ID: rec.ID, OrganizationID: rec.OrganizationID,
		Name: rec.Name, Driver: rec.Driver, Host: host, Port: port,
		Database: database, User: user, SSLMode: rec.SSLMode,
		CreatedAt: rec.CreatedAt, UpdatedAt: rec.UpdatedAt,
	}
}

func buildDSN(driver, host string, port int, database, user, password, sslMode string) string {
	if sslMode == "" {
		sslMode = "disable"
	}
	if driver == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database)
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, password, host, port, database, sslMode)
}

func fallback(val, def string) string {
	if val != "" {
		return val
	}
	return def
}
