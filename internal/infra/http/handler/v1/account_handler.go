package v1handler

import (
	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/internal/domain"
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
		listAccount := usecase.ListAccount{AccountRepository: accountRepository}
		input := domain.ListAccountInput{Page: 1, Limit: 12}
		accounts, err := listAccount.ListAccountUseCase(input)

		if err != nil {
			c.JSON(500, gin.H{"code": domain.UnknownError, "message": "Unknown error has happened"})
			return
		}
		c.JSON(200, accounts)
	}
}

func get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		findAccount := usecase.FindAccount{AccountRepository: accountRepository}
		account, err := findAccount.FindAccountUseCase(id)

		if err != nil {
			detail := err.(*domain.Error)

			if detail.Code == domain.AccountNotFound {
				c.JSON(404, detail.ToDomain())
				return
			}

			c.JSON(500, gin.H{"code": domain.UnknownError, "message": "Unknown error has happened"})
			return
		}
		c.JSON(200, account)
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
