package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"DCS/internal/auth"
	"DCS/internal/http/middleware"
	"DCS/internal/scenarios"
	"DCS/internal/sessions"
)

type RouterDeps struct {
	DB               *pgxpool.Pool
	AuthHandler      *auth.Handler
	ScenariosHandler *scenarios.Handler
	SessionsHandler  *sessions.Handler
	JWTManager       *auth.JWTManager
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
			sessionsGroup.POST("", deps.SessionsHandler.Create)
			sessionsGroup.GET("/:id", deps.SessionsHandler.GetByID)
		}
	}

	return router
}
