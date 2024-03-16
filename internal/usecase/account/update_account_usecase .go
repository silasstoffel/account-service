package usecase

import (
	"log"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	errorDomain "github.com/silasstoffel/account-service/internal/domain/exception"
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
}

func (ref *UpdateAccount) checkInput(input UpdateAccountInput, accountId string) error {
	var account accountDomain.Account
	var err error

	if input.Email != "" {
		account, err = ref.AccountRepository.FindByEmail(input.Email)
		if err != nil {
			detail := err.(*errorDomain.Error)
			if detail.Code != accountDomain.AccountNotFound {
				return errorDomain.NewError(errorDomain.UnknownError, "Unknown error has happened", err)
			}
		}

		if !account.IsEmpty() && account.Id != accountId {
			return errorDomain.NewError(accountDomain.AccountEmailAlreadyExists, "Email already registered", err)
		}
	}

	if input.Phone != "" {
		account, err = ref.AccountRepository.FindByPhone(input.Phone)
		if err != nil {
			detail := err.(*errorDomain.Error)
			if detail.Code != accountDomain.AccountNotFound {
				return errorDomain.NewError(errorDomain.UnknownError, "Unknown error has happened", err)
			}
		}

		if !account.IsEmpty() && account.Id != accountId {
			return errorDomain.NewError(accountDomain.AccountEmailAlreadyExists, "Phone already registered", err)
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

	log.Println(loggerPrefix, "Account updated", "id:", id)

	return updatedAccount.ToDomain(), nil
}
