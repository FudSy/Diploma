package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) userIdentity(c *gin.Context) {
	token, err := extractAccessToken(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := h.services.Authorization.ParseToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.services.Authorization.GetUserById(userID)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "пользователь не найден")
		return
	}

	c.Set(userCtx, userID)
	c.Set(roleCtx, user.Role)
	c.Next()
}

func extractAccessToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		headerParts := strings.SplitN(authHeader, " ", 2)
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return "", errors.New("некорректный заголовок авторизации")
		}
		return headerParts[1], nil
	}

	token, err := c.Cookie(authCookieName)
	if err == nil && token != "" {
		return token, nil
	}

	return "", errors.New("токен авторизации не найден")
}

func (h *Handler) adminIdentity(c *gin.Context) {
	role, ok := c.Get(roleCtx)
	if !ok {
		newErrorResponse(c, http.StatusForbidden, "роль не найдена в контексте")
		return
	}

	userRole, ok := role.(string)
	if !ok {
		newErrorResponse(c, http.StatusForbidden, "некорректный тип роли")
		return
	}
	if userRole != "ADMIN" {
		newErrorResponse(c, http.StatusForbidden, "требуются права администратора")
		return
	}

	c.Next()
}
