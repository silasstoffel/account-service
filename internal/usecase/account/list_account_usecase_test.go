package usecase_test

import (
	"errors"
	"testing"
	"time"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/logger"
	"github.com/silasstoffel/account-service/internal/test/mock"
	usecase "github.com/silasstoffel/account-service/internal/usecase/account"
	"go.uber.org/mock/gomock"
)

func TestListAccountUseCase(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var l = logger.Logger{Env: "testing", Service: "account-service"}

	t.Run("Should return an error", func(t *testing.T) {
		t.Parallel()

		accountRepository := mock.NewMockAccountRepository(ctrl)

		accountUseCase := usecase.NewAccountUseCase(
			accountRepository,
			nil,
			nil,
			&l,
		)

		accountRepository.EXPECT().List(gomock.Any()).Return(nil, errors.New("error"))

		_, err := accountUseCase.ListAccountUseCase(domain.ListAccountInput{})
		if err == nil {
			t.Errorf("Error should not be nil")
		}
	})

	t.Run("Should return accounts", func(t *testing.T) {
		t.Parallel()

		accountRepository := mock.NewMockAccountRepository(ctrl)
		accountPermRepository := mock.NewMockAccountPermissionRepository(ctrl)
		messaging := mock.NewMockEventProducer(ctrl)
		accountUseCase := usecase.NewAccountUseCase(
			accountRepository,
			accountPermRepository,
			messaging,
			&l,
		)

		input := domain.ListAccountInput{Page: 1, Limit: 10}
		account := domain.Account{
			Id:          "ulid:1",
			Name:        "Bruce",
			LastName:    "Wayne",
			Email:       "batman@dc.com",
			Phone:       "123456789",
			Active:      true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Permissions: []domain.AccountPermission{},
			HashedPwd:   "123456",
		}
		permission := domain.AccountPermission{
			AccountId:    account.Id,
			PermissionId: "ulid:1",
			Scope:        "admin",
			Active:       true,
			CreatedAt:    time.Now(),
		}
		accountRepository.EXPECT().List(input).Return([]domain.Account{account}, nil)
		accountPermRepository.EXPECT().FindByAccountId(account.Id).Return([]domain.AccountPermission{
			permission,
		}, nil)

		accounts, _ := accountUseCase.ListAccountUseCase(input)
		if len(accounts) == 0 {
			t.Errorf("Account should not be empty")
		}
		if condition := accounts[0].FullName == "Bruce Wayne"; !condition {
			t.Errorf("Expected Bruce Wayne, got %s", accounts[0].FullName)
		}
	})
}
