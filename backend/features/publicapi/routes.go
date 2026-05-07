package publicapi

import "github.com/gin-gonic/gin"

func RegisterRoutes(pub *gin.RouterGroup, ks keySvc, os orgSvc, collRepo collectionLookup, connSvc conduitOpener) {
	h := NewHandler(collRepo, connSvc)
	g := pub.Group("/:orgSlug/:table", APIKeyAuth(ks, os))
	{
		g.GET("", h.list)
		g.POST("", h.create)
		g.GET("/:id", h.get)
		g.PATCH("/:id", h.update)
		g.DELETE("/:id", h.delete)
	}
}
