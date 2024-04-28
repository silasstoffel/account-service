package middleware

import (
	"bytes"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/silasstoffel/account-service/internal/infra/helper"
)

type BodyValidatorMiddlewareParams struct {
	Input any
}

func NewBodyValidatorMiddleware(i any) *BodyValidatorMiddlewareParams {
	return &BodyValidatorMiddlewareParams{Input: i}
}

func (ref *BodyValidatorMiddlewareParams) BodyValidatorMiddleware(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.InvalidInputFormat())
		return
	}

	inputType := reflect.TypeOf(ref.Input)
	if inputType.Kind() != reflect.Ptr {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.InvalidInputFormat())
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	data := reflect.New(inputType.Elem()).Interface()
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.InvalidInputFormat())
		return
	}

	validate := validator.New()
	err = validate.Struct(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.ValidationFailure(err.Error()))
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	c.Next()
}
