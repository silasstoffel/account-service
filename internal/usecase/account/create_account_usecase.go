package usecase

import (
	"fmt"
	"log"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/exception"
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
		detail := err.(*exception.Exception)
		if detail.Code != accountDomain.AccountNotFound {
			return exception.New(exception.UnknownError, "Unknown error has happened", err, exception.HttpInternalError)
		}
	}

	if !account.IsEmpty() {
		return exception.New(accountDomain.AccountEmailAlreadyExists, "Email already registered", err, exception.HttpClientError)
	}

	account, err = ref.AccountRepository.FindByPhone(input.Phone)
	if err != nil {
		detail := err.(*exception.Exception)
		if detail.Code != accountDomain.AccountNotFound {
			return exception.New(exception.UnknownError, "Unknown error has happened", err, exception.HttpInternalError)
		}
		return exception.New(accountDomain.AccountPhoneAlreadyExists, "Phone already registered", err, exception.HttpClientError)
	}

	if !account.IsEmpty() {
		return exception.New(accountDomain.AccountEmailAlreadyExists, "Phone already registered", err, exception.HttpClientError)
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
