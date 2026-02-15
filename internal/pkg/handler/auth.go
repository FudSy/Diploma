package handler

import (
	"errors"
	"net/http"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

const userCtx = "userID"
const roleCtx = "role"
const authCookieName = "access_token"
const authCookieMaxAge = 12 * 60 * 60 // 12h

func (h *Handler) register(c *gin.Context) {
	var input dto.RegisterRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) registerAdmin(c *gin.Context) {
	var input dto.RegisterRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Authorization.CreateAdmin(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}


func (h *Handler) login(c *gin.Context) {
	var input dto.LoginRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.Login(input.Login, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Also persist token in cookie so clients can use auth without manually setting Authorization header.
	c.SetCookie(authCookieName, token, authCookieMaxAge, "/", "", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) me(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.services.Authorization.GetUserById(userID)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "user not found")
		return
	}

	c.JSON(http.StatusOK, dto.MeResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}

func (h *Handler) logout(c *gin.Context) {
	c.SetCookie(authCookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

func (h *Handler) adminCheck(c *gin.Context) {
	c.JSON(http.StatusOK, statusResponse{Status: "admin access granted"})
}

func getUserID(c *gin.Context) (uuid.UUID, error) {
	userID, ok := c.Get(userCtx)
	if !ok {
		return uuid.Nil, errors.New("user id not found in context")
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user id has invalid type")
	}

	return id, nil
}
