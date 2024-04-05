package main

import (
	"fmt"
	"log"

	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/infra/database"
	router "github.com/silasstoffel/account-service/internal/infra/http"
)

func main() {
	log.Println("Starting account-service REST API...")
	config := configs.NewConfigFromEnvVars()
	cnx, err := database.OpenConnection(config)
	if err != nil {
		log.Fatalf("Failed to open connection to database: %v", err)
		return
	}
	routes := router.BuildRouter(config, cnx)
	routes.Run(fmt.Sprintf(":%s", config.App.ApiPort))
}
