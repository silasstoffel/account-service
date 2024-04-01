package usecase

import (
	"fmt"
	"log"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/service"
)

type AccountPermissionInput struct {
	AppId string `json:"appId"`
	Scope string `json:"scope"`
}
type CreateAccountInput struct {
	Name        string
	LastName    string
	Email       string
	Phone       string
	Password    string
	Permissions []AccountPermissionInput
}

type CreateAccount struct {
	AccountRepository           accountDomain.AccountRepository
	PermissionAccountRepository accountDomain.AccountPermissionRepository
	Messaging                   event.EventProducer
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
	}

	if !account.IsEmpty() {
		return exception.New(accountDomain.AccountEmailAlreadyExists, "Phone already registered", err, exception.HttpClientError)
	}

	return nil
}

func (ref *CreateAccount) CreateAccountUseCase(input CreateAccountInput) (accountDomain.Account, error) {
	const loggerPrefix = "[create-account-usecase]"

	if err := ref.checkInput(input); err != nil {
		return accountDomain.Account{}, err
	}

	pwd, err := service.CreateHash(input.Password)
	if err != nil {
		log.Println(loggerPrefix, "Error creating password hash", "detail:", err)
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
		log.Println(loggerPrefix, "Error when creating account. Detail:", err)
		return accountDomain.Account{}, err
	}

	if len(input.Permissions) > 0 {
		ref.PermissionAccountRepository.DeleteByAccount(createdAccount.Id)
		createdAccount.Permissions = []accountDomain.AccountPermission{}
		for _, permission := range input.Permissions {
			p := accountDomain.AccountPermission{
				AppId:     permission.AppId,
				Scope:     permission.Scope,
				AccountId: createdAccount.Id,
			}
			err = ref.PermissionAccountRepository.Create(p)
			if err != nil {
				log.Println(loggerPrefix, "Error when creating account permission. Detail:", err)
				return accountDomain.Account{}, err
			}
			createdAccount.Permissions = append(createdAccount.Permissions, p)
		}
	}
	data := createdAccount.ToDomain()

	go ref.Messaging.Publish(event.AccountCreated, data, "account-service")

	return data, nil
}
