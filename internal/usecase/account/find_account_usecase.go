package usecase

import (
	"log"

	"github.com/silasstoffel/account-service/internal/domain"
)

type FindAccount struct {
	AccountRepository domain.AccountRepository
}

func (ref *FindAccount) FindAccountUseCase(id string) (domain.Account, error) {
	const prefix = "[find-account-usecase]"
	log.Println(prefix, "finding account", id)

	account, err := ref.AccountRepository.FindById(id)

	if err != nil {
		log.Println(prefix, "An error happens when find account", id)
		return domain.Account{}, err
	}

	log.Println(prefix, "found account", id)

	return account, nil
}
