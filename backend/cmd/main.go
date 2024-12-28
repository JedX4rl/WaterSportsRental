package main

import (
	"WaterSportsRental/internal/configs/serverConfig"
	"WaterSportsRental/internal/configs/storageConfig"
	"WaterSportsRental/internal/handler"
	"WaterSportsRental/internal/logger/sl"
	"WaterSportsRental/internal/postgres"
	"WaterSportsRental/internal/repository"
	"WaterSportsRental/internal/service"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	logger "log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalf("Error loading .env file")
	}
}

func main() {

	serverCfg := serverConfig.MustLoadServerConfig()
	storageCfg := storageConfig.MustLoadStorageConfig()

	log := setupLogger(serverCfg.Env)

	dataBase, err := postgres.NewPostgresDb(storageCfg)

	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}

	log.Info("connected to database")

	defer func() {
		err = dataBase.Close()
		if err != nil {
			log.Error("got error when closing the DB connection", sl.Err(err))
			os.Exit(1)
		}
	}()

	repos := repository.NewRepository(dataBase)

	services := service.NewServices(repos)

	handlers := handler.NewHandler(services)

	log.Info("starting server")

	server := &http.Server{
		Addr:    serverCfg.Address,
		Handler: handlers.InitRoutes(),
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("starting server...", slog.String("env", serverCfg.Env))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen and serve error", sl.Err(err))
		}
	}()

	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", sl.Err(err))
	}

	log.Info("server exiting")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
