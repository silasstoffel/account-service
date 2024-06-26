package usecase

import (
	"time"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/domain/auth"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	loggerContract "github.com/silasstoffel/account-service/internal/logger/contract"
	"github.com/silasstoffel/account-service/internal/service"
)

type AuthParams struct {
	AccountRepository           domain.AccountRepository
	AccountPermissionRepository domain.AccountPermissionRepository
	Messaging                   event.EventProducer
	TokenService                auth.TokenManagerService
	Logger                      loggerContract.Logger
}

type AuthInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthOutput struct {
	AccessToken string    `json:"accessToken"`
	ExpiresIn   time.Time `json:"expiresIn"`
	Permissions []string  `json:"permissions"`
}

func (ref *AuthParams) AuthenticateUseCase(data *AuthInput) (*AuthOutput, error) {
	account, err := ref.AccountRepository.FindByEmail(data.Email)

	if err != nil {
		detail := err.(*exception.Exception)
		if detail.Code != exception.AccountNotFound {
			message := "Error when find account by e-mail"
			ref.Logger.Error(message, err, nil)
			return nil, exception.New(exception.UnknownError, &err)
		}
		return nil, exception.New(exception.InvalidUserOrPassword, &err)
	}

	if err := service.CompareHash(data.Password, account.HashedPwd); err != nil {
		return nil, exception.New(exception.InvalidUserOrPassword, &err)
	}

	token, err := ref.TokenService.CreateToken(account.Id)
	if err != nil {
		message := "Error when create token"
		ref.Logger.Error(message, err, nil)
		return nil, exception.New(exception.UnknownError, &err)
	}

	result, err := ref.AccountPermissionRepository.FindByAccountId(account.Id)
	var permissions []string
	if err == nil {
		for _, p := range result {
			permissions = append(permissions, p.Scope)
		}
	}

	go ref.Messaging.Publish(event.AccountLogged, account.ToDomain(), "account-service")

	return &AuthOutput{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		Permissions: permissions,
	}, nil
}
