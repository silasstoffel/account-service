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
	usecase "github.com/silasstoffel/account-service/internal/usecase/account"
)

var accountRepository *database.AccountRepository
var messagingProducer *messaging.MessagingProducer
var accountPermissionRepository *database.AccountPermissionRepository

func GetAccountHandler(router *gin.RouterGroup, config *configs.Config, db *sql.DB) {
	accountRepository = database.NewAccountRepository(db)
	accountPermissionRepository = database.NewAccountPermissionRepository(db)

	messagingProducer = messaging.NewDefaultMessagingProducerFromConfig(config)

	permissions := make(map[string]string)
	permissions["GET|/v1/accounts/"] = "account-service:list-accounts,account-service:*"
	permissions["POST|/v1/accounts/"] = "account-service:create-account,account-service:*"
	permissions["PUT|/v1/accounts/:id"] = "account-service:update-account,account-service:*"
	permissions["GET|/v1/accounts/:id"] = "account-service:get-account,account-service:*"
	authorizer := middleware.NewAuthorizerMiddleware(permissions)

	group := router.Group("/accounts")
	group.GET("/", authorizer.AuthorizerMiddleware, listAccount())
	group.GET("/:id", authorizer.AuthorizerMiddleware, getAccount())
	group.POST("/", authorizer.AuthorizerMiddleware, createAccount())
	group.PUT("/:id", authorizer.AuthorizerMiddleware, updateAccount())
}

func listAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		listAccount := usecase.ListAccount{
			AccountRepository:           accountRepository,
			AccountPermissionRepository: accountPermissionRepository,
		}
		input := domain.ListAccountInput{Page: 1, Limit: 12}
		accounts, err := listAccount.ListAccountUseCase(input)

		if err != nil {
			c.JSON(500, gin.H{"code": exception.UnknownError, "message": "Unknown error has happened"})
			return
		}
		c.JSON(200, accounts)
	}
}

func getAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		findAccount := usecase.FindAccount{
			AccountRepository:           accountRepository,
			AccountPermissionRepository: accountPermissionRepository,
		}
		account, err := findAccount.FindAccountUseCase(id)

		if err != nil {
			detail := err.(*exception.Exception)
			status := detail.HttpStatusCode

			if status < 500 {
				c.JSON(status, detail.ToDomain())
				return
			}

			c.JSON(detail.HttpStatusCode, gin.H{"code": exception.UnknownError, "message": "Unknown error has happened"})
			return
		}
		c.JSON(200, account)
	}
}

func createAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		createAccount := usecase.CreateAccount{
			AccountRepository:           accountRepository,
			AccountPermissionRepository: accountPermissionRepository,
			Messaging:                   messagingProducer,
		}
		var input usecase.CreateAccountInput

		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, helper.InvalidInputFormat())
			return
		}

		account, err := createAccount.CreateAccountUseCase(input)
		if err != nil {
			detail := err.(*exception.Exception)
			c.JSON(detail.HttpStatusCode, detail.ToDomain())
			return
		}

		c.JSON(201, account)
	}
}

func updateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		updateAccountInstance := usecase.UpdateAccount{
			AccountRepository:           accountRepository,
			Messaging:                   messagingProducer,
			AccountPermissionRepository: accountPermissionRepository,
		}
		var input usecase.UpdateAccountInput

		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, helper.InvalidInputFormat())
			return
		}

		account, err := updateAccountInstance.UpdateAccountUseCase(c.Param("id"), input)
		if err != nil {
			detail := err.(*exception.Exception)
			c.JSON(detail.HttpStatusCode, detail.ToDomain())
			return
		}

		c.JSON(200, account)
	}
}
