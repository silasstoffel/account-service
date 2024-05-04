package domain

import "time"

type Permission struct {
	Id        string    `json:"id"`
	Scope     string    `json:"scope"`
	Active    time.Time `json:"active"`
	CreatedAt string    `json:"createdAt,omitempty"`
}

type AccountPermission struct {
	AccountId    string    `json:"accountId"`
	PermissionId string    `json:"permissionId"`
	Scope        string    `json:"scope"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	Active       bool      `json:"active,omitempty"`
}
