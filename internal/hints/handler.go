package hints

import (
	"errors"
	"net/http"
	"strings"

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

func (h *Handler) Create(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	sessionIDParam := strings.TrimSpace(c.Param("id"))

	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid session id",
		})
		return
	}

	var req CreateHintRequest

	_ = c.ShouldBindJSON(&req)

	result, err := h.service.CreateHint(
		c.Request.Context(),
		sessionID,
		userID,
		req,
	)
	if err != nil {
		if errors.Is(err, ErrHintSessionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "session not found",
			})
			return
		}

		if errors.Is(err, ErrInvalidHintLevel) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid hint level",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user is not authorized",
		})
		return uuid.Nil, false
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user context",
		})
		return uuid.Nil, false
	}

	return userID, true
}
