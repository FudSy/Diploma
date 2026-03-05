package handler

import (
	"errors"
	"net/http"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CapacityInput struct {
	Delta int `json:"delta" binding:"required,min=1"`
}

// createResource godoc
// @Summary Create resource
// @Description Creates a new resource. Requires ADMIN role.
// @Tags resources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateResourceRequest true "Resource payload"
// @Success 200 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resources/ [post]
func (h *Handler) createResource(c *gin.Context) {
	var input dto.CreateResourceRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Resource.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, IDResponse{ID: id})
}

// getAllResources godoc
// @Summary List resources
// @Description Returns all resources.
// @Tags resources
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ResourceResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resources/ [get]
func (h *Handler) getAllResources(c *gin.Context) {
	resources, err := h.services.Resource.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resources)
}

// getResourceByID godoc
// @Summary Get resource by ID
// @Description Returns one resource by UUID.
// @Tags resources
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} dto.ResourceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resources/{id} [get]
func (h *Handler) getResourceByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid resource id")
		return
	}

	resource, err := h.services.Resource.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "resource not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resource)
}

// updateResource godoc
// @Summary Update resource
// @Description Updates resource fields. Requires ADMIN role.
// @Tags resources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Param input body dto.UpdateResourceRequest true "Resource update payload"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resources/{id} [put]
func (h *Handler) updateResource(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid resource id")
		return
	}

	var input dto.UpdateResourceRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.services.Resource.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteResource godoc
// @Summary Delete resource
// @Description Deletes resource by UUID. Requires ADMIN role.
// @Tags resources
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resources/{id} [delete]
func (h *Handler) deleteResource(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid resource id")
		return
	}

	if err := h.services.Resource.Delete(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// increaseResourceCapacity godoc
// @Summary Increase resource capacity
// @Description Increases resource capacity by delta. Requires ADMIN role.
// @Tags resources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Param input body CapacityInput true "Capacity delta"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resources/{id}/capacity/increase [patch]
func (h *Handler) increaseResourceCapacity(c *gin.Context) {
	h.changeResourceCapacity(c, true)
}

// decreaseResourceCapacity godoc
// @Summary Decrease resource capacity
// @Description Decreases resource capacity by delta. Requires ADMIN role.
// @Tags resources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Param input body CapacityInput true "Capacity delta"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resources/{id}/capacity/decrease [patch]
func (h *Handler) decreaseResourceCapacity(c *gin.Context) {
	h.changeResourceCapacity(c, false)
}

func (h *Handler) changeResourceCapacity(c *gin.Context, increase bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid resource id")
		return
	}

	var input CapacityInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if increase {
		err = h.services.Resource.IncreaseCapacity(id, input.Delta)
	} else {
		err = h.services.Resource.DecreaseCapacity(id, input.Delta)
	}
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			newErrorResponse(c, http.StatusNotFound, "resource not found")
		case err.Error() == "insufficient resource capacity":
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
