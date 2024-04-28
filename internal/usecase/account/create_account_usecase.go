package usecase

import (
	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	loggerContract "github.com/silasstoffel/account-service/internal/logger/contract"
	"github.com/silasstoffel/account-service/internal/service"
)

type CreateAccountInput struct {
	Name        string `validate:"required,min=3,max=20"`
	LastName    string `validate:"required"`
	Email       string `validate:"required,email"`
	Phone       string `validate:"required"`
	Password    string `validate:"required,min=6"`
	Permissions []string
}

type AccountUseCase struct {
	AccountRepository           accountDomain.AccountRepository
	AccountPermissionRepository accountDomain.AccountPermissionRepository
	Messaging                   event.EventProducer
	Logger                      loggerContract.Logger
}

func NewAccountUseCase(
	accountRepository accountDomain.AccountRepository,
	accountPermissionRepository accountDomain.AccountPermissionRepository,
	messaging event.EventProducer,
	logger loggerContract.Logger,
) *AccountUseCase {
	return &AccountUseCase{
		AccountRepository:           accountRepository,
		AccountPermissionRepository: accountPermissionRepository,
		Messaging:                   messaging,
		Logger:                      logger,
	}
}

func (ref *AccountUseCase) checkInput(input CreateAccountInput) error {
	var account accountDomain.Account
	var err error

	account, err = ref.AccountRepository.FindByEmail(input.Email)
	if err != nil {
		detail := err.(*exception.Exception)
		if detail.Code != exception.AccountNotFound {
			return exception.NewUnknownError(&err)
		}
	}

	if !account.IsEmpty() {
		return exception.New(exception.AccountEmailAlreadyExists, &err)
	}

	account, err = ref.AccountRepository.FindByPhone(input.Phone)
	if err != nil {
		detail := err.(*exception.Exception)
		if detail.Code != exception.AccountNotFound {
			return exception.NewUnknownError(&err)
		}
	}

	if !account.IsEmpty() {
		return exception.New(exception.AccountEmailAlreadyExists, &err)
	}

	return nil
}

func (ref *AccountUseCase) CreateAccountUseCase(input CreateAccountInput) (accountDomain.Account, error) {
	const loggerPrefix = "[create-account-usecase]"

	if err := ref.checkInput(input); err != nil {
		return accountDomain.Account{}, err
	}

	pwd, err := service.CreateHash(input.Password)
	if err != nil {
		ref.Logger.Error(loggerPrefix+"Error creating password hash", err, nil)
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
		ref.Logger.Error(loggerPrefix+"Error when creating account", err, nil)
		return accountDomain.Account{}, err
	}

	permissions, err := createAccountPermissions(input.Permissions, createdAccount.Id, ref.AccountPermissionRepository)
	if err != nil {
		ref.Logger.Error(loggerPrefix+"Error when creating account", err, nil)
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
				return nil, exception.New(exception.DbCommandError, &err)
			}
		}
		createdPermissions, err = accountPermissionRepository.FindByAccountId(accountId)
		if err != nil {
			return nil, exception.New(exception.DbCommandError, &err)
		}
	}
	return createdPermissions, nil
}
