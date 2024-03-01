package database

import (
	"log"
	"time"

	"github.com/silasstoffel/account-service/internal/domain"
	"github.com/silasstoffel/account-service/internal/infra/helper"
)

const loggerPrefix = "[account-repository]"

type AccountRepository struct{}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

func (repository *AccountRepository) Create(account domain.Account) (domain.Account, error) {
	log.Println(loggerPrefix, "Creating account...")
	now := time.Now().UTC()

	account.Id = helper.NewULID()
	account.CreatedAt = now
	account.UpdatedAt = now

	log.Println(loggerPrefix, "Account created with id", account.Id)
	return account, nil
}
