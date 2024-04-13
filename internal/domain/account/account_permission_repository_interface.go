package domain

type CreateAccountPermissionInput struct {
	AccountId    string `json:"accountId"`
	PermissionId string `json:"permissionId"`
}

type ListAccountPermissionInput struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type AccountPermissionReadRepository interface {
	FindByAccountId(accountId string) ([]AccountPermission, error)
}

type AccountPermissionWriteRepository interface {
	Create(data CreateAccountPermissionInput) error
	DeleteByAccount(accountId string) error
}

type AccountPermissionRepository interface {
	AccountPermissionReadRepository
	AccountPermissionWriteRepository
}
