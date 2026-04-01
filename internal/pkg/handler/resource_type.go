package handler

import (
	"net/http"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getResourceTypes godoc
// @Summary List resource types
// @Description Returns all resource types with their options.
// @Tags resource-types
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ResourceTypeResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resource-types/ [get]
func (h *Handler) getResourceTypes(c *gin.Context) {
	types, err := h.services.ResourceType.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, types)
}

// createResourceType godoc
// @Summary Create resource type
// @Description Creates a new resource type with optional options. Requires ADMIN role.
// @Tags resource-types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateResourceTypeRequest true "Resource type payload"
// @Success 200 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resource-types/ [post]
func (h *Handler) createResourceType(c *gin.Context) {
	var input dto.CreateResourceTypeRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	id, err := h.services.ResourceType.Create(input.Name, input.Options)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, IDResponse{ID: id})
}

// deleteResourceType godoc
// @Summary Delete resource type
// @Description Deletes resource type by UUID. Requires ADMIN role.
// @Tags resource-types
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource type ID (UUID)"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resource-types/{id} [delete]
func (h *Handler) deleteResourceType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор типа ресурса")
		return
	}

	if err := h.services.ResourceType.Delete(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// addResourceTypeOption godoc
// @Summary Add option to resource type
// @Description Adds a new option to an existing resource type. Requires ADMIN role.
// @Tags resource-types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource type ID (UUID)"
// @Param input body dto.ResourceTypeOptionRequest true "Option payload"
// @Success 200 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resource-types/{id}/options [post]
func (h *Handler) addResourceTypeOption(c *gin.Context) {
	rtID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор типа ресурса")
		return
	}

	var input dto.ResourceTypeOptionRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректное тело запроса")
		return
	}

	id, err := h.services.ResourceType.AddOption(rtID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, IDResponse{ID: id})
}

// deleteResourceTypeOption godoc
// @Summary Delete option from resource type
// @Description Deletes an option by UUID. Requires ADMIN role.
// @Tags resource-types
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource type ID (UUID)"
// @Param optionId path string true "Option ID (UUID)"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resource-types/{id}/options/{optionId} [delete]
func (h *Handler) deleteResourceTypeOption(c *gin.Context) {
	_, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор типа ресурса")
		return
	}

	optionID, err := uuid.Parse(c.Param("optionId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор опции")
		return
	}

	if err := h.services.ResourceType.DeleteOption(optionID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
