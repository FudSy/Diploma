package handler

import (
	"errors"
	"net/http"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const userCtx = "userID"
const roleCtx = "role"
const authCookieName = "access_token"
const authCookieMaxAge = 12 * 60 * 60 // 12h

// register godoc
// @Summary Register user
// @Description Creates a regular user account.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.RegisterRequest true "Registration payload"
// @Success 200 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
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

	c.JSON(http.StatusOK, IDResponse{ID: id})
}

// registerAdmin godoc
// @Summary Register admin
// @Description Creates an admin account. Requires ADMIN role.
// @Tags auth-admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.RegisterRequest true "Registration payload"
// @Success 200 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/admin/register [post]
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

	c.JSON(http.StatusOK, IDResponse{ID: id})
}

// login godoc
// @Summary Login
// @Description Authenticates user and returns access token.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LoginRequest true "Login payload"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
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

	c.JSON(http.StatusOK, TokenResponse{Token: token})
}

// me godoc
// @Summary Get current user
// @Description Returns current authorized user profile.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MeResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/me [get]
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
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Surname: user.Surname,
		Role:    user.Role,
	})
}

// logout godoc
// @Summary Logout
// @Description Clears auth cookie.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StatusResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	c.SetCookie(authCookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// adminCheck godoc
// @Summary Admin access check
// @Description Checks that current user has admin access.
// @Tags auth-admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StatusResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /auth/admin/check [get]
func (h *Handler) adminCheck(c *gin.Context) {
	c.JSON(http.StatusOK, StatusResponse{Status: "admin access granted"})
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
