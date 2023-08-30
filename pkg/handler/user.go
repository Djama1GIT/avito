package handler

import (
	"avito/pkg/structures"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get User History
// @Description You can also use the request body to send data, but not here :)
// @Description p.s. For example, via curl
// @Tags user
// @ID get-user-history
// @Accpet json
// @Produce json
// @Param input query structures.UserHistory true "User History Data"
// @Success 200 {object} validGetUserHistoryResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/history/ [get]
func (h *Handler) getUserHistory(c *gin.Context) {
	var input structures.UserHistory
	var err error

	userId := c.Query("user_id")
	if userId != "" {
		input.Id, err = strconv.Atoi(userId)
		input.YearMonth = c.Query("year_month")
		if input.YearMonth == "" {
			NewErrorResponse(c, http.StatusBadRequest, "Key: 'UserHistory.YearMonth' Error:Field validation for 'YearMonth' failed on the 'required' tag")
			return
		}
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

	report, err := h.services.User.GetUserHistory(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, validGetUserHistoryResponse{
		UserId: input.Id,
		Report: "http://localhost:8000/files/" + report,
	})
}

// @Summary Delete Expired User Segments
// @Tags user
// @ID delete-user-expired-segments
// @Accpet json
// @Produce json
// @Success 200 integer 1
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/expired-segments/ [delete]
func (h *Handler) deleteExpiredSegments(c *gin.Context) {
	log.Print("Garbage collection request received")
	if err := h.services.User.DeleteExpiredSegments(); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
}
