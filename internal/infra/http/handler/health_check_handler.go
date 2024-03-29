package handler

import "github.com/gin-gonic/gin"

func HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "app is running",
		})
	}
}
