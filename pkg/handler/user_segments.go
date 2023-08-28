package handler

import (
	"avito/pkg/structures"
	"avito/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Patch Segment
// @Tags user-segments
// @ID patch-segment
// @Accept  json
// @Produce  json
// @Param input body structures.UserSegments true "Patch data"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /segments/ [patch]
func (h *Handler) patchSegment(c *gin.Context) {
	var input structures.UserSegments

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	for _, segment := range input.SegmentsToAdd {
		if err := utils.ValidateSlug(segment); err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error()+" (segment to add: "+segment+")")
			return
		}
	}

	for _, segment := range input.SegmentsToDelete {
		if err := utils.ValidateSlug(segment); err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error()+" (segment to delete: "+segment+")")
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

// @Summary Get Users In Segment
// @Tags user-segments
// @ID get-users-in-segment
// @Accpet json
// @Produce json
// @Param user_id query integer true "User data"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /segments/ [get]
func (h *Handler) getUsersInSegment(c *gin.Context) {
	var input structures.User
	var err error

	userId := c.Query("user_id")
	if userId != "" {
		input.Id, err = strconv.Atoi(userId)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err := c.BindJSON(&input); err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
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
		"user_id":  input.Id,
		"segments": segments,
	})
}
