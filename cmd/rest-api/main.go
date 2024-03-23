package main

import (
	"fmt"

	"github.com/silasstoffel/account-service/configs"
	router "github.com/silasstoffel/account-service/internal/infra/http"
)

func main() {
	config := configs.NewConfigFromEnvVars()
	routes := router.BuildRouter(config)
	routes.Run(fmt.Sprintf(":%s", config.App.ApiPort))
}
