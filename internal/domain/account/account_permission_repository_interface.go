package domain

type CreateAccountPermissionInput struct {
	AccountId    string `json:"accountId"`
	PermissionId string `json:"permissionId"`
}
type AccountPermissionRepository interface {
	FindByAccountId(accountId string) ([]AccountPermission, error)
	Create(data CreateAccountPermissionInput) error
	DeleteByAccount(accountId string) error
}
