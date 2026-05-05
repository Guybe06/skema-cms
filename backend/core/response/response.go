package response

import "github.com/gin-gonic/gin"

type Body struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ListBody struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

func OK(c *gin.Context, message string, data any) {
	c.JSON(StatusOK, Body{Success: true, Message: message, Data: data})
}

func Created(c *gin.Context, message string, data any) {
	c.JSON(StatusCreated, Body{Success: true, Message: message, Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(StatusNoContent)
}

/*
 * List envoie une réponse paginée avec les métadonnées de pagination.
 *
 * Attend  : les données, le total d'éléments, la page courante et la taille de page.
 * Retourne: une réponse JSON avec data et pagination.
 */

func List(c *gin.Context, message string, data any, total, page, perPage int) {
	totalPages := total / perPage
	if total%perPage != 0 {
		totalPages++
	}

	c.JSON(StatusOK, ListBody{
		Success: true,
		Message: message,
		Data:    data,
		Pagination: PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(StatusBadRequest, ErrorBody{Success: false, Message: message})
}

func ValidationError(c *gin.Context, message string, errors any) {
	c.JSON(StatusBadRequest, ErrorBody{Success: false, Message: message, Errors: errors})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(StatusUnauthorized, ErrorBody{Success: false, Message: message})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(StatusForbidden, ErrorBody{Success: false, Message: message})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(StatusNotFound, ErrorBody{Success: false, Message: message})
}

func Conflict(c *gin.Context, message string) {
	c.JSON(StatusConflict, ErrorBody{Success: false, Message: message})
}

func Internal(c *gin.Context, message string) {
	c.JSON(StatusInternal, ErrorBody{Success: false, Message: message})
}

func TooManyRequests(c *gin.Context) {
	c.JSON(StatusTooManyRequests, ErrorBody{Success: false, Message: MsgTooManyRequests})
}

func RequestTooLarge(c *gin.Context) {
	c.JSON(StatusBadRequest, ErrorBody{Success: false, Message: MsgRequestTooLarge})
}

func RequestTimeout(c *gin.Context) {
	c.JSON(StatusBadRequest, ErrorBody{Success: false, Message: MsgRequestTimeout})
}
