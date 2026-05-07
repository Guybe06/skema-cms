package publicapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"skema-api/core/response"
	contentsvc "skema-api/features/content/service"
	pubconstants "skema-api/features/publicapi/constants"
)

func (h *Handler) create(c *gin.Context) {
	k := getKey(c)
	if !hasPermission(k, "create") {
		response.Forbidden(c, pubconstants.MsgPermCreate)
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
		response.BadRequest(c, pubconstants.ErrInvalidJSON)
		return
	}
	query, args, err := contentsvc.BuildInsert(coll.TableName, fields, data)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	var id string
	if err := conn.QueryRow(c.Request.Context(), query, args...).Scan(&id); err != nil {
		response.Internal(c, pubconstants.ErrCreateFailed)
		return
	}
	entry, err := fetchEntry(c.Request.Context(), conn, coll.TableName, fields, id)
	if err != nil {
		response.Internal(c, pubconstants.ErrFetchFailed)
		return
	}
	response.Created(c, pubconstants.MsgEntryCreated, entry)
}

func (h *Handler) update(c *gin.Context) {
	k := getKey(c)
	if !hasPermission(k, "update") {
		response.Forbidden(c, pubconstants.MsgPermUpdate)
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
		response.BadRequest(c, pubconstants.ErrInvalidJSON)
		return
	}
	query, args := contentsvc.BuildUpdate(coll.TableName, fields, data, c.Param("id"))
	if res, err := conn.Exec(c.Request.Context(), query, args...); err != nil || res.RowsAffected() == 0 {
		response.NotFound(c, pubconstants.ErrEntryNotFound)
		return
	}
	entry, _ := fetchEntry(c.Request.Context(), conn, coll.TableName, fields, c.Param("id"))
	response.OK(c, pubconstants.MsgEntryUpdated, entry)
}

func (h *Handler) delete(c *gin.Context) {
	k := getKey(c)
	if !hasPermission(k, "delete") {
		response.Forbidden(c, pubconstants.MsgPermDelete)
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
		response.NotFound(c, pubconstants.ErrEntryNotFound)
		return
	}
	response.NoContent(c)
}
