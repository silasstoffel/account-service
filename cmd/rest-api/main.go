package main

import (
	"fmt"

	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/infra/database"
	router "github.com/silasstoffel/account-service/internal/infra/http"
	"github.com/silasstoffel/account-service/internal/logger"
)

func main() {
	config := configs.NewConfigFromEnvVars()
	log := logger.NewLogger(config)

	log.Info("Starting account-service REST API...", nil)

	cnx, err := database.OpenConnection(config)
	if err != nil {
		log.Error("Failed to open connection to database", err, nil)
		return
	}
	routes := router.BuildRouter(config, cnx)
	routes.Run(fmt.Sprintf(":%s", config.App.ApiPort))
}
