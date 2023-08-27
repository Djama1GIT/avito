package handler

import (
	"avito/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
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
