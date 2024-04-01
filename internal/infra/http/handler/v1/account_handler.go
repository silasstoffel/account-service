package v1handler

import (
	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/configs"
	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/infra/messaging"
	usecase "github.com/silasstoffel/account-service/internal/usecase/account"
)

var accountRepository *database.AccountRepository
var messagingProducer *messaging.MessagingProducer
var accountPermissionRepository *database.AccountPermissionRepository

func GetAccountHandler(router *gin.RouterGroup, config *configs.Config) {
	cnx := database.OpenConnection(config)
	accountRepository = database.NewAccountRepository(cnx)
	accountPermissionRepository = database.NewAccountPermissionRepository(cnx)

	messagingProducer = messaging.NewMessagingProducer(
		config.Aws.AccountServiceTopicArn,
		config.Aws.Endpoint,
	)

	group := router.Group("/accounts")
	group.GET("/", list())
	group.GET("/:id", get())
	group.POST("/", create())
	group.PUT("/:id", update())
	group.PATCH("/:id/disabled", create())
}

func list() gin.HandlerFunc {
	return func(c *gin.Context) {
		listAccount := usecase.ListAccount{
			AccountRepository:           accountRepository,
			PermissionAccountRepository: accountPermissionRepository,
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

func get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		findAccount := usecase.FindAccount{
			AccountRepository:           accountRepository,
			PermissionAccountRepository: accountPermissionRepository,
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

func create() gin.HandlerFunc {
	return func(c *gin.Context) {
		createAccount := usecase.CreateAccount{
			AccountRepository:           accountRepository,
			PermissionAccountRepository: accountPermissionRepository,
			Messaging:                   messagingProducer,
		}
		var input usecase.CreateAccountInput

		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, gin.H{"code": "INVALID_INPUT_FORMAT", "message": "Invalid input format"})
			return
		}

		account, err := createAccount.CreateAccountUseCase(input)
		if err != nil {
			detail := err.(*exception.Exception)
			c.JSON(400, detail.ToDomain())
			return
		}

		c.JSON(201, account)
	}
}

func update() gin.HandlerFunc {
	return func(c *gin.Context) {
		updateAccountInstance := usecase.UpdateAccount{
			AccountRepository:           accountRepository,
			Messaging:                   messagingProducer,
			PermissionAccountRepository: accountPermissionRepository,
		}
		var input usecase.UpdateAccountInput

		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, gin.H{"code": "INVALID_INPUT_FORMAT", "message": "Invalid input format"})
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
