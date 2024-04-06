package service

import (
	"log"

	"github.com/silasstoffel/account-service/internal/exception"
	"golang.org/x/crypto/bcrypt"
)

const (
	FailureWhenCreatePassword  = "FAILURE_WHEN_CREATE_PASSWORD"
	FailureWhenComparePassword = "FAILURE_WHEN_COMPARE_PASSWORD"
)

func CreateHash(value string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), 15)
	if err != nil {
		log.Println("Failure create password hash", err.Error())
		return "", exception.New(FailureWhenCreatePassword, "Failure when create password", err, exception.HttpInternalError)
	}
	return string(hash), nil
}

func CompareHash(value string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	if err != nil {
		message := "Failure when compare hash"
		log.Println(message, err.Error())
		return exception.New(FailureWhenComparePassword, message, err, exception.HttpInternalError)
	}
	return nil
}
