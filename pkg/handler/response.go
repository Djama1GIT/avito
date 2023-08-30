package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

type validGetUserHistoryResponse struct {
	Report string `json:"report" example:"http://localhost:8000/files/reports/user_history_YYYY-MM_0.csv"`
	UserId int    `json:"user_id"`
}

type validGetUserSegmentsResponse struct {
	Segments []string `json:"segments"`
	UserId   int      `json:"user_id"`
}

type validPatchResponse struct {
	UserId int `json:"user_id"`
}

type validCreateSegmentResponse struct {
	Segment string `json:"slug"`
}

type validDeleteSegmentResponse struct {
	Segment string `json:"slug"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Print(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
