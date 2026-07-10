package reports

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetSessionReport(c *gin.Context) {
	userIDValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user is not authorized",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user context",
		})
		return
	}

	sessionIDParam := c.Param("id")

	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid session id",
		})
		return
	}

	report, err := h.service.GetSessionReport(c.Request.Context(), sessionID, userID)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "session not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, report)
}
