package usecase

import (
	"log"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
)

type FindAccount struct {
	AccountRepository           accountDomain.AccountRepository
	AccountPermissionRepository accountDomain.AccountPermissionRepository
}

func (ref *FindAccount) FindAccountUseCase(id string) (accountDomain.Account, error) {
	account, err := ref.AccountRepository.FindById(id)

	if err != nil {
		log.Println("[find-account-usecase] An error happens when find account", id)
		return accountDomain.Account{}, err
	}

	p, _ := ref.AccountPermissionRepository.FindByAccountId(account.Id)
	account.Permissions = p

	return account, nil
}
