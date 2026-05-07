package service

import (
	"context"
	"errors"
	"fmt"

	"skema-api/features/content/constants"
)

func (s *Service) List(ctx context.Context, requesterID, orgSlug, collectionID string, page, perPage int, sort, order string) ([]map[string]any, int, error) {
	c, conn, fields, err := s.getContext(ctx, requesterID, orgSlug, collectionID)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()

	if sort == "" {
		sort = "created_at"
	}
	if order != "asc" {
		order = "desc"
	}

	selectCols := BuildSelectCols(fields)
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s %s LIMIT $1 OFFSET $2",
		selectCols, quoteIdent(c.TableName), quoteIdent(sort), order)

	rows, err := conn.Query(ctx, query, perPage, (page-1)*perPage)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	entries, err := ScanRows(rows)
	if err != nil {
		return nil, 0, err
	}

	var total int
	conn.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", quoteIdent(c.TableName))).Scan(&total)
	return entries, total, nil
}

func (s *Service) Get(ctx context.Context, requesterID, orgSlug, collectionID, entryID string) (map[string]any, error) {
	c, conn, fields, err := s.getContext(ctx, requesterID, orgSlug, collectionID)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return fetchByID(ctx, conn, c.TableName, fields, entryID)
}

func (s *Service) Create(ctx context.Context, requesterID, orgSlug, collectionID string, data map[string]any) (map[string]any, error) {
	c, conn, fields, err := s.getContext(ctx, requesterID, orgSlug, collectionID)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query, args, err := BuildInsert(c.TableName, fields, data)
	if err != nil {
		return nil, err
	}

	var id string
	if err := conn.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return nil, err
	}
	return fetchByID(ctx, conn, c.TableName, fields, id)
}

func (s *Service) Update(ctx context.Context, requesterID, orgSlug, collectionID, entryID string, data map[string]any) (map[string]any, error) {
	c, conn, fields, err := s.getContext(ctx, requesterID, orgSlug, collectionID)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query, args := BuildUpdate(c.TableName, fields, data, entryID)
	if res, err := conn.Exec(ctx, query, args...); err != nil || res.RowsAffected() == 0 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(constants.ErrEntryNotFound)
	}
	return fetchByID(ctx, conn, c.TableName, fields, entryID)
}

func (s *Service) Delete(ctx context.Context, requesterID, orgSlug, collectionID, entryID string) error {
	c, conn, _, err := s.getContext(ctx, requesterID, orgSlug, collectionID)
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := conn.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE id=$1::uuid", quoteIdent(c.TableName)), entryID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New(constants.ErrEntryNotFound)
	}
	return nil
}
