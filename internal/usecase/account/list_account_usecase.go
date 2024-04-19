package usecase

import (
	domain "github.com/silasstoffel/account-service/internal/domain/account"
)

func (ref *AccountUseCase) ListAccountUseCase(input domain.ListAccountInput) ([]domain.Account, error) {
	const listLoggerPrefix = "[list-account-usecase]"

	accounts, err := ref.AccountRepository.List(input)

	if err != nil {
		ref.Logger.Error(listLoggerPrefix+" Error when listing accounts", err, nil)
		return nil, err
	}

	for key := range accounts {
		p, err := ref.AccountPermissionRepository.FindByAccountId(accounts[key].Id)
		if err != nil {
			ref.Logger.Error(listLoggerPrefix+" Error when listing account permissions", err, nil)
			return nil, err
		}
		accounts[key].Permissions = p
	}

	return accounts, nil
}
