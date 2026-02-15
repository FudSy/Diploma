package handler

import (
	"errors"
	"net/http"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h *Handler) createBooking(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input dto.CreateBookingRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Booking.Create(userID, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

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
		newErrorResponse(c, http.StatusBadRequest, "invalid booking id")
		return
	}

	booking, err := h.services.Booking.GetById(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "booking not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if isAdmin != "ADMIN" && booking.UserID != userID {
		newErrorResponse(c, http.StatusForbidden, "forbidden")
		return
	}

	c.JSON(http.StatusOK, booking)
}

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
		newErrorResponse(c, http.StatusBadRequest, "invalid booking id")
		return
	}

	booking, err := h.services.Booking.GetById(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "booking not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if isAdmin != "ADMIN" && booking.UserID != userID {
		newErrorResponse(c, http.StatusForbidden, "forbidden")
		return
	}

	var input dto.UpdateBookingRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.services.Booking.Update(userID, bookingID, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

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
		newErrorResponse(c, http.StatusBadRequest, "invalid booking id")
		return
	}

	booking, err := h.services.Booking.GetById(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "booking not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if isAdmin != "ADMIN" && booking.UserID != userID {
		newErrorResponse(c, http.StatusForbidden, "forbidden")
		return
	}

	if err := h.services.Booking.Delete(userID, bookingID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
