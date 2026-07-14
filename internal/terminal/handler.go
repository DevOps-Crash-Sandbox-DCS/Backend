package terminal

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"DCS/internal/actions"
	"DCS/internal/sandbox"
	"DCS/internal/sessions"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler struct {
	sessionsService *sessions.Service
	actionsService  *actions.Service
	terminalService *Service
	sandboxManager  *sandbox.Manager
}

func NewHandler(
	sessionsService *sessions.Service,
	actionsService *actions.Service,
	terminalService *Service,
	sandboxManager *sandbox.Manager,
) *Handler {
	return &Handler{
		sessionsService: sessionsService,
		actionsService:  actionsService,
		terminalService: terminalService,
		sandboxManager:  sandboxManager,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) Connect(c *gin.Context) {
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

	session, err := h.sessionsService.GetByID(c.Request.Context(), sessionID, userID)
	if err != nil {
		if errors.Is(err, sessions.ErrSessionNotFound) {
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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	sandboxStarted := false

	defer func() {
		_ = conn.Close()

		if sandboxStarted {
			stopSandboxInBackground(h.sandboxManager, sessionID)
		}
	}()

	welcomeOutput := "Connected to incident simulator terminal.\n"

	if image, ok := dockerImageForScenario(session.ScenarioID); ok && h.sandboxManager != nil {
		_, err := h.sandboxManager.EnsureStarted(
			c.Request.Context(),
			sandbox.StartSandboxRequest{
				SessionID:  sessionID,
				ScenarioID: session.ScenarioID,
				Image:      image,
			},
		)
		if err != nil {
			_ = conn.WriteJSON(ServerMessage{
				Type:  ServerMessageTypeError,
				Error: "failed to start sandbox: " + err.Error(),
			})
			return
		}

		sandboxStarted = true
		welcomeOutput += "Docker sandbox started.\n"
	}

	err = conn.WriteJSON(ServerMessage{
		Type:          ServerMessageTypeWelcome,
		SessionID:     session.ID.String(),
		CurrentStepID: session.CurrentStepID,
		Output:        welcomeOutput,
	})
	if err != nil {
		return
	}

	for {
		var msg ClientMessage

		if err := conn.ReadJSON(&msg); err != nil {
			return
		}

		if msg.Type != ClientMessageTypeCommand {
			_ = conn.WriteJSON(ServerMessage{
				Type:  ServerMessageTypeError,
				Error: "unsupported message type",
			})
			continue
		}

		command := strings.TrimSpace(msg.Command)
		if command == "" {
			_ = conn.WriteJSON(ServerMessage{
				Type:  ServerMessageTypeError,
				Error: "command is empty",
			})
			continue
		}

		currentSession, err := h.sessionsService.GetByID(c.Request.Context(), sessionID, userID)
		if err != nil {
			_ = conn.WriteJSON(ServerMessage{
				Type:  ServerMessageTypeError,
				Error: "failed to refresh session",
			})
			continue
		}

		if currentSession.CurrentStepID == nil {
			_ = conn.WriteJSON(ServerMessage{
				Type:      ServerMessageTypeSessionClosed,
				SessionID: currentSession.ID.String(),
				Output:    "Session has no current step.\n",
			})
			continue
		}

		var commandResult *CommandResult

		if image, ok := dockerImageForScenario(currentSession.ScenarioID); ok && h.sandboxManager != nil {
			_, err := h.sandboxManager.EnsureStarted(
				c.Request.Context(),
				sandbox.StartSandboxRequest{
					SessionID:  sessionID,
					ScenarioID: currentSession.ScenarioID,
					Image:      image,
				},
			)
			if err != nil {
				_ = conn.WriteJSON(ServerMessage{
					Type:  ServerMessageTypeError,
					Error: "failed to start sandbox: " + err.Error(),
				})
				continue
			}

			sandboxStarted = true

			sandboxResult, err := h.sandboxManager.Exec(
				c.Request.Context(),
				sandbox.ExecCommandRequest{
					SessionID: sessionID,
					Command:   command,
				},
			)
			if err != nil {
				_ = conn.WriteJSON(ServerMessage{
					Type:  ServerMessageTypeError,
					Error: "failed to execute sandbox command: " + err.Error(),
				})
				continue
			}

			commandResult = &CommandResult{
				Stdout:   sandboxResult.Stdout,
				Stderr:   sandboxResult.Stderr,
				ExitCode: sandboxResult.ExitCode,
			}
		} else {
			commandResult, err = h.terminalService.GetCommandOutput(
				c.Request.Context(),
				*currentSession.CurrentStepID,
				command,
			)
			if err != nil {
				_ = conn.WriteJSON(ServerMessage{
					Type:  ServerMessageTypeError,
					Error: "failed to get command output",
				})
				continue
			}
		}

		combinedOutput := commandResult.Stdout
		if commandResult.Stderr != "" {
			combinedOutput += commandResult.Stderr
		}

		exitCode := commandResult.ExitCode

		err = conn.WriteJSON(ServerMessage{
			Type:     ServerMessageTypeOutput,
			Output:   combinedOutput,
			Stdout:   commandResult.Stdout,
			Stderr:   commandResult.Stderr,
			ExitCode: &exitCode,
		})
		if err != nil {
			return
		}

		actionResult, err := h.actionsService.Submit(
			c.Request.Context(),
			sessionID,
			userID,
			actions.SubmitActionRequest{
				StepID:  *currentSession.CurrentStepID,
				Command: command,
			},
		)
		if err != nil {
			_ = conn.WriteJSON(ServerMessage{
				Type:  ServerMessageTypeError,
				Error: err.Error(),
			})
			continue
		}

		updatedSession, err := h.sessionsService.GetByID(c.Request.Context(), sessionID, userID)
		if err != nil {
			_ = conn.WriteJSON(ServerMessage{
				Type:         ServerMessageTypeActionResult,
				ActionResult: actionResult,
			})
			continue
		}

		if updatedSession.Status == "completed" {
			go stopSandboxInBackground(h.sandboxManager, sessionID)
		}

		err = conn.WriteJSON(ServerMessage{
			Type:          ServerMessageTypeActionResult,
			SessionID:     updatedSession.ID.String(),
			CurrentStepID: updatedSession.CurrentStepID,
			ActionResult:  actionResult,
		})
		if err != nil {
			return
		}
	}
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

func dockerImageForScenario(scenarioID string) (string, bool) {
	switch scenarioID {
	case "permissions-junior":
		return "dcs-scenario-permissions-junior:dev", true
	default:
		return "", false
	}
}

func stopSandboxInBackground(manager *sandbox.Manager, sessionID uuid.UUID) {
	if manager == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = manager.Stop(ctx, sessionID)
}
