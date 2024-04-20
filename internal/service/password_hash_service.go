package service

import (
	"github.com/silasstoffel/account-service/internal/exception"
	"golang.org/x/crypto/bcrypt"
)

func CreateHash(value string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), 15)
	if err != nil {
		return "", exception.New(exception.FailureToComparHash, &err)
	}
	return string(hash), nil
}

func CompareHash(value string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	if err != nil {
		return exception.New(exception.FailureToComparHash, &err)
	}
	return nil
}
