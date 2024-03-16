package usecase

import (
	"log"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/service"
)

type UpdateAccountInput struct {
	Name     string
	LastName string
	Email    string
	Phone    string
	Password string
}

type UpdateAccount struct {
	AccountRepository accountDomain.AccountRepository
	Messaging         event.EventService
}

func (ref *UpdateAccount) checkInput(input UpdateAccountInput, accountId string) error {
	var account accountDomain.Account
	var err error

	if input.Email != "" {
		account, err = ref.AccountRepository.FindByEmail(input.Email)
		if err != nil {
			detail := err.(*exception.Exception)
			if detail.Code != accountDomain.AccountNotFound {
				return exception.New(exception.UnknownError, "Unknown error has happened", err, exception.HttpInternalError)
			}
		}

		if !account.IsEmpty() && account.Id != accountId {
			return exception.New(accountDomain.AccountEmailAlreadyExists, "Email already registered", err, exception.HttpClientError)
		}
	}

	if input.Phone != "" {
		account, err = ref.AccountRepository.FindByPhone(input.Phone)
		if err != nil {
			detail := err.(*exception.Exception)
			if detail.Code != accountDomain.AccountNotFound {
				return exception.New(exception.UnknownError, "Unknown error has happened", err, exception.HttpInternalError)
			}
		}

		if !account.IsEmpty() && account.Id != accountId {
			return exception.New(accountDomain.AccountEmailAlreadyExists, "Phone already registered", err, exception.HttpClientError)
		}
	}

	return nil
}

func (ref *UpdateAccount) UpdateAccountUseCase(id string, input UpdateAccountInput) (accountDomain.Account, error) {
	const loggerPrefix = "[update-account-usecase]"
	log.Println(loggerPrefix, "Updating account...")

	if err := ref.checkInput(input, id); err != nil {
		return accountDomain.Account{}, err
	}

	var pwd string
	if input.Password != "" {
		var err error
		pwd, err = service.CreateHash(input.Password)
		if err != nil {
			return accountDomain.Account{}, err
		}
	}

	account := accountDomain.Account{
		Name:      input.Name,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		HashedPwd: pwd,
	}
	account.BuildFullName()
	updatedAccount, err := ref.AccountRepository.Update(id, account)

	if err != nil {
		return accountDomain.Account{}, err
	}

	data := updatedAccount.ToDomain()
	go ref.Messaging.Publish(event.AccountUpdated, data, "account-service")
	log.Println(loggerPrefix, "Account updated", "id:", id)

	return data, nil
}
