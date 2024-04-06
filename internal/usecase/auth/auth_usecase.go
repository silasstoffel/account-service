package usecase

import (
	"log"
	"time"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/domain/auth"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/service"
)

type AuthParams struct {
	AccountRepository           domain.AccountRepository
	PermissionAccountRepository domain.AccountPermissionRepository
	Messaging                   event.EventProducer
	TokenService                auth.TokenManagerService
}

type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
		if detail.Code != domain.AccountNotFound {
			message := "Error when find account by e-mail"
			log.Println(message, err)
			return nil, exception.New(exception.UnknownError, message, err, exception.HttpInternalError)
		}
		return nil, exception.New(domain.InvalidUserOrPassword, "Invalid user or password", err, exception.HttpUnauthorized)
	}

	if err := service.CompareHash(data.Password, account.HashedPwd); err != nil {
		return nil, exception.New(domain.InvalidUserOrPassword, "Invalid user or password", err, exception.HttpUnauthorized)
	}

	token, err := ref.TokenService.CreateToken(account.Id)
	if err != nil {
		message := "Error when create token"
		log.Println(message, err)
		return nil, exception.New(exception.UnknownError, message, err, exception.HttpInternalError)
	}

	result, err := ref.PermissionAccountRepository.FindByAccountId(account.Id)
	var permissions []string
	if err == nil {
		for _, p := range result {
			permissions = append(permissions, p.Scope)
		}
	}

	return &AuthOutput{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		Permissions: permissions,
	}, nil
}
