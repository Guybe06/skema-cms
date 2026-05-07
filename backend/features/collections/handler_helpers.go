package collections

import (
	"github.com/gin-gonic/gin"
	"skema-api/core/response"
	"skema-api/features/collections/constants"
	"skema-api/features/collections/types"
)

func toCollectionResp(c *types.Collection) types.CollectionResponse {
	r := types.CollectionResponse{
		ID: c.ID, ConnectionID: c.ConnectionID, OrganizationID: c.OrganizationID,
		Name: c.Name, TableName: c.TableName, DisplayName: c.DisplayName,
		Description: c.Description, CreatedAt: c.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
	for _, f := range c.Fields {
		r.Fields = append(r.Fields, toFieldResp(f))
	}
	return r
}

func toFieldResp(f *types.Field) *types.FieldResponse {
	return &types.FieldResponse{
		ID: f.ID, Name: f.Name, ColumnName: f.ColumnName, Type: f.Type,
		Required: f.Required, IsUnique: f.IsUnique,
		DefaultValue: f.DefaultValue, Options: f.Options, Position: f.Position,
	}
}

func handleErr(c *gin.Context, err error) {
	switch err.Error() {
	case constants.ErrNotAuthorized:
		response.Forbidden(c, err.Error())
	case constants.ErrCollectionNotFound, constants.ErrFieldNotFound:
		response.NotFound(c, err.Error())
	case constants.ErrTableNameTaken, constants.ErrColumnNameTaken, constants.ErrSchemaFailed:
		response.BadRequest(c, err.Error())
	default:
		response.Internal(c, constants.MsgInternalError)
	}
}
