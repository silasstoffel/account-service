package usecase

import (
	"log"

	"github.com/silasstoffel/account-service/internal/domain"
)

const loggerPrefix = "[create-account-usecase]"

type CreateAccount struct {
	AccountRepository domain.AccountRepository
}

func (ref *CreateAccount) CreateAccountUseCase(input CreateAccountInput) (domain.Account, error) {
	log.Println(loggerPrefix, "Creating account...")

	account := domain.Account{
		Name:      input.Name,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		HashedPwd: input.HashedPwd,
		Active:    true,
	}

	createdAccount, _ := ref.AccountRepository.Create(account)
	log.Println(loggerPrefix, "Account created", "id:", createdAccount.Id)

	return createdAccount.ToDomain(), nil
}
