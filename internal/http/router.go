package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"DCS/internal/actions"
	"DCS/internal/auth"
	"DCS/internal/hints"
	"DCS/internal/http/middleware"
	"DCS/internal/reports"
	"DCS/internal/sandbox"
	"DCS/internal/scenarios"
	"DCS/internal/sessions"
	"DCS/internal/terminal"
)

type RouterDeps struct {
	DB               *pgxpool.Pool
	AuthHandler      *auth.Handler
	ScenariosHandler *scenarios.Handler
	SessionsHandler  *sessions.Handler
	ActionsHandler   *actions.Handler
	ReportsHandler   *reports.Handler
	JWTManager       *auth.JWTManager
	TerminalHandler  *terminal.Handler
	SandboxHandler   *sandbox.Handler
	HintsHandler     *hints.Handler
}

func NewRouter(deps RouterDeps) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		if err := deps.DB.Ping(c.Request.Context()); err != nil {
			c.JSON(nethttp.StatusServiceUnavailable, gin.H{
				"status": "error",
				"error":  "database is unavailable",
			})
			return
		}

		c.JSON(nethttp.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api")

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", deps.AuthHandler.Register)
		authGroup.POST("/login", deps.AuthHandler.Login)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(deps.JWTManager))
	{
		protected.GET("/me", func(c *gin.Context) {
			userID, _ := c.Get("userId")
			email, _ := c.Get("email")
			role, _ := c.Get("role")

			c.JSON(nethttp.StatusOK, gin.H{
				"id":    userID,
				"email": email,
				"role":  role,
			})
		})

		scenariosGroup := protected.Group("/scenarios")
		{
			scenariosGroup.GET("", deps.ScenariosHandler.GetAll)
			scenariosGroup.GET("/:id", deps.ScenariosHandler.GetByID)
		}

		sessionsGroup := protected.Group("/sessions")
		{
			sessionsGroup.GET("", deps.SessionsHandler.GetHistory)
			sessionsGroup.POST("", deps.SessionsHandler.Create)
			sessionsGroup.GET("/:id", deps.SessionsHandler.GetByID)
			sessionsGroup.POST("/:id/actions", deps.ActionsHandler.Submit)
			sessionsGroup.GET("/:id/report", deps.ReportsHandler.GetSessionReport)
			sessionsGroup.GET("/:id/terminal", deps.TerminalHandler.Connect)
			sessionsGroup.POST("/:id/hints", deps.HintsHandler.Create)
		}

		sandboxGroup := protected.Group("/sandbox")
		{
			sandboxGroup.GET("/:sessionId", deps.SandboxHandler.Get)
			sandboxGroup.POST("/cleanup", deps.SandboxHandler.Cleanup)
			sandboxGroup.POST("/start", deps.SandboxHandler.Start)
			sandboxGroup.POST("/:sessionId/exec", deps.SandboxHandler.Exec)
			sandboxGroup.DELETE("/:sessionId", deps.SandboxHandler.Stop)
		}
	}

	return router
}
