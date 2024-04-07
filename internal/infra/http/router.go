package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/infra/http/handler"
	v1handler "github.com/silasstoffel/account-service/internal/infra/http/handler/v1"
	"github.com/silasstoffel/account-service/internal/infra/http/middleware"
	"github.com/silasstoffel/account-service/internal/infra/service/token"
)

func BuildRouter(config *configs.Config, db *sql.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	route := gin.Default()

	route.GET("/health-check", handler.HealthCheckHandler())

	v1handler.GetAuthHandler(route, config, db)

	tokenManagerService := &token.TokenService{
		Secret:           config.AuthSecret,
		EmittedBy:        "account-service",
		ExpiresInMinutes: 60,
	}
	accountPermissionRepository := database.NewAccountPermissionRepository(db)

	verifyToken := middleware.NewVerifyTokenMiddleware(tokenManagerService, accountPermissionRepository)

	// protected routes
	v1Group := route.Group("/v1")
	v1Group.Use(verifyToken.VerifyTokenMiddleware)
	v1handler.GetAccountHandler(v1Group, config, db)

	return route
}
