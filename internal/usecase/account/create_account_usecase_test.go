package usecase

import (
	"testing"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/logger"
)

type AccountRepositoryMock struct{}
type MessagingMock struct{}
type AccountPermissionRepositoryMock struct{}

var createdAccountMock domain.Account

func (ref AccountPermissionRepositoryMock) FindByAccountId(accountId string) ([]domain.AccountPermission, error) {
	return []domain.AccountPermission{}, nil
}

func (ref AccountPermissionRepositoryMock) Create(data domain.CreateAccountPermissionInput) error {
	return nil
}

func (ref AccountPermissionRepositoryMock) DeleteByAccount(accountId string) error {
	return nil
}

func (ref AccountRepositoryMock) FindByEmail(email string) (domain.Account, error) {
	return domain.Account{Id: ""}, nil
}

func (ref AccountRepositoryMock) FindByPhone(phone string) (domain.Account, error) {
	return domain.Account{Id: ""}, nil
}

func (ref AccountRepositoryMock) Create(account domain.Account) (domain.Account, error) {
	return createdAccountMock, nil
}

func (ref AccountRepositoryMock) Update(id string, account domain.Account) (domain.Account, error) {
	return domain.Account{}, nil
}

func (ref AccountRepositoryMock) FindById(id string) (domain.Account, error) {
	return domain.Account{}, nil
}

func (ref AccountRepositoryMock) List(input domain.ListAccountInput) ([]domain.Account, error) {
	return []domain.Account{}, nil
}

func (ref *MessagingMock) Publish(eventType string, data interface{}, source string) error {
	return nil
}

func TestCreateAccountUseCase(t *testing.T) {
	repositoryMock := AccountRepositoryMock{}
	AccountPermissionRepositoryMock := AccountPermissionRepositoryMock{}
	messagingMock := MessagingMock{}
	logger := logger.Logger{
		Env:     "testing",
		Service: "account-service",
	}

	createdAccountMock = domain.Account{
		Id:        "123",
		Name:      "Silas",
		LastName:  "Stoffel",
		Email:     "mail@mail.com",
		Phone:     "123456789",
		HashedPwd: "HashedPassword",
		Permissions: []domain.AccountPermission{
			{
				AccountId:    "123",
				PermissionId: "1",
				Scope:        "admin",
				Active:       true,
				CreatedAt:    "2021-01-01",
			},
		},
	}

	usecase := NewAccountUseCase(repositoryMock, AccountPermissionRepositoryMock, &messagingMock, &logger)
	result, err := usecase.CreateAccountUseCase(CreateAccountInput{
		Name:        createdAccountMock.Name,
		LastName:    createdAccountMock.LastName,
		Email:       createdAccountMock.Email,
		Phone:       createdAccountMock.Phone,
		Password:    "123456",
		Permissions: []string{"admin"},
	})

	if err != nil {
		t.Error("Expected nil, got ", err)
	}
	if result.Id != "123" {
		t.Error("Expected id, got ", result.Id)
	}
	if result.HashedPwd != "" {
		t.Error("Expected hashed password, got ", result.HashedPwd)
	}
	if result.FullName != "Silas Stoffel" {
		t.Error("Expected full name, got ", result.FullName)
	}
}
