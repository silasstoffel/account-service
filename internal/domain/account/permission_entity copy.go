package domain

type ListPermissionInput struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ReadPermissionRepository interface {
	List(input ListPermissionInput) ([]*Permission, error)
}

type PermissionRepository interface {
	ReadPermissionRepository
}
