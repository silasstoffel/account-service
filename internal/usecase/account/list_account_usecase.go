package usecase

import (
	"log"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
)

const listLoggerPrefix = "[list-account-usecase]"

type ListAccount struct {
	AccountRepository domain.AccountRepository
}

func (ref *ListAccount) ListAccountUseCase(input domain.ListAccountInput) ([]domain.Account, error) {
	log.Println(listLoggerPrefix, "Listing accounts")

	accounts, err := ref.AccountRepository.List(input)

	if err != nil {
		return []domain.Account{}, err
	}

	log.Println(listLoggerPrefix, "Listed accounts")

	return accounts, nil
}
