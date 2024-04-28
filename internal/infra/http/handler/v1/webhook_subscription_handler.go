package v1handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/domain/webhook"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/infra/helper"
	"github.com/silasstoffel/account-service/internal/infra/http/middleware"
	"github.com/silasstoffel/account-service/internal/infra/messaging"
	"github.com/silasstoffel/account-service/internal/logger"
	usecase "github.com/silasstoffel/account-service/internal/usecase/webhook"
	"github.com/silasstoffel/account-service/internal/utility"
)

var WebHookSubscriptionUse usecase.WebHookSubscriptionUseCaseParams

func GetWebHookSubscriptionHandler(router *gin.RouterGroup, config *configs.Config, db *sql.DB) {
	logger := logger.NewLogger(config)
	webhookSubscriptionRepository := database.NewSubscriptionRepository(db, logger)
	messagingProducer := messaging.NewDefaultMessagingProducerFromConfig(config, logger)
	WebHookSubscriptionUse = usecase.WebHookSubscriptionUseCaseParams{
		Messaging:                     messagingProducer,
		WebhookSubscriptionRepository: webhookSubscriptionRepository,
	}

	permissions := map[string]string{
		"POST|/v1/webhooks/subscriptions/":              "account-service:create-webhook-subscription,account-service:*",
		"PUT|/v1/webhooks/subscriptions/:id":            "account-service:update-webhook-subscription,account-service:*",
		"GET|/v1/webhooks/subscriptions/:id":            "account-service:find-webhook-subscription,account-service:*",
		"GET|/v1/webhooks/subscriptions/":               "account-service:find-webhook-subscription,account-service:*",
		"PATCH|/v1/webhooks/subscriptions/:id/active":   "account-service:update-webhook-subscription,account-service:*",
		"PATCH|/v1/webhooks/subscriptions/:id/inactive": "account-service:update-webhook-subscription,account-service:*",
	}
	authorizer := middleware.NewAuthorizerMiddleware(permissions)

	group := router.Group("/webhooks/subscriptions")
	group.POST("/", authorizer.AuthorizerMiddleware, createWebHookSubscription())
	group.PUT("/:id", authorizer.AuthorizerMiddleware, updateWebHookSubscription())
	group.GET("/", authorizer.AuthorizerMiddleware, listWebHookSubscriptions())
	group.GET("/:id", authorizer.AuthorizerMiddleware, findWebHookSubscription())
	group.PATCH("/:id/active", authorizer.AuthorizerMiddleware, changeWebHookSubscriptionStatus(true))
	group.PATCH("/:id/inactive", authorizer.AuthorizerMiddleware, changeWebHookSubscriptionStatus(false))
}

func createWebHookSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input webhook.CreateSubscriptionInput

		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, helper.InvalidInputFormat())
			return
		}

		sub, err := WebHookSubscriptionUse.CreateSubscriptionUseCase(input)
		if err != nil {
			e := err.(*exception.Exception)
			c.JSON(e.StatusCode, e)
			return
		}

		c.JSON(201, sub)
	}
}

func updateWebHookSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input webhook.UpdateSubscriptionInput

		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, helper.InvalidInputFormat())
			return
		}
		input.Active = true
		sub, err := WebHookSubscriptionUse.UpdateSubscriptionUseCase(c.Param("id"), input)
		if err != nil {
			e := err.(*exception.Exception)
			c.JSON(e.StatusCode, e)
			return
		}

		c.JSON(200, sub)
	}
}

func changeWebHookSubscriptionStatus(status bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := WebHookSubscriptionUse.ChangSubscriptionStatusUseCase(c.Param("id"), status)
		if err != nil {
			e := err.(*exception.Exception)
			c.JSON(e.StatusCode, e)
			return
		}
		c.JSON(204, nil)
	}
}

func findWebHookSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, err := WebHookSubscriptionUse.FindSubscriptionUseCase(c.Param("id"))
		if err != nil {
			e := err.(*exception.Exception)
			c.JSON(e.StatusCode, e)
			return
		}

		c.JSON(200, sub)
	}
}

func listWebHookSubscriptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		subs, err := WebHookSubscriptionUse.ListSubscriptionUseCase(webhook.ListSubscriptionInput{
			Page:  utility.StrToInt(c.Query("page"), 1),
			Limit: utility.StrToInt(c.Query("limit"), 10),
		})

		if err != nil {
			e := err.(*exception.Exception)
			c.JSON(e.StatusCode, e)
			return
		}

		c.JSON(200, subs)
	}
}
