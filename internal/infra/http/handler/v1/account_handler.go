package v1handler

import (
	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/internal/infra/database"
	usecase "github.com/silasstoffel/account-service/internal/usecase/account"
)

var accountRepository *database.AccountRepository

func GetAccountHandler(router *gin.RouterGroup) {
	cnx := database.OpenConnection()
	accountRepository = database.NewAccountRepository(cnx)

	group := router.Group("/accounts")
	group.GET("/", list())
	group.GET("/:id", get())
	group.POST("/", create())
	group.PUT("/:id", update())
	group.PATCH("/:id/disabled", create())
}

func list() gin.HandlerFunc {
	return func(c *gin.Context) {
		createAccount := usecase.CreateAccount{
			AccountRepository: accountRepository,
		}

		input := usecase.CreateAccountInput{
			Name:     "Silas",
			LastName: "Stoffel",
			Email:    "email@email.com",
			Phone:    "+55996354103",
			Password: "123456",
		}

		createdAccount, err := createAccount.CreateAccountUseCase(input)

		if err != nil {
			c.JSON(500, gin.H{"message": "internal server error", "error": err})
			return
		}

		c.JSON(200, gin.H{
			"data": createdAccount,
		})
	}
}

func get() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get account",
		})
	}
}

func create() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get account",
		})
	}
}

func update() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get account",
		})
	}
}
