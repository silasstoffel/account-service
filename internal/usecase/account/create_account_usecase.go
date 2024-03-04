package usecase

import (
	"fmt"
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
		HashedPwd: input.Password,
		Active:    true,
		FullName:  fmt.Sprintf("%s %s", input.Name, input.LastName),
	}

	createdAccount, err := ref.AccountRepository.Create(account)

	if err != nil {
		return createdAccount, err
	}

	log.Println(loggerPrefix, "Account created", "id:", createdAccount.Id)

	return createdAccount.ToDomain(), nil
}
