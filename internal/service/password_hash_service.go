package service

import (
	"log"

	exception "github.com/silasstoffel/account-service/internal/domain/exception"
	"golang.org/x/crypto/bcrypt"
)

const (
	FailureWhenCreatePassword = "FAILURE_WHEN_CREATE_PASSWORD"
)

func CreateHash(value string) (string, error) {
	log.Println("[password-hash-service]", "Creating hash for password...")
	hash, err := bcrypt.GenerateFromPassword([]byte(value), 15)
	if err != nil {
		log.Println("[password-hash-service]", "Failure create password hash", err.Error())
		return "", exception.NewError(FailureWhenCreatePassword, "Failure when create password", err)
	}
	log.Println("[password-hash-service]", "Created hash for password")
	return string(hash), nil
}
