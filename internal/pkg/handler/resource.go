package handler

import (
	"errors"
	"net/http"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type capacityInput struct {
	Delta int `json:"delta" binding:"required,min=1"`
}

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

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllResources(c *gin.Context) {
	resources, err := h.services.Resource.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resources)
}

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

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

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

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

func (h *Handler) increaseResourceCapacity(c *gin.Context) {
	h.changeResourceCapacity(c, true)
}

func (h *Handler) decreaseResourceCapacity(c *gin.Context) {
	h.changeResourceCapacity(c, false)
}

func (h *Handler) changeResourceCapacity(c *gin.Context, increase bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid resource id")
		return
	}

	var input capacityInput
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

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
