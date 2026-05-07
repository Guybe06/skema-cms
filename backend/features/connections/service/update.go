package service

import (
	"context"
	"errors"

	"skema-api/features/connections/constants"
	"skema-api/features/connections/types"
)

func (s *Service) Update(ctx context.Context, requesterID, orgSlug, id string, req types.UpdateConnectionRequest) (*types.Connection, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New(constants.ErrOrgNotFound)
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

	newRec, err := s.encrypt(host, port, database, user, req.Password)
	if err != nil {
		return nil, err
	}
	if req.Password == "" {
		newRec.PasswordEncrypted = rec.PasswordEncrypted
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
