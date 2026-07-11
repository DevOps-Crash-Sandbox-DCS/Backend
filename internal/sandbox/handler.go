package sandbox

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	manager *Manager
}

func NewHandler(manager *Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

type StartRequest struct {
	SessionID  string `json:"sessionId" binding:"required"`
	ScenarioID string `json:"scenarioId" binding:"required"`
	Image      string `json:"image" binding:"required"`
}

type ExecRequest struct {
	Command string `json:"command" binding:"required"`
}

type CleanupRequest struct {
	OlderThanMinutes int `json:"olderThanMinutes"`
}

func (h *Handler) Start(c *gin.Context) {
	var req StartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	sessionID, err := uuid.Parse(strings.TrimSpace(req.SessionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid session id",
		})
		return
	}

	sandbox, err := h.manager.Start(c.Request.Context(), StartSandboxRequest{
		SessionID:  sessionID,
		ScenarioID: req.ScenarioID,
		Image:      req.Image,
	})
	if err != nil {
		if errors.Is(err, ErrSandboxAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "sandbox already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, toSandboxResponse(sandbox))
}

func (h *Handler) Get(c *gin.Context) {
	sessionID, ok := parseSessionIDParam(c)
	if !ok {
		return
	}

	sandbox, err := h.manager.Get(c.Request.Context(), sessionID)
	if err != nil {
		if errors.Is(err, ErrSandboxNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "sandbox not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, toSandboxResponse(sandbox))
}

func (h *Handler) Exec(c *gin.Context) {
	sessionID, ok := parseSessionIDParam(c)
	if !ok {
		return
	}

	var req ExecRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	result, err := h.manager.Exec(c.Request.Context(), ExecCommandRequest{
		SessionID: sessionID,
		Command:   req.Command,
	})
	if err != nil {
		if errors.Is(err, ErrSandboxNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "sandbox not found",
			})
			return
		}

		if errors.Is(err, ErrInvalidCommand) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid command",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stdout":   result.Stdout,
		"stderr":   result.Stderr,
		"exitCode": result.ExitCode,
	})
}

func (h *Handler) Stop(c *gin.Context) {
	sessionID, ok := parseSessionIDParam(c)
	if !ok {
		return
	}

	err := h.manager.Stop(c.Request.Context(), sessionID)
	if err != nil {
		if errors.Is(err, ErrSandboxNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "sandbox not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "stopped",
	})
}

func (h *Handler) Cleanup(c *gin.Context) {
	var req CleanupRequest

	_ = c.ShouldBindJSON(&req)

	olderThanMinutes := req.OlderThanMinutes
	if olderThanMinutes <= 0 {
		olderThanMinutes = 60
	}

	result, err := h.manager.CleanupOldRunning(
		c.Request.Context(),
		time.Duration(olderThanMinutes)*time.Minute,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func parseSessionIDParam(c *gin.Context) (uuid.UUID, bool) {
	sessionIDParam := strings.TrimSpace(c.Param("sessionId"))

	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid session id",
		})
		return uuid.Nil, false
	}

	return sessionID, true
}

func toSandboxResponse(sandbox *Sandbox) gin.H {
	return gin.H{
		"id":            sandbox.ID,
		"sessionId":     sandbox.SessionID,
		"scenarioId":    sandbox.ScenarioID,
		"containerName": sandbox.ContainerName,
		"image":         sandbox.Image,
		"status":        sandbox.Status,
		"startedAt":     sandbox.StartedAt,
		"stoppedAt":     sandbox.StoppedAt,
		"lastSeenAt":    sandbox.LastSeenAt,
		"createdAt":     sandbox.CreatedAt,
		"updatedAt":     sandbox.UpdatedAt,
	}
}
