package v1handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/configs"
	domain "github.com/silasstoffel/account-service/internal/domain/account"
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
	accountRepository := database.NewAccountRepository(db)
	permissionAccountRepository := database.NewAccountPermissionRepository(db)
	authUseCaseParams = &usecase.AuthParams{
		AccountRepository:           accountRepository,
		PermissionAccountRepository: permissionAccountRepository,
		Messaging:                   messagingProducer,
		TokenService:                tokenManagerService,
	}
	router.POST("/auth", auth())
	router.GET("/auth/verify", verify(tokenManagerService, accountRepository, permissionAccountRepository))
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

func verify(
	tokenManagerService *token.TokenService,
	accountRepository domain.AccountRepository,
	accountPermissionRepository domain.AccountPermissionRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 8 {
			c.JSON(400, helper.InvalidInputFormat())
			return
		}
		token := authHeader[7:]
		data, err := tokenManagerService.VerifyToken(token)
		if err != nil {
			detail := err.(*exception.Exception)
			c.JSON(detail.HttpStatusCode, detail.ToDomain())
			return
		}

		account, err := accountRepository.FindById(data.Sub)
		if err != nil {
			detail := err.(*exception.Exception)
			c.JSON(detail.HttpStatusCode, detail.ToDomain())
			return
		}

		var permissions []string
		items, err := accountPermissionRepository.FindByAccountId(account.Id)
		if err == nil {
			for _, p := range items {
				permissions = append(permissions, p.Scope)
			}
		}
		c.JSON(200, gin.H{
			"account":     account.Id,
			"permissions": permissions,
		})
	}
}
