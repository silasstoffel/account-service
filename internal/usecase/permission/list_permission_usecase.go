package usecase

import (
	domain "github.com/silasstoffel/account-service/internal/domain/account"
	loggerContract "github.com/silasstoffel/account-service/internal/logger/contract"
)

type PermissionUseCase struct {
	PermissionRepository domain.PermissionRepository
	Logger               loggerContract.Logger
}

func NewPermissionUseCase(permissionRepository domain.PermissionRepository, logger loggerContract.Logger) *PermissionUseCase {
	return &PermissionUseCase{
		PermissionRepository: permissionRepository,
		Logger:               logger,
	}
}

func (ref *PermissionUseCase) ListPermissionUseCase(input domain.ListPermissionInput) ([]*domain.Permission, error) {
	p, err := ref.PermissionRepository.List(input)

	if err != nil {
		ref.Logger.Error("[list-permission-usecase] Error when listing accounts", err, nil)
		return nil, err
	}

	return p, nil
}
