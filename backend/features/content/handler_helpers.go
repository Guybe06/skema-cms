package content

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"skema-api/core/response"
	"skema-api/features/content/constants"
)

func parsePagination(c *gin.Context) (page, perPage int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ = strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return
}

func handleErr(c *gin.Context, err error) {
	switch err.Error() {
	case constants.ErrNotAuthorized:
		response.Forbidden(c, err.Error())
	case constants.ErrEntryNotFound:
		response.NotFound(c, err.Error())
	default:
		response.Internal(c, constants.MsgInternalError)
	}
}
