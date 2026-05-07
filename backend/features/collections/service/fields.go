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

func (s *Service) AddField(ctx context.Context, requesterID, orgSlug, collectionID string, req types.AddFieldRequest) (*types.Field, error) {
	if !constants.ValidFieldTypes[req.Type] {
		return nil, errors.New("type de champ invalide")
	}

	c, err := s.Get(ctx, requesterID, orgSlug, collectionID)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.ColumnExists(ctx, c.ID, req.ColumnName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ErrColumnNameTaken)
	}

	f := &types.Field{
		ID:           uuid.New().String(),
		CollectionID: c.ID,
		Name:         req.Name,
		ColumnName:   req.ColumnName,
		Type:         req.Type,
		Required:     req.Required,
		IsUnique:     req.IsUnique,
		DefaultValue: req.DefaultValue,
		Options:      req.Options,
		Position:     req.Position,
		CreatedAt:    time.Now(),
	}

	conn, err := s.connSvc.OpenConduit(ctx, c.ConnectionID)
	if err != nil {
		return nil, errors.New(constants.ErrSchemaFailed)
	}
	defer conn.Close()

	builder, err := schema.New(conn.Driver())
	if err != nil {
		return nil, errors.New(constants.ErrSchemaFailed)
	}
	if _, err := conn.Exec(ctx, builder.AddColumn(c.TableName, f)); err != nil {
		return nil, errors.New(constants.ErrSchemaFailed)
	}

	if err := s.repo.AddField(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) RemoveField(ctx context.Context, requesterID, orgSlug, collectionID, fieldID string) error {
	c, err := s.Get(ctx, requesterID, orgSlug, collectionID)
	if err != nil {
		return err
	}

	f, err := s.repo.FindField(ctx, fieldID)
	if err != nil {
		return err
	}
	if f == nil || f.CollectionID != c.ID {
		return errors.New(constants.ErrFieldNotFound)
	}

	conn, err := s.connSvc.OpenConduit(ctx, c.ConnectionID)
	if err != nil {
		return errors.New(constants.ErrSchemaFailed)
	}
	defer conn.Close()

	builder, err := schema.New(conn.Driver())
	if err != nil {
		return errors.New(constants.ErrSchemaFailed)
	}
	if _, err := conn.Exec(ctx, builder.DropColumn(c.TableName, f.ColumnName)); err != nil {
		return errors.New(constants.ErrSchemaFailed)
	}

	return s.repo.DeleteField(ctx, f.ID)
}
