package publicapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"skema-api/core/response"
	contentsvc "skema-api/features/content/service"
	pubconstants "skema-api/features/publicapi/constants"
)

func (h *Handler) list(c *gin.Context) {
	k := getKey(c)
	if !hasPermission(k, "read") {
		response.Forbidden(c, pubconstants.MsgPermRead)
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
		response.Internal(c, pubconstants.ErrQueryFailed)
		return
	}
	defer rows.Close()

	entries, err := contentsvc.ScanRows(rows)
	if err != nil {
		response.Internal(c, pubconstants.ErrScanFailed)
		return
	}

	var total int
	conn.QueryRow(c.Request.Context(),
		fmt.Sprintf("SELECT COUNT(*) FROM %s", quoteIdent(coll.TableName))).Scan(&total)

	response.List(c, pubconstants.MsgEntriesFound, entries, total, page, perPage)
}

func (h *Handler) get(c *gin.Context) {
	k := getKey(c)
	if !hasPermission(k, "read") {
		response.Forbidden(c, pubconstants.MsgPermRead)
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
		response.NotFound(c, pubconstants.ErrEntryNotFound)
		return
	}
	defer rows.Close()

	entries, _ := contentsvc.ScanRows(rows)
	if len(entries) == 0 {
		response.NotFound(c, pubconstants.ErrEntryNotFound)
		return
	}
	response.OK(c, pubconstants.MsgEntryFound, entries[0])
}
