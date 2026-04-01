package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
		newErrorResponse(c, http.StatusBadRequest, "некорректное тело запроса")
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
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор ресурса")
		return
	}

	resource, err := h.services.Resource.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "ресурс не найден")
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
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор ресурса")
		return
	}

	var input dto.UpdateResourceRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректное тело запроса")
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
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор ресурса")
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

// uploadResourcePhoto godoc
// @Summary Upload resource photo
// @Description Uploads a photo for a resource. Requires ADMIN role.
// @Tags resources
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resource ID (UUID)"
// @Param photo formData file true "Photo file (jpeg/png/webp, max 5MB)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /resources/{id}/photo [post]
func (h *Handler) uploadResourcePhoto(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор ресурса")
		return
	}

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "файл не найден в запросе")
		return
	}
	defer file.Close()

	if header.Size > 5<<20 {
		newErrorResponse(c, http.StatusBadRequest, "размер файла не должен превышать 5MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		newErrorResponse(c, http.StatusBadRequest, "допустимые форматы: jpg, jpeg, png, webp")
		return
	}

	uploadDir := "./uploads/resources"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "ошибка создания директории для загрузки")
		return
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(header, savePath); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "ошибка сохранения файла")
		return
	}

	photoURL := fmt.Sprintf("/uploads/resources/%s", filename)
	if err := h.services.Resource.UpdatePhoto(id, photoURL); err != nil {
		os.Remove(savePath)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "ресурс не найден")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"photo_url": photoURL})
}

func (h *Handler) changeResourceCapacity(c *gin.Context, increase bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректный идентификатор ресурса")
		return
	}

	var input CapacityInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "некорректное тело запроса")
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
			newErrorResponse(c, http.StatusNotFound, "ресурс не найден")
		case err.Error() == "недостаточная вместимость ресурса":
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
