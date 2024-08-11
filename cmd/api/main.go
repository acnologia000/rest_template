package main

import (
	"log"
	"net/http"
	"rest_template/internal/api"
	"rest_template/internal/config"
	"rest_template/pkg/database"
	"rest_template/shared/logger"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	loggerConfig := logger.Config{
		Level:        parseLogLevel(cfg.Logging.Level),
		OutputType:   cfg.Logging.OutputType,
		FilePath:     cfg.Logging.FilePath,
		CustomWriter: &logger.DBWriter{DB: db},
	}

	appLogger, err := logger.New(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	server := api.NewServer(cfg, appLogger, db)

	appLogger.Info("Starting server", map[string]interface{}{"address": cfg.GetListenAddress()})
	if err := http.ListenAndServe(cfg.GetListenAddress(), server); err != nil {
		appLogger.Error("Server failed to start", map[string]interface{}{"error": err.Error()})
	}
}

func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "debug":
		return logger.DEBUG
	case "info":
		return logger.INFO
	case "warn":
		return logger.WARN
	case "error":
		return logger.ERROR
	default:
		return logger.INFO
	}
}
