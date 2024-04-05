package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/infra/http/handler"
	v1handler "github.com/silasstoffel/account-service/internal/infra/http/handler/v1"
)

func BuildRouter(config *configs.Config, db *sql.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	route := gin.Default()

	route.GET("/health-check", handler.HealthCheckHandler())

	v1handler.GetAuthHandler(route, config, db)

	v1Group := route.Group("/v1")
	v1handler.GetAccountHandler(v1Group, config, db)

	return route
}
