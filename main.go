package main

import (
	"crypto-project/config"
	_ "crypto-project/docs"
	"crypto-project/internal/controller"
	"crypto-project/internal/usecase"
	"crypto-project/internal/usecase/repo/sqlitedb"
	"crypto-project/pkg/logger"
	"log/slog"
)

// @title Crypto-Project API
// @version 1.0
// @description This is a server for cryptocurrency with login and registration
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("Config error", slog.Any("error", err))
		return
	}
	l := logger.New(cfg.Log.Level)
	db, err := sqlitedb.New(cfg.DB.Path, l)
	if err != nil {
		l.Error("Database error", slog.Any("error", err))
		return

	}
	l.Info("connection to database")
	defer db.Close(l)
	u := usecase.New(db)
	server := controller.New(u, l)
	server.Run(cfg.HTTP.Port)
}
