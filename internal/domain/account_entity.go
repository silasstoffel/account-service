package domain

import (
	"fmt"
	"time"
)

type Account struct {
	Id        string
	Name      string
	LastName  string
	fullName  string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
	HashedPwd string
}

func (account Account) ToDomain() Account {
	account.fullName = fmt.Sprintf("%s %s", account.Name, account.LastName)
	account.HashedPwd = ""

	return account
}
