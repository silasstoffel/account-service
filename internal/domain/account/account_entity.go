package domain

import (
	"fmt"
	"time"
)

type Account struct {
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	LastName    string              `json:"lastName"`
	FullName    string              `json:"fullName"`
	Email       string              `json:"email"`
	Phone       string              `json:"phone"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Active      bool                `json:"active"`
	HashedPwd   string              `json:"-"`
	Permissions []AccountPermission `json:"permissions,omitempty"`
}

func (account Account) ToDomain() Account {
	(account).BuildFullName()
	account.HashedPwd = ""

	return account
}

func (account *Account) IsEmpty() bool {
	return account.Id == ""
}

func (account *Account) BuildFullName() {
	if (account.Name != "" && account.LastName != "") && account.FullName == "" {
		account.FullName = fmt.Sprintf("%s %s", account.Name, account.LastName)
	}
}

// Error codes
const (
	AccountEmailAlreadyExists = "ACCOUNT_EMAIL_ALREADY_EXISTS"
	AccountPhoneAlreadyExists = "ACCOUNT_PHONE_ALREADY_EXISTS"
	AccountNotFound           = "ACCOUNT_NOT_FOUND"
)
