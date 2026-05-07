package publicapi

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"skema-api/core/conduit"
	"skema-api/core/response"
	colltypes "skema-api/features/collections/types"
	contentsvc "skema-api/features/content/service"
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

func (h *Handler) openCollection(c *gin.Context) (*colltypes.Collection, conduit.Conduit, []*colltypes.Field, error) {
	orgID := c.GetString(orgIDCtxKey)
	coll, err := h.collRepo.FindByOrgAndTable(c.Request.Context(), orgID, c.Param("table"))
	if err != nil || coll == nil {
		return nil, nil, nil, errors.New("collection introuvable")
	}

	fields, err := h.collRepo.ListFields(c.Request.Context(), coll.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	conn, err := h.connSvc.OpenConduit(c.Request.Context(), coll.ConnectionID)
	if err != nil {
		return nil, nil, nil, errors.New("connexion impossible")
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

func (h *Handler) list(c *gin.Context) {
	k := getKey(c)
	if !hasPermission(k, "read") {
		response.Forbidden(c, "Permission lecture manquante.")
		return
	}

	coll, conn, fields, err := h.openCollection(c)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	defer conn.Close()

	page, perPage, sort, order := parsePub(c)
	selectCols := contentsvc.BuildSelectCols(fields)
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s %s LIMIT $1 OFFSET $2",
		selectCols, quoteIdent(coll.TableName), quoteIdent(sort), order)

	rows, err := conn.Query(c.Request.Context(), query, perPage, (page-1)*perPage)
	if err != nil {
		response.Internal(c, "Erreur lors de la requête.")
		return
	}
	defer rows.Close()

	entries, err := contentsvc.ScanRows(rows)
	if err != nil {
		response.Internal(c, "Erreur lors du scan.")
		return
	}

	var total int
	conn.QueryRow(c.Request.Context(),
		fmt.Sprintf("SELECT COUNT(*) FROM %s", quoteIdent(coll.TableName))).Scan(&total)

	response.List(c, "Entrées récupérées.", entries, total, page, perPage)
}

func (h *Handler) get(c *gin.Context) {
	k := getKey(c)
	if !hasPermission(k, "read") {
		response.Forbidden(c, "Permission lecture manquante.")
		return
	}

	coll, conn, fields, err := h.openCollection(c)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	defer conn.Close()

	selectCols := contentsvc.BuildSelectCols(fields)
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id=$1::uuid", selectCols, quoteIdent(coll.TableName))
	rows, err := conn.Query(c.Request.Context(), query, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Entrée introuvable.")
		return
	}
	defer rows.Close()

	entries, _ := contentsvc.ScanRows(rows)
	if len(entries) == 0 {
		response.NotFound(c, "Entrée introuvable.")
		return
	}
	response.OK(c, "Entrée récupérée.", entries[0])
}

func (h *Handler) create(c *gin.Context) {
	k := getKey(c)
	if !hasPermission(k, "create") {
		response.Forbidden(c, "Permission création manquante.")
		return
	}

	coll, conn, fields, err := h.openCollection(c)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	defer conn.Close()

	var data map[string]any
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "Corps JSON invalide.")
		return
	}

	query, args, err := contentsvc.BuildInsert(coll.TableName, fields, data)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var id string
	if err := conn.QueryRow(c.Request.Context(), query, args...).Scan(&id); err != nil {
		response.Internal(c, "Erreur lors de la création.")
		return
	}

	entry, err := fetchEntry(c.Request.Context(), conn, coll.TableName, fields, id)
	if err != nil {
		response.Internal(c, "Erreur lors de la récupération.")
		return
	}
	response.Created(c, "Entrée créée.", entry)
}

func (h *Handler) update(c *gin.Context) {
	k := getKey(c)
	if !hasPermission(k, "update") {
		response.Forbidden(c, "Permission modification manquante.")
		return
	}

	coll, conn, fields, err := h.openCollection(c)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	defer conn.Close()

	var data map[string]any
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "Corps JSON invalide.")
		return
	}

	query, args := contentsvc.BuildUpdate(coll.TableName, fields, data, c.Param("id"))
	if res, err := conn.Exec(c.Request.Context(), query, args...); err != nil || res.RowsAffected() == 0 {
		response.NotFound(c, "Entrée introuvable.")
		return
	}

	entry, _ := fetchEntry(c.Request.Context(), conn, coll.TableName, fields, c.Param("id"))
	response.OK(c, "Entrée modifiée.", entry)
}

func (h *Handler) delete(c *gin.Context) {
	k := getKey(c)
	if !hasPermission(k, "delete") {
		response.Forbidden(c, "Permission suppression manquante.")
		return
	}

	coll, conn, _, err := h.openCollection(c)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	defer conn.Close()

	res, err := conn.Exec(c.Request.Context(),
		fmt.Sprintf("DELETE FROM %s WHERE id=$1::uuid", quoteIdent(coll.TableName)),
		c.Param("id"))
	if err != nil || res.RowsAffected() == 0 {
		response.NotFound(c, "Entrée introuvable.")
		return
	}
	response.NoContent(c)
}

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
		return nil, errors.New("entrée introuvable")
	}
	return entries[0], nil
}

func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}
