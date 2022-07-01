package main

import (
	server "allincecup-server"
	"allincecup-server/pkg/handler"
	"allincecup-server/pkg/repository"
	"allincecup-server/pkg/service"
	"log"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
