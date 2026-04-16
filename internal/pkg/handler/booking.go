package handler

import (
	"errors"
	"net/http"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// getAllBookingsAdmin godoc
// @Summary List all bookings (admin)
// @Description Returns all bookings with user and resource details. Requires ADMIN role.
// @Tags bookings
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.AdminBookingResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings/all [get]
func (h *Handler) getAllBookingsAdmin(c *gin.Context) {
	bookings, err := h.services.Booking.GetAllAdmin()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// createBooking godoc
// @Summary Create booking
// @Description Creates a booking for current user.
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateBookingRequest true "Booking payload"
// @Success 200 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /bookings/ [post]
func (h *Handler) createBooking(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input dto.CreateBookingRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	id, err := h.services.Booking.Create(userID, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, IDResponse{ID: id})
}

// getMyBookings godoc
// @Summary List my bookings
// @Description Returns bookings of current user.
// @Tags bookings
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.BookingResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings/ [get]
func (h *Handler) getMyBookings(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bookings, err := h.services.Booking.GetAll(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// getBookingByID godoc
// @Summary Get booking by ID
// @Description Returns booking by ID. Available for owner or admin.
// @Tags bookings
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID (UUID)"
// @Success 200 {object} dto.BookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings/{id} [get]
func (h *Handler) getBookingByID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	role, _ := c.Get(roleCtx)
	isAdmin, _ := role.(string)

	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор бронирования")
		return
	}

	booking, err := h.services.Booking.GetById(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "бронирование не найдено")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if isAdmin != "ADMIN" && booking.UserID != userID {
		newErrorResponse(c, http.StatusForbidden, "доступ запрещён")
		return
	}

	c.JSON(http.StatusOK, booking)
}

// updateBooking godoc
// @Summary Update booking
// @Description Updates booking by ID. Available for owner or admin.
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID (UUID)"
// @Param input body dto.UpdateBookingRequest true "Booking update payload"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings/{id} [put]
func (h *Handler) updateBooking(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	role, _ := c.Get(roleCtx)
	isAdmin, _ := role.(string)

	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор бронирования")
		return
	}

	booking, err := h.services.Booking.GetById(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "бронирование не найдено")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if isAdmin != "ADMIN" && booking.UserID != userID {
		newErrorResponse(c, http.StatusForbidden, "доступ запрещён")
		return
	}

	var input dto.UpdateBookingRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	if err := h.services.Booking.Update(userID, bookingID, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteBooking godoc
// @Summary Delete booking
// @Description Deletes booking by ID. Available for owner or admin.
// @Tags bookings
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID (UUID)"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings/{id} [delete]
func (h *Handler) deleteBooking(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	role, _ := c.Get(roleCtx)
	isAdmin, _ := role.(string)

	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор бронирования")
		return
	}

	booking, err := h.services.Booking.GetById(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "бронирование не найдено")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if isAdmin != "ADMIN" && booking.UserID != userID {
		newErrorResponse(c, http.StatusForbidden, "доступ запрещён")
		return
	}

	if err := h.services.Booking.Delete(userID, bookingID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
