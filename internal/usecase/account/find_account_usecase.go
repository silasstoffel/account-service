package usecase

import (
	"log"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
)

type FindAccount struct {
	AccountRepository accountDomain.AccountRepository
}

func (ref *FindAccount) FindAccountUseCase(id string) (accountDomain.Account, error) {
	const prefix = "[find-account-usecase]"
	log.Println(prefix, "finding account", id)

	account, err := ref.AccountRepository.FindById(id)

	if err != nil {
		log.Println(prefix, "An error happens when find account", id)
		return accountDomain.Account{}, err
	}

	log.Println(prefix, "found account", id)

	return account, nil
}
