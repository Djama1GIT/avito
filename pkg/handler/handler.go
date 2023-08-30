package handler

import (
	"avito/pkg/service"

	_ "avito/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.Default())

	files := router.Group("/files")
	{
		files.Static("/reports", "./reports")
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		segments := api.Group("/segments")
		{
			segments.POST("/", h.createSegment)
			segments.DELETE("/", h.deleteSegment)
			segments.PATCH("/", h.patchSegment)
			segments.GET("/", h.getUsersInSegment)
		}

		users := api.Group(("/users"))
		{
			history := users.Group("/history")
			{
				history.GET("/", h.getUserHistory)
			}
		}
	}

	return router
}
