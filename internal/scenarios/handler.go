package scenarios

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	items, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid scenario id",
		})
		return
	}

	scenario, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrScenarioNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "scenario not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, scenario)
}
