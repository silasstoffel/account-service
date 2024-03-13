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

	_, err := ref.AccountRepository.FindByEmail(input.Email)
	if err != nil {
		return domain.Account{}, domain.NewError(domain.AccountEmailAlreadyExists, "Email already registered", err)
	}

	_, err = ref.AccountRepository.FindByPhone(input.Phone)
	if err != nil {
		return domain.Account{}, domain.NewError(domain.AccountPhoneAlreadyExists, "Phone already registered", err)
	}

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
		return domain.Account{}, err
	}

	log.Println(loggerPrefix, "Account created", "id:", createdAccount.Id)

	return createdAccount.ToDomain(), nil
}
