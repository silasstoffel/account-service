package v1handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/configs"
	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/logger"
	usecase "github.com/silasstoffel/account-service/internal/usecase/permission"
	"github.com/silasstoffel/account-service/internal/utility"
)

func GetPermissionHandler(router *gin.RouterGroup, config *configs.Config, db *sql.DB) {
	logger := logger.NewLogger(config)
	var permissionsRepository = database.NewPermissionRepository(db, logger)
	var permissionUseCase = usecase.NewPermissionUseCase(permissionsRepository, logger)

	group := router.Group("/permissions")
	group.GET("/", func(c *gin.Context) {
		subs, err := permissionUseCase.ListPermissionUseCase(domain.ListPermissionInput{
			Page:  utility.StrToInt(c.Query("page"), 1),
			Limit: utility.StrToInt(c.Query("limit"), 10),
		})

		if err != nil {
			e := err.(*exception.Exception)
			c.JSON(e.StatusCode, e)
			return
		}

		c.JSON(200, subs)
	})
}
