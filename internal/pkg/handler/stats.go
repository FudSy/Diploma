package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getStatsOverview godoc
// @Summary Get analytics overview
// @Description Returns booking and resource utilization statistics. Requires ADMIN role.
// @Tags stats
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.StatsOverview
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /stats/overview [get]
func (h *Handler) getStatsOverview(c *gin.Context) {
	stats, err := h.services.Analytics.GetOverview()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, stats)
}

// getResourceAvailability godoc
// @Summary Get busy time slots for a resource
// @Description Returns all booked (non-cancelled) time slots for a resource on a given date.
// @Tags resources
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Param date query string true "Date in YYYY-MM-DD format"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resources/{id}/availability [get]
func (h *Handler) getResourceAvailability(c *gin.Context) {
	resourceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор ресурса")
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "дата должна быть в формате YYYY-MM-DD")
		return
	}

	slots, err := h.services.Booking.GetBusySlots(resourceID, dateStr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"date": dateStr, "busy_slots": slots})
}
