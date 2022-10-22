package main

import (
	_ "github.com/lib/pq"
	server "github.com/zh0vtyj/allincecup-server"
	"github.com/zh0vtyj/allincecup-server/internal/adapters/handler"
	"github.com/zh0vtyj/allincecup-server/internal/config"
	"github.com/zh0vtyj/allincecup-server/internal/domain/repository"
	"github.com/zh0vtyj/allincecup-server/internal/domain/service"
	"github.com/zh0vtyj/allincecup-server/pkg/client/minio"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
)

// @title AllianceCup API
// @version 1.0
// @description API Server for AllianceCup Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logger := logging.GetLogger()

	logger.Info("config initializing...")
	cfg := config.GetConfig()

	logger.Info("storage initializing...")
	storage, err := postgres.NewPostgresDB(cfg.Storage)
	if err != nil {
		logger.Fatalf("error occured while initializing db: %v", err)
	}

	logger.Info("minio client initializing...")
	_, err = minioPkg.NewClient(cfg.MinIO)
	if err != nil {
		logger.Fatalf("error occured while initializing minio client: %v", err)
	}

	logger.Info("repository initializing...")
	repos := repository.NewRepository(storage, logger)

	logger.Info("service initializing...")
	services := service.NewService(repos, logger)

	logger.Info("handler initializing...")
	handlers := handler.NewHandler(services, logger)

	logger.Info("running the server...")
	srv := new(server.Server)
	if err = srv.Run(cfg.AppPort, handlers.InitRoutes()); err != nil {
		logger.Fatalf("error occured while running http server: %s", err.Error())
	}
}
