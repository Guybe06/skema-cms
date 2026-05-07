package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"skema-api/features/connections/constants"
	"skema-api/features/connections/types"
)

func (s *Service) Create(ctx context.Context, requesterID, orgSlug string, req types.CreateConnectionRequest) (*types.Connection, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New(constants.ErrOrgNotFound)
	}
	if err := s.checkAccess(ctx, org.ID, org.OwnerID, requesterID); err != nil {
		return nil, err
	}

	rec, err := s.encrypt(req.Host, req.Port, req.Database, req.User, req.Password)
	if err != nil {
		return nil, err
	}

	sslMode := req.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	now := time.Now()
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
		return nil, errors.New(constants.ErrOrgNotFound)
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
	return s.decryptRecord(rec)
}

func (s *Service) Delete(ctx context.Context, requesterID, orgSlug, id string) error {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return errors.New(constants.ErrOrgNotFound)
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
