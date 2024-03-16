package usecase

import (
	"log"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
)

type ListAccount struct {
	AccountRepository domain.AccountRepository
}

func (ref *ListAccount) ListAccountUseCase(input domain.ListAccountInput) ([]domain.Account, error) {
	const listLoggerPrefix = "[list-account-usecase]"
	log.Println(listLoggerPrefix, "Listing accounts")

	accounts, err := ref.AccountRepository.List(input)

	if err != nil {
		return []domain.Account{}, err
	}

	log.Println(listLoggerPrefix, "Listed accounts")

	return accounts, nil
}
