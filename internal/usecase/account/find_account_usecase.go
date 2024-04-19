package usecase

import (
	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
)

func (ref *AccountUseCase) FindAccountUseCase(id string) (accountDomain.Account, error) {
	account, err := ref.AccountRepository.FindById(id)

	if err != nil {
		ref.Logger.Error(
			"[find-account-usecase] An error happens when find account",
			err,
			map[string]interface{}{"id": id},
		)
		return accountDomain.Account{}, err
	}

	p, _ := ref.AccountPermissionRepository.FindByAccountId(account.Id)
	account.Permissions = p

	return account, nil
}
