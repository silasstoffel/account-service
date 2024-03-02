package main

import (
	router "github.com/silasstoffel/account-service/internal/infra/http"
)

func main() {
	routes := router.BuildRouter()
	routes.Run(":8008")
}
