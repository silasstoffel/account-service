package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/silasstoffel/account-service/internal/domain"
	"github.com/silasstoffel/account-service/internal/infra/helper"
)

const loggerPrefix = "[account-repository]"

type AccountRepository struct {
	Db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		Db: db,
	}
}

func (repository *AccountRepository) Create(account domain.Account) (domain.Account, error) {
	log.Println(loggerPrefix, "Creating account...")
	now := time.Now().UTC()

	account.Id = helper.NewULID()
	account.CreatedAt = now
	account.UpdatedAt = now

	stmt := `INSERT INTO accounts (id, name, last_name, email, phone, created_at, updated_at, active, full_name, hashed_pwd)
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	result, err := repository.Db.Exec(
		stmt,
		account.Id,
		account.Name,
		account.LastName,
		account.Email,
		account.Phone,
		account.CreatedAt,
		account.UpdatedAt,
		account.Active,
		account.FullName,
		account.HashedPwd,
	)

	if err != nil {
		return account, err
	}
	affected, _ := result.RowsAffected()
	log.Println(loggerPrefix, "Affected row", affected)
	log.Println(loggerPrefix, "Account created with id", account.Id)
	return account, nil
}
