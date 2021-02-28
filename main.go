package main

import (
	"log"

	"github.com/leonwright/reactor/logger"
	"github.com/leonwright/reactor/server"
	"github.com/leonwright/reactor/utils"
)

var deb = logger.SugaredLogger()

func main() {
	deb.Info("Starting Application...")
	deb.Info("Getting application configuration files.")
	var cfg utils.Config
	utils.ReadFile(&cfg)
	utils.ReadEnv(&cfg)
	if !utils.CheckForManagementAPIToken() {
		log.Println("Updating Management API Token...")
		utils.RequestNewManagementAPIToken(cfg)
		log.Println("Token successfully updated.")
	}

	go func() {
		server.StartServer(cfg)
	}()
	server.StartGrpcServer(cfg)
	log.Println("Exit method main()")
}
