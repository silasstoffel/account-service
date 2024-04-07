package usecase

import (
	"log"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/service"
)

type CreateAccountInput struct {
	Name        string
	LastName    string
	Email       string
	Phone       string
	Password    string
	Permissions []string
}

type CreateAccount struct {
	AccountRepository           accountDomain.AccountRepository
	AccountPermissionRepository accountDomain.AccountPermissionRepository
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
	}
	account.BuildFullName()

	createdAccount, err := ref.AccountRepository.Create(account)

	if err != nil {
		log.Println(loggerPrefix, "Error when creating account. Detail:", err)
		return accountDomain.Account{}, err
	}

	permissions, err := createAccountPermissions(input.Permissions, createdAccount.Id, ref.AccountPermissionRepository)
	if err != nil {
		log.Println(loggerPrefix, "Error when creating account. Detail:", err)
		return accountDomain.Account{}, err
	}
	createdAccount.Permissions = permissions

	data := createdAccount.ToDomain()

	go ref.Messaging.Publish(event.AccountCreated, data, "account-service")

	return data, nil
}

func createAccountPermissions(
	permissions []string,
	accountId string,
	accountPermissionRepository accountDomain.AccountPermissionRepository,
) ([]accountDomain.AccountPermission, error) {
	var createdPermissions []accountDomain.AccountPermission
	var err error
	if len(permissions) > 0 {
		accountPermissionRepository.DeleteByAccount(accountId)
		for _, permission := range permissions {
			p := accountDomain.CreateAccountPermissionInput{
				AccountId:    accountId,
				PermissionId: permission,
			}
			err := accountPermissionRepository.Create(p)
			if err != nil {
				message := "Error when creating account permission"
				log.Println(message, "Detail:", err)
				return nil, exception.New(exception.DbCommandError, message, err, 500)
			}
		}
		createdPermissions, err = accountPermissionRepository.FindByAccountId(accountId)
		if err != nil {
			message := "Error when querying account permission"
			log.Println(message, "Detail:", err)
			return nil, exception.New(exception.DbCommandError, message, err, 500)
		}
	}
	return createdPermissions, nil
}
