package actions

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

func (h *Handler) Submit(c *gin.Context) {
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

	var req SubmitActionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	res, err := h.service.Submit(c.Request.Context(), sessionID, userID, req)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "session not found",
			})
			return
		}

		if errors.Is(err, ErrStepNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "step not found",
			})
			return
		}

		if errors.Is(err, ErrInvalidStep) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid step for current session",
			})
			return
		}

		if errors.Is(err, ErrInvalidCommand) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid command",
			})
			return
		}

		if errors.Is(err, ErrInvalidSessionState) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "session is not in progress",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
