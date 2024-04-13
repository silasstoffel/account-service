package usecase

import (
	"log"

	domain "github.com/silasstoffel/account-service/internal/domain/account"
)

type PermissionUseCase struct {
	PermissionRepository domain.PermissionRepository
}

func NewPermissionUseCase(permissionRepository domain.PermissionRepository) *PermissionUseCase {
	return &PermissionUseCase{
		PermissionRepository: permissionRepository,
	}
}

func (ref *PermissionUseCase) ListPermissionUseCase(input domain.ListPermissionInput) ([]*domain.Permission, error) {
	p, err := ref.PermissionRepository.List(input)

	if err != nil {
		log.Println("[list-permission-usecase] Error when listing accounts", err)
		return nil, err
	}

	return p, nil
}
