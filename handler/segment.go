package handler

import (
	"avito/structures"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createSegment(c *gin.Context) {
	var input structures.Segment

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	slug, err := h.services.Segment.Create(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"slug": slug,
	})
}

func (h *Handler) deleteSegment(c *gin.Context) {

}

func (h *Handler) patchSegment(c *gin.Context) {

}

func (h *Handler) getUsersInSegment(c *gin.Context) {

}
