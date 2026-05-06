package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"skema-api/core/conduit/adapters/factory"
	"skema-api/core/crypto"
	"skema-api/features/connections/constants"
	"skema-api/features/connections/types"
)

func (s *Service) Create(ctx context.Context, requesterID, orgSlug string, req types.CreateConnectionRequest) (*types.Connection, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New("organisation introuvable")
	}
	if err := s.checkAccess(ctx, org.ID, org.OwnerID, requesterID); err != nil {
		return nil, err
	}

	rec, err := s.encrypt(req.Host, req.Port, req.Database, req.User, req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sslMode := req.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	rec.ID = uuid.New().String()
	rec.OrganizationID = org.ID
	rec.Name = req.Name
	rec.Driver = req.Driver
	rec.SSLMode = sslMode
	rec.CreatedAt = now
	rec.UpdatedAt = now

	if err := s.repo.Create(ctx, rec); err != nil {
		return nil, err
	}
	return s.toConn(rec, req.Host, req.Port, req.Database, req.User), nil
}

func (s *Service) List(ctx context.Context, requesterID, orgSlug string) ([]*types.Connection, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New("organisation introuvable")
	}
	if err := s.checkAccess(ctx, org.ID, org.OwnerID, requesterID); err != nil {
		return nil, err
	}

	recs, err := s.repo.ListByOrg(ctx, org.ID)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Connection, 0, len(recs))
	for _, rec := range recs {
		c, err := s.decryptRecord(rec)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (s *Service) Get(ctx context.Context, requesterID, orgSlug, id string) (*types.Connection, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New("organisation introuvable")
	}
	if err := s.checkAccess(ctx, org.ID, org.OwnerID, requesterID); err != nil {
		return nil, err
	}

	rec, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rec == nil || rec.OrganizationID != org.ID {
		return nil, errors.New(constants.ErrConnectionNotFound)
	}
	return s.decryptRecord(rec)
}

func (s *Service) Update(ctx context.Context, requesterID, orgSlug, id string, req types.UpdateConnectionRequest) (*types.Connection, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New("organisation introuvable")
	}
	if err := s.checkAccess(ctx, org.ID, org.OwnerID, requesterID); err != nil {
		return nil, err
	}

	rec, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rec == nil || rec.OrganizationID != org.ID {
		return nil, errors.New(constants.ErrConnectionNotFound)
	}

	// Déchiffre les champs existants pour pouvoir les mettre à jour partiellement.
	current, err := s.decryptRecord(rec)
	if err != nil {
		return nil, err
	}

	host := fallback(req.Host, current.Host)
	port := current.Port
	if req.Port != 0 {
		port = req.Port
	}
	database := fallback(req.Database, current.Database)
	user := fallback(req.User, current.User)
	if req.Name != "" {
		rec.Name = req.Name
	}
	if req.SSLMode != "" {
		rec.SSLMode = req.SSLMode
	}

	password := req.Password
	if password == "" {
		// Réutilise le mot de passe chiffré existant.
		newRec, err := s.encrypt(host, port, database, user, "placeholder")
		if err != nil {
			return nil, err
		}
		newRec.PasswordEncrypted = rec.PasswordEncrypted
		newRec.ID = rec.ID
		newRec.OrganizationID = rec.OrganizationID
		newRec.Name = rec.Name
		newRec.Driver = rec.Driver
		newRec.SSLMode = rec.SSLMode
		newRec.CreatedAt = rec.CreatedAt
		if err := s.repo.Update(ctx, newRec); err != nil {
			return nil, err
		}
		return s.toConn(newRec, host, port, database, user), nil
	}

	newRec, err := s.encrypt(host, port, database, user, password)
	if err != nil {
		return nil, err
	}
	newRec.ID = rec.ID
	newRec.OrganizationID = rec.OrganizationID
	newRec.Name = rec.Name
	newRec.Driver = rec.Driver
	newRec.SSLMode = rec.SSLMode
	newRec.CreatedAt = rec.CreatedAt

	if err := s.repo.Update(ctx, newRec); err != nil {
		return nil, err
	}
	return s.toConn(newRec, host, port, database, user), nil
}

func (s *Service) Delete(ctx context.Context, requesterID, orgSlug, id string) error {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return errors.New("organisation introuvable")
	}
	if err := s.checkAccess(ctx, org.ID, org.OwnerID, requesterID); err != nil {
		return err
	}

	rec, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if rec == nil || rec.OrganizationID != org.ID {
		return errors.New(constants.ErrConnectionNotFound)
	}
	return s.repo.Delete(ctx, rec.ID)
}

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
		ID:             rec.ID,
		OrganizationID: rec.OrganizationID,
		Name:           rec.Name,
		Driver:         rec.Driver,
		Host:           host,
		Port:           port,
		Database:       database,
		User:           user,
		SSLMode:        rec.SSLMode,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      rec.UpdatedAt,
	}
}

func buildDSN(driver, host string, port int, database, user, password, sslMode string) string {
	if sslMode == "" {
		sslMode = "disable"
	}
	switch driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database)
	default:
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, password, host, port, database, sslMode)
	}
}

func fallback(val, def string) string {
	if val != "" {
		return val
	}
	return def
}
