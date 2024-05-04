package usecase_test

import (
	"testing"
	"time"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/event"
	"github.com/silasstoffel/account-service/internal/logger"
	"github.com/silasstoffel/account-service/internal/test/mock"
	usecase "github.com/silasstoffel/account-service/internal/usecase/account"
	"go.uber.org/mock/gomock"
)

func TestUpdateAccountUseCase(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var l = logger.Logger{Env: "testing", Service: "account-service"}

	var account = domain.Account{
		Id:          "ulid:1",
		Name:        "Bruce",
		LastName:    "Wayne",
		Email:       "batman@dc.com",
		Phone:       "+1 222 333-4444",
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Permissions: []domain.AccountPermission{},
		HashedPwd:   "123456",
	}

	t.Run("Should update an account", func(t *testing.T) {
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

		input := usecase.UpdateAccountInput{
			Name:        account.Name,
			LastName:    account.LastName,
			Email:       account.Email,
			Phone:       "+1 222 333-4445",
			Password:    "123456",
			Permissions: []string{},
		}
		updatedAt := time.Now()
		accountRepository.EXPECT().FindByEmail(account.Email).Return(account, nil)
		accountRepository.EXPECT().FindByPhone(input.Phone).Return(account, nil)
		updatedAccount := account
		updatedAccount.Phone = input.Phone
		updatedAccount.UpdatedAt = updatedAt

		accountRepository.EXPECT().Update(account.Id, gomock.Any()).Return(updatedAccount, nil)
		messaging.EXPECT().Publish(event.AccountUpdated, gomock.Any(), gomock.Any()).Return(nil)

		act, err := accountUseCase.UpdateAccountUseCase(account.Id, input)
		if err != nil {
			t.Errorf("Error should be nil")
		}
		if act.Phone != input.Phone {
			t.Errorf("Phone should be updated")
		}
		if act.UpdatedAt != updatedAt {
			t.Errorf("UpdatedAt should be updated")
		}
	})

	t.Run("Should return error when email already exists", func(t *testing.T) {
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

		input := usecase.UpdateAccountInput{
			Name:        account.Name,
			LastName:    account.LastName,
			Email:       account.Email,
			Phone:       "+1 222 333-4445",
			Password:    "123456",
			Permissions: []string{},
		}
		existentAccount := account
		existentAccount.Id = "ulid:2"
		accountRepository.EXPECT().FindByEmail(account.Email).Return(existentAccount, nil)

		accountRepository.EXPECT().Update(account.Id, gomock.Any()).MaxTimes(0)
		messaging.EXPECT().Publish(event.AccountUpdated, gomock.Any(), gomock.Any()).Return(nil)

		_, err := accountUseCase.UpdateAccountUseCase(account.Id, input)
		if err != nil && err.Error() != "The email is already in use" {
			t.Errorf("Error should be AccountEmailAlreadyExists but received %s", err.Error())
		}
	})

	t.Run("Should return error when phone already exists", func(t *testing.T) {
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

		input := usecase.UpdateAccountInput{
			Name:        account.Name,
			LastName:    account.LastName,
			Email:       account.Email,
			Phone:       "+1 222 333-4445",
			Password:    "123456",
			Permissions: []string{},
		}
		existentAccount := account
		existentAccount.Id = "ulid:2"
		accountRepository.EXPECT().FindByEmail(account.Email).Return(account, nil)
		accountRepository.EXPECT().FindByPhone(input.Phone).Return(existentAccount, nil)

		accountRepository.EXPECT().Update(account.Id, gomock.Any()).MaxTimes(0)
		messaging.EXPECT().Publish(event.AccountUpdated, gomock.Any(), gomock.Any()).MaxTimes(0)

		_, err := accountUseCase.UpdateAccountUseCase(account.Id, input)
		if err != nil && err.Error() != "The phone is already in use" {
			t.Errorf("Error should be The email is already in use, but received %s", err.Error())
		}
	})
}
