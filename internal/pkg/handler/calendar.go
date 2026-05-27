package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// downloadBookingICS godoc
// @Summary Download booking as .ics
// @Description Returns a single booking as RFC 5545 iCalendar file. Owner or admin only.
// @Tags calendar
// @Produce text/calendar
// @Security BearerAuth
// @Param id path string true "Booking ID (UUID)"
// @Success 200 {string} string "iCalendar (RFC 5545) document"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /bookings/{id}/ics [get]
func (h *Handler) downloadBookingICS(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	role, _ := c.Get(roleCtx)
	isAdmin := role == "ADMIN"

	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор бронирования")
		return
	}

	booking, err := h.services.Calendar.GetBookingForUser(bookingID, userID, isAdmin)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "бронирование не найдено")
			return
		}
		if err.Error() == "доступ запрещён" {
			newErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	body := h.services.Calendar.BuildICS([]dto.CalendarBooking{booking}, "Booking")
	writeICSResponse(c, "booking-"+bookingID.String()+".ics", body)
}

// googleCalendarLink godoc
// @Summary Get "Add to Google Calendar" link
// @Description Returns a Google Calendar deep-link to add this booking. Owner or admin only.
// @Tags calendar
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID (UUID)"
// @Success 200 {object} dto.GoogleLinkResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /bookings/{id}/google-link [get]
func (h *Handler) googleCalendarLink(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	role, _ := c.Get(roleCtx)
	isAdmin := role == "ADMIN"

	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор бронирования")
		return
	}

	booking, err := h.services.Calendar.GetBookingForUser(bookingID, userID, isAdmin)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "бронирование не найдено")
			return
		}
		if err.Error() == "доступ запрещён" {
			newErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.GoogleLinkResponse{URL: h.services.Calendar.BuildGoogleLink(booking)})
}

// getMyCalendarFeed godoc
// @Summary Get personal calendar feed URLs
// @Description Returns subscription URLs (HTTPS + webcal://) for the user's bookings feed. Creates a token on first call.
// @Tags calendar
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CalendarFeedResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /me/calendar [get]
func (h *Handler) getMyCalendarFeed(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := h.services.Calendar.EnsureToken(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, buildFeedResponse(c, token))
}

// rotateMyCalendarFeed godoc
// @Summary Rotate calendar feed token
// @Description Issues a fresh feed token, invalidating any previous subscription URL.
// @Tags calendar
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CalendarFeedResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /me/calendar/rotate [post]
func (h *Handler) rotateMyCalendarFeed(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := h.services.Calendar.RotateToken(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, buildFeedResponse(c, token))
}

// publicCalendarFeed godoc
// @Summary Public calendar feed (token-based, no auth)
// @Description Returns the user's bookings as an iCalendar feed. URL is opaque — for use in calendar apps (Google Calendar, Apple Calendar, Outlook).
// @Tags calendar
// @Produce text/calendar
// @Param token path string true "Calendar feed token"
// @Success 200 {string} string "iCalendar (RFC 5545) document"
// @Failure 404 {object} ErrorResponse
// @Router /calendar/feed/{token} [get]
func (h *Handler) publicCalendarFeed(c *gin.Context) {
	token := strings.TrimSuffix(c.Param("token"), ".ics")
	bookings, user, err := h.services.Calendar.BookingsByToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "календарный фид не найден")
		return
	}

	calName := "My Bookings"
	if user.Name != "" {
		calName = fmt.Sprintf("%s — Bookings", user.Name)
	}

	body := h.services.Calendar.BuildICS(bookings, calName)
	writeICSResponse(c, "bookings.ics", body)
}

func buildFeedResponse(c *gin.Context, token string) dto.CalendarFeedResponse {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.Request.Host
	feedURL := fmt.Sprintf("%s://%s/calendar/feed/%s.ics", scheme, host, token)
	webcalURL := "webcal://" + host + "/calendar/feed/" + token + ".ics"
	return dto.CalendarFeedResponse{
		Token:     token,
		FeedURL:   feedURL,
		WebcalURL: webcalURL,
	}
}

func writeICSResponse(c *gin.Context, filename, body string) {
	c.Header("Content-Type", "text/calendar; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Header("Cache-Control", "no-cache, must-revalidate")
	c.String(http.StatusOK, body)
}
