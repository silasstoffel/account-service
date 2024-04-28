package v1handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/configs"
	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/infra/helper"
	"github.com/silasstoffel/account-service/internal/infra/http/middleware"
	"github.com/silasstoffel/account-service/internal/infra/messaging"
	"github.com/silasstoffel/account-service/internal/logger"
	usecase "github.com/silasstoffel/account-service/internal/usecase/account"
	"github.com/silasstoffel/account-service/internal/utility"
)

var accountUseCase *usecase.AccountUseCase

func GetAccountHandler(router *gin.RouterGroup, config *configs.Config, db *sql.DB) {
	logger := logger.NewLogger(config)
	accountRepository := database.NewAccountRepository(db, logger)
	accountPermissionRepository := database.NewAccountPermissionRepository(db, logger)
	messagingProducer := messaging.NewDefaultMessagingProducerFromConfig(config, logger)
	accountUseCase = usecase.NewAccountUseCase(accountRepository, accountPermissionRepository, messagingProducer, logger)

	permissions := map[string]string{
		"GET|/v1/accounts/":    "account-service:list-accounts,account-service:*",
		"POST|/v1/accounts/":   "account-service:create-account,account-service:*",
		"PUT|/v1/accounts/:id": "account-service:update-account,account-service:*",
		"GET|/v1/accounts/:id": "account-service:get-account,account-service:*",
	}
	authorizer := middleware.NewAuthorizerMiddleware(permissions)
	var createAccountSchema *usecase.CreateAccountInput
	createAccountValidator := middleware.NewBodyValidatorMiddleware(createAccountSchema)

	group := router.Group("/accounts")
	group.GET("/", authorizer.AuthorizerMiddleware, listAccount())
	group.GET("/:id", authorizer.AuthorizerMiddleware, getAccount())
	group.POST("/", authorizer.AuthorizerMiddleware, createAccountValidator.BodyValidatorMiddleware, createAccount())
	group.PUT("/:id", authorizer.AuthorizerMiddleware, updateAccount())
}

func listAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		accounts, err := accountUseCase.ListAccountUseCase(
			domain.ListAccountInput{
				Page:  utility.StrToInt(c.Query("page"), 1),
				Limit: utility.StrToInt(c.Query("limit"), 10),
			},
		)

		if err != nil {
			e := err.(*exception.Exception)
			c.JSON(e.GetStatusCode(), e)
			return
		}

		c.JSON(200, accounts)
	}
}

func getAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		account, err := accountUseCase.FindAccountUseCase(id)

		if err != nil {
			detail := err.(*exception.Exception)
			c.JSON(detail.StatusCode, detail)
			return
		}
		c.JSON(200, account)
	}
}

func createAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input usecase.CreateAccountInput

		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, helper.InvalidInputFormat())
			return
		}

		account, err := accountUseCase.CreateAccountUseCase(input)
		if err != nil {
			detail := err.(*exception.Exception)
			c.JSON(detail.StatusCode, detail)
			return
		}

		c.JSON(201, account)
	}
}

func updateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input usecase.UpdateAccountInput

		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, helper.InvalidInputFormat())
			return
		}

		account, err := accountUseCase.UpdateAccountUseCase(c.Param("id"), input)
		if err != nil {
			detail := err.(*exception.Exception)
			c.JSON(detail.StatusCode, detail)
			return
		}

		c.JSON(200, account)
	}
}
