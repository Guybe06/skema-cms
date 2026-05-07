package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"skema-api/features/collections/constants"
	"skema-api/features/collections/schema"
	"skema-api/features/collections/types"
)

func (s *Service) Create(ctx context.Context, requesterID, orgSlug string, req types.CreateCollectionRequest) (*types.Collection, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New(constants.ErrOrgNotFound)
	}
	if err := s.checkAccess(ctx, org.ID, org.OwnerID, requesterID); err != nil {
		return nil, err
	}

	exists, err := s.repo.TableExists(ctx, req.ConnectionID, req.TableName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ErrTableNameTaken)
	}

	conn, err := s.connSvc.OpenConduit(ctx, req.ConnectionID)
	if err != nil {
		return nil, errors.New(constants.ErrSchemaFailed)
	}
	defer conn.Close()

	builder, err := schema.New(conn.Driver())
	if err != nil {
		return nil, errors.New(constants.ErrSchemaFailed)
	}
	if _, err := conn.Exec(ctx, builder.CreateTable(req.TableName)); err != nil {
		return nil, errors.New(constants.ErrSchemaFailed)
	}

	now := time.Now()
	c := &types.Collection{
		ID: uuid.New().String(), ConnectionID: req.ConnectionID,
		OrganizationID: org.ID, Name: req.Name, TableName: req.TableName,
		DisplayName: req.DisplayName, Description: req.Description,
		CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) Update(ctx context.Context, requesterID, orgSlug, id string, req types.UpdateCollectionRequest) (*types.Collection, error) {
	c, err := s.Get(ctx, requesterID, orgSlug, id)
	if err != nil {
		return nil, err
	}

	name := req.Name
	if name == "" {
		name = c.Name
	}
	displayName := req.DisplayName
	if displayName == "" {
		displayName = c.DisplayName
	}
	description := req.Description
	if description == "" {
		description = c.Description
	}

	if err := s.repo.Update(ctx, c.ID, name, displayName, description); err != nil {
		return nil, err
	}
	c.Name = name
	c.DisplayName = displayName
	c.Description = description
	return c, nil
}

func (s *Service) Delete(ctx context.Context, requesterID, orgSlug, id string) error {
	c, err := s.Get(ctx, requesterID, orgSlug, id)
	if err != nil {
		return err
	}

	conn, err := s.connSvc.OpenConduit(ctx, c.ConnectionID)
	if err == nil {
		builder, berr := schema.New(conn.Driver())
		if berr == nil {
			conn.Exec(ctx, builder.DropTable(c.TableName))
		}
		conn.Close()
	}

	return s.repo.Delete(ctx, c.ID)
}
