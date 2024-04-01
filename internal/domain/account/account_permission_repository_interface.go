package domain

type AccountPermissionRepository interface {
	FindByAccountId(accountId string) ([]AccountPermission, error)
	Create(data AccountPermission) error
	DeleteByAccount(accountId string) error
}
