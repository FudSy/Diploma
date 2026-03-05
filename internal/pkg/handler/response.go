package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type IDResponse struct {
	ID uuid.UUID `json:"id"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Error().Msg(message)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}
