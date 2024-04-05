package v1handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/infra/messaging"
)

func GetAuthHandler(router *gin.Engine, config *configs.Config, db *sql.DB) {
	accountRepository = database.NewAccountRepository(db)
	accountPermissionRepository = database.NewAccountPermissionRepository(db)

	messagingProducer = messaging.NewMessagingProducer(
		config.Aws.AccountServiceTopicArn,
		config.Aws.Endpoint,
	)

	router.POST("/auth", auth())
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{})
	}
}
