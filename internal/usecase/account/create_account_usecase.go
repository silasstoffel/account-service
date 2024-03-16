package usecase

import (
	"fmt"
	"log"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	errorDomain "github.com/silasstoffel/account-service/internal/domain/exception"
	"github.com/silasstoffel/account-service/internal/service"
)

type CreateAccountInput struct {
	Name     string
	LastName string
	Email    string
	Phone    string
	Password string
}

type CreateAccount struct {
	AccountRepository accountDomain.AccountRepository
}

func (ref *CreateAccount) checkInput(input CreateAccountInput) error {
	var account accountDomain.Account
	var err error

	account, err = ref.AccountRepository.FindByEmail(input.Email)
	if err != nil {
		detail := err.(*errorDomain.Error)
		if detail.Code != accountDomain.AccountNotFound {
			return errorDomain.NewError(errorDomain.UnknownError, "Unknown error has happened", err)
		}
	}

	if !account.IsEmpty() {
		return errorDomain.NewError(accountDomain.AccountEmailAlreadyExists, "Email already registered", err)
	}

	account, err = ref.AccountRepository.FindByPhone(input.Phone)
	if err != nil {
		detail := err.(*errorDomain.Error)
		if detail.Code != accountDomain.AccountNotFound {
			return errorDomain.NewError(errorDomain.UnknownError, "Unknown error has happened", err)
		}
		return errorDomain.NewError(accountDomain.AccountPhoneAlreadyExists, "Phone already registered", err)
	}

	if !account.IsEmpty() {
		return errorDomain.NewError(accountDomain.AccountEmailAlreadyExists, "Phone already registered", err)
	}

	return nil
}

func (ref *CreateAccount) CreateAccountUseCase(input CreateAccountInput) (accountDomain.Account, error) {
	const loggerPrefix = "[create-account-usecase]"
	log.Println(loggerPrefix, "Creating account...")

	if err := ref.checkInput(input); err != nil {
		return accountDomain.Account{}, err
	}

	pwd, err := service.CreateHash(input.Password)
	if err != nil {
		return accountDomain.Account{}, err
	}

	account := accountDomain.Account{
		Name:      input.Name,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		HashedPwd: pwd,
		Active:    true,
		FullName:  fmt.Sprintf("%s %s", input.Name, input.LastName),
	}

	createdAccount, err := ref.AccountRepository.Create(account)

	if err != nil {
		return accountDomain.Account{}, err
	}

	log.Println(loggerPrefix, "Account created", "id:", createdAccount.Id)

	return createdAccount.ToDomain(), nil
}
