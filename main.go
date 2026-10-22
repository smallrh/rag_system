package main

import (
	"go_study/internal/config"
	"go_study/internal/db"
	handlerAuth "go_study/internal/handler/auth"
	repositoryAuth "go_study/internal/repository/auth"
	"go_study/internal/router"
	serviceAuth "go_study/internal/services/auth"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return
	}

	database, err := db.NewPostgres(cfg.Database.URL)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		return
	}
	defer func(database *sqlx.DB) {
		err := database.Close()
		if err != nil {

		}
	}(database)

	users := repositoryAuth.NewRepository(database)
	authService := serviceAuth.NewAuthService(users, []byte(cfg.JWTSecret))
	authHandler := handlerAuth.NewAuthHandler(authService)

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	router.RegisterRouters(e, authHandler)

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
