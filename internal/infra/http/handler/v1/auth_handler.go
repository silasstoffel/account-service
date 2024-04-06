package v1handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/infra/helper"
	"github.com/silasstoffel/account-service/internal/infra/messaging"
	"github.com/silasstoffel/account-service/internal/infra/service/token"
	usecase "github.com/silasstoffel/account-service/internal/usecase/auth"
)

var authUseCaseParams *usecase.AuthParams

func GetAuthHandler(router *gin.Engine, config *configs.Config, db *sql.DB) {
	messagingProducer := messaging.NewMessagingProducer(
		config.Aws.AccountServiceTopicArn,
		config.Aws.Endpoint,
	)
	tokenManagerService := &token.TokenService{
		Secret:           config.AuthSecret,
		EmittedBy:        "account-service",
		ExpiresInMinutes: 60,
	}
	authUseCaseParams = &usecase.AuthParams{
		AccountRepository:           database.NewAccountRepository(db),
		PermissionAccountRepository: database.NewAccountPermissionRepository(db),
		Messaging:                   messagingProducer,
		TokenService:                tokenManagerService,
	}
	router.POST("/auth", auth())
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input usecase.AuthInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, helper.InvalidInputFormat())
			return
		}

		auth, err := authUseCaseParams.AuthenticateUseCase(&input)
		if err != nil {
			detail := err.(*exception.Exception)
			c.JSON(detail.HttpStatusCode, detail.ToDomain())
			return
		}

		c.JSON(200, auth)
	}
}
