package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		segments := api.Group("/segments")
		{
			segments.POST("/", h.createSegment)
			segments.DELETE("/", h.deleteSegment)
			segments.PATCH("/", h.patchSegment)
			segments.GET("/", h.getUsersInSegment)
		}
	}

	return router
}
