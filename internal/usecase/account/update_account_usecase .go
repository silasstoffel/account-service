package usecase

import (
	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/service"
)

type UpdateAccountInput struct {
	Name        string
	LastName    string
	Email       string
	Phone       string
	Password    string
	Permissions []string
}

func (ref *AccountUseCase) checkUpdateInput(input UpdateAccountInput, accountId string) error {
	var account accountDomain.Account
	var err error

	if input.Email != "" {
		account, err = ref.AccountRepository.FindByEmail(input.Email)
		if err != nil {
			detail := err.(*exception.Exception)
			if detail.Code != exception.AccountNotFound {
				return exception.New(exception.UnknownError, &err)
			}
		}

		if !account.IsEmpty() && account.Id != accountId {
			return exception.New(exception.AccountEmailAlreadyExists, &err)
		}
	}

	if input.Phone != "" {
		account, err = ref.AccountRepository.FindByPhone(input.Phone)
		if err != nil {
			detail := err.(*exception.Exception)
			if detail.Code != exception.AccountNotFound {
				return exception.New(exception.UnknownError, &err)
			}
		}

		if !account.IsEmpty() && account.Id != accountId {
			return exception.New(exception.AccountEmailAlreadyExists, &err)
		}
	}

	return nil
}

func (ref *AccountUseCase) UpdateAccountUseCase(id string, input UpdateAccountInput) (accountDomain.Account, error) {
	const loggerPrefix = "[update-account-usecase]"

	if err := ref.checkUpdateInput(input, id); err != nil {
		ref.Logger.Error(loggerPrefix+"Error when creating password", err, nil)
		return accountDomain.Account{}, err
	}

	var pwd string
	if input.Password != "" {
		var err error
		pwd, err = service.CreateHash(input.Password)
		if err != nil {
			ref.Logger.Error(loggerPrefix+"Error creating password hash", err, nil)
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
		ref.Logger.Error(loggerPrefix+"Error when updating account", err, nil)
		return accountDomain.Account{}, err
	}

	permissions, err := createAccountPermissions(input.Permissions, updatedAccount.Id, ref.AccountPermissionRepository)
	if err != nil {
		ref.Logger.Error(loggerPrefix+"Error when creating account", err, nil)
		return accountDomain.Account{}, err
	}
	updatedAccount.Permissions = permissions
	data := updatedAccount.ToDomain()

	go ref.Messaging.Publish(event.AccountUpdated, data, "account-service")

	return data, nil
}
