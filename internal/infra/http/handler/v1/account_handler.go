package v1handler

import "github.com/gin-gonic/gin"

func GetAccountHandler(router *gin.RouterGroup) {
	group := router.Group("/accounts")
	group.GET("/", list())
	group.GET("/:id", get())
	group.POST("/", create())
	group.PUT("/:id", update())
	group.PATCH("/:id/disabled", create())
}

func list() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "list accounts",
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
