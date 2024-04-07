package usecase

import (
	"log"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	domain "github.com/silasstoffel/account-service/internal/domain/account"
)

type ListAccount struct {
	AccountRepository           domain.AccountRepository
	AccountPermissionRepository accountDomain.AccountPermissionRepository
}

func (ref *ListAccount) ListAccountUseCase(input domain.ListAccountInput) ([]domain.Account, error) {
	const listLoggerPrefix = "[list-account-usecase]"

	accounts, err := ref.AccountRepository.List(input)

	if err != nil {
		log.Println(listLoggerPrefix, "Error when listing accounts", err)
		return []domain.Account{}, err
	}

	for key, _ := range accounts {
		p, err := ref.AccountPermissionRepository.FindByAccountId(accounts[key].Id)
		if err != nil {
			log.Println(listLoggerPrefix, "Error when listing account permissions", err)
			return []domain.Account{}, err
		}
		accounts[key].Permissions = p
	}

	return accounts, nil
}
