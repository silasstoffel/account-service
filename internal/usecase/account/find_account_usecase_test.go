package usecase_test

import (
	"errors"
	"testing"
	"time"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/exception"
	"github.com/silasstoffel/account-service/internal/logger"
	"github.com/silasstoffel/account-service/internal/test/mock"
	usecase "github.com/silasstoffel/account-service/internal/usecase/account"
	"go.uber.org/mock/gomock"
)

func TestFindAccountUseCaseSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var l = logger.Logger{Env: "testing", Service: "account-service"}
	accountRepository := mock.NewMockAccountRepository(ctrl)
	accountPermRepository := mock.NewMockAccountPermissionRepository(ctrl)

	id := "ulid:1"
	accountUseCase := usecase.NewAccountUseCase(
		accountRepository,
		accountPermRepository,
		nil,
		&l,
	)
	accountPermRepository.EXPECT().FindByAccountId(id).Return([]domain.AccountPermission{}, nil)
	accountRepository.EXPECT().FindById(id).Return(domain.Account{
		Id:          id,
		Name:        "Bruce",
		LastName:    "Wayne",
		Email:       "batman@dc.com",
		Phone:       "123456789",
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Permissions: []domain.AccountPermission{},
		HashedPwd:   "123456",
	}, nil)

	account, err := accountUseCase.FindAccountUseCase(id)
	if err != nil {
		t.Errorf("Error should be nil")
	}

	expectName := "Bruce Wayne"
	if account.FullName != expectName {
		t.Errorf("Expected %s, got %s", expectName, account.FullName)
	}
}

func TestFindAccountUseCaseUnSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var l = logger.Logger{Env: "testing", Service: "account-service"}
	accountRepository := mock.NewMockAccountRepository(ctrl)

	id := "ulid:2"
	accountUseCase := usecase.NewAccountUseCase(
		accountRepository,
		nil,
		nil,
		&l,
	)
	accountRepository.EXPECT().FindById(id).Return(
		domain.Account{},
		errors.New(exception.AccountNotFound),
	)
	_, err := accountUseCase.FindAccountUseCase(id)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
