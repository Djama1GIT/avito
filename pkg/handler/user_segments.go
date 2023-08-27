package handler

import (
	"avito/pkg/structures"
	"avito/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) patchSegment(c *gin.Context) {
	var input structures.UserSegments

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	for _, segment := range input.SegmentsToAdd {
		if err := utils.ValidateSlug(segment); err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	for _, segment := range input.SegmentsToDelete {
		if err := utils.ValidateSlug(segment); err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	user_id, err := h.services.UserSegments.Patch(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": user_id,
	})
}

func (h *Handler) getUsersInSegment(c *gin.Context) {
	var input structures.User

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	segments, err := h.services.UserSegments.GetUsersInSegment(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if segments == nil {
		segments = []string{}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"segments": segments,
	})
}
