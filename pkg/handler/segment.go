package handler

import (
	"avito/pkg/structures"
	"avito/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createSegment(c *gin.Context) {
	var input structures.Segment

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.ValidateSlug(input.Slug); err != nil {
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
	var input structures.Segment

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.ValidateSlug(input.Slug); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	slug, err := h.services.Segment.Delete(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"slug": slug,
	})
}
