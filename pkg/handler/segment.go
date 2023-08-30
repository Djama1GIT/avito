package handler

import (
	"avito/pkg/structures"
	"avito/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create Segment
// @Tags segment
// @ID create-segment
// @Accept  json
// @Produce  json
// @Param input body structures.Segment true "Slug of segment"
// @Success 200 {object} validCreateSegmentResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /segments/ [post]
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

	c.JSON(http.StatusOK, validCreateSegmentResponse{
		Segment: slug,
	})
}

// @Summary Delete Segment
// @Tags segment
// @ID delete-segment
// @Accept  json
// @Produce  json
// @Param input body structures.Segment true "Slug of segment"
// @Success 200 {object} validDeleteSegmentResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /segments/ [delete]
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

	c.JSON(http.StatusOK, validDeleteSegmentResponse{
		Segment: slug,
	})
}
