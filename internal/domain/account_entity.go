package domain

import (
	"fmt"
	"time"
)

type Account struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"lastName"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Active    bool      `json:"active"`
	HashedPwd string    `json:"-"`
}

func (account Account) ToDomain() Account {
	account.FullName = fmt.Sprintf("%s %s", account.Name, account.LastName)
	account.HashedPwd = ""

	return account
}
