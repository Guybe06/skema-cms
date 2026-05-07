package publicapi

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"skema-api/core/conduit"
	colltypes "skema-api/features/collections/types"
	contentsvc "skema-api/features/content/service"
	pubconstants "skema-api/features/publicapi/constants"
)

type collectionLookup interface {
	FindByOrgAndTable(ctx context.Context, orgID, tableName string) (*colltypes.Collection, error)
	ListFields(ctx context.Context, collectionID string) ([]*colltypes.Field, error)
}

type conduitOpener interface {
	OpenConduit(ctx context.Context, connectionID string) (conduit.Conduit, error)
}

type Handler struct {
	collRepo collectionLookup
	connSvc  conduitOpener
}

func NewHandler(collRepo collectionLookup, connSvc conduitOpener) *Handler {
	return &Handler{collRepo: collRepo, connSvc: connSvc}
}

/*
 * openCollection résout la collection et ouvre la connexion Conduit associée.
 *
 * Attend  : un contexte Gin avec orgIDCtxKey et le paramètre :table.
 * Retourne: la collection, la connexion Conduit, les champs, ou une erreur.
 */

func (h *Handler) openCollection(c *gin.Context) (*colltypes.Collection, conduit.Conduit, []*colltypes.Field, error) {
	orgID := c.GetString(orgIDCtxKey)
	coll, err := h.collRepo.FindByOrgAndTable(c.Request.Context(), orgID, c.Param("table"))
	if err != nil || coll == nil {
		return nil, nil, nil, errors.New(pubconstants.ErrCollNotFound)
	}

	fields, err := h.collRepo.ListFields(c.Request.Context(), coll.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	conn, err := h.connSvc.OpenConduit(c.Request.Context(), coll.ConnectionID)
	if err != nil {
		return nil, nil, nil, errors.New(pubconstants.ErrConnFailed)
	}
	return coll, conn, fields, nil
}

func parsePub(c *gin.Context) (page, perPage int, sort, order string) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ = strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	sort = c.DefaultQuery("sort", "created_at")
	order = c.DefaultQuery("order", "desc")
	if order != "asc" {
		order = "desc"
	}
	return
}

func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

/*
 * fetchEntry récupère une entrée par son identifiant depuis la base client.
 *
 * Attend  : un contexte, une connexion Conduit, le nom de table, les champs et l'ID.
 * Retourne: la ligne sous forme de map, ou une erreur.
 */

func fetchEntry(ctx context.Context, conn conduit.Conduit, table string, fields []*colltypes.Field, id string) (map[string]any, error) {
	selectCols := contentsvc.BuildSelectCols(fields)
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id=$1::uuid", selectCols, quoteIdent(table))
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries, err := contentsvc.ScanRows(rows)
	if err != nil || len(entries) == 0 {
		return nil, errors.New(pubconstants.ErrEntryNotFound)
	}
	return entries[0], nil
}
