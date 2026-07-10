package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"

	"DCS/internal/actions"
	"DCS/internal/auth"
	"DCS/internal/config"
	"DCS/internal/database"
	apphttp "DCS/internal/http"
	"DCS/internal/scenarios"
	"DCS/internal/sessions"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	cfg := config.Load()

	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	jwtManager := auth.NewJWTManager(
		cfg.JWTSecret,
		time.Duration(cfg.JWTTTLHours)*time.Hour,
	)

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, jwtManager)
	authHandler := auth.NewHandler(authService)

	scenariosRepo := scenarios.NewRepository(db)
	scenariosService := scenarios.NewService(scenariosRepo)
	scenariosHandler := scenarios.NewHandler(scenariosService)

	sessionsRepo := sessions.NewRepository(db)
	sessionsService := sessions.NewService(sessionsRepo)
	sessionsHandler := sessions.NewHandler(sessionsService)

	actionsRepo := actions.NewRepository(db)
	actionsService := actions.NewService(actionsRepo)
	actionsHandler := actions.NewHandler(actionsService)

	router := apphttp.NewRouter(apphttp.RouterDeps{
		DB:               db,
		AuthHandler:      authHandler,
		ScenariosHandler: scenariosHandler,
		SessionsHandler:  sessionsHandler,
		ActionsHandler:   actionsHandler,
		JWTManager:       jwtManager,
	})

	addr := ":" + cfg.HTTPPort

	log.Printf("server started on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
