package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	errorDomain "github.com/silasstoffel/account-service/internal/domain/exception"
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

func (repository *AccountRepository) Create(account accountDomain.Account) (accountDomain.Account, error) {
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
		return account, errorDomain.NewError(errorDomain.DbCommandError, "Error when creating account", err)
	}
	affected, _ := result.RowsAffected()
	log.Println(loggerPrefix, "Affected row", affected)
	log.Println(loggerPrefix, "Account created with id", account.Id)
	return account, nil
}

func (repository *AccountRepository) FindByEmail(email string) (accountDomain.Account, error) {
	log.Println(loggerPrefix, "Finding account by email")
	var account accountDomain.Account

	stmt := `SELECT id, name, last_name, email, phone, created_at, updated_at, active, full_name, COALESCE(hashed_pwd, '')
	         FROM accounts
	         WHERE email = $1`

	row := repository.Db.QueryRow(stmt, email)
	err := scanRow(row, &account)
	if err != nil {
		log.Println(loggerPrefix, "Error find account by email. Detail", err.Error())
		if err == sql.ErrNoRows {
			return account, errorDomain.NewError(accountDomain.AccountNotFound, "Account not found", nil)
		}
		return account, errorDomain.NewError(errorDomain.DbCommandError, "Error when finding account by e-mail", err)
	}

	log.Println(loggerPrefix, "Account found with id", account.Id)
	return account, nil
}

func (repository *AccountRepository) FindByPhone(phone string) (accountDomain.Account, error) {
	log.Println(loggerPrefix, "Finding account by phone")
	var account accountDomain.Account

	stmt := `SELECT id, name, last_name, email, phone, created_at, updated_at, active, full_name, COALESCE(hashed_pwd, '')
	         FROM accounts
	         WHERE phone = $1`

	row := repository.Db.QueryRow(stmt, phone)
	err := scanRow(row, &account)
	if err != nil {
		log.Println(loggerPrefix, "Error find account by phone. Detail", err.Error())
		if err == sql.ErrNoRows {
			log.Println(loggerPrefix, "There is no row to return")
			return account, errorDomain.NewError(accountDomain.AccountNotFound, "Account not found", nil)
		}
		return account, errorDomain.NewError(errorDomain.DbCommandError, "Error when finding account by e-mail", err)
	}

	log.Println(loggerPrefix, "Account found with id", account.Id)
	return account, nil
}

func (repository *AccountRepository) List(input accountDomain.ListAccountInput) ([]accountDomain.Account, error) {
	log.Println(loggerPrefix, "Listing accounts. Page:", input.Page, "Limit:", input.Limit)

	stmt := `SELECT id, name, last_name, email, phone, created_at, updated_at, active, full_name, COALESCE(hashed_pwd, '')
	         FROM accounts
	         LIMIT $1 OFFSET $2`
	page, limit := input.Page, input.Limit

	if page <= 1 {
		page = 1
	}

	if limit <= 1 {
		limit = 12
	}
	offset := (page - 1) * limit
	rows, err := repository.Db.Query(stmt, limit, offset)
	var accounts []accountDomain.Account
	if err != nil {
		log.Println(loggerPrefix, "error when execute command on database.", err.Error())
		return accounts, errorDomain.NewError(errorDomain.DbCommandError, "Error when listing accounts", err)
	}

	defer rows.Close()
	for rows.Next() {
		var account accountDomain.Account
		if err := scanRow(rows, &account); err != nil {
			log.Println(loggerPrefix, "error when scan result", err.Error())
			return accounts, errorDomain.NewError(errorDomain.DbCommandError, "Error when listing accounts", err)
		}
		log.Println(loggerPrefix, "account", account)
		accounts = append(accounts, account)
	}

	log.Println(loggerPrefix, "Listed accounts", "total", len(accounts))
	return accounts, nil
}

func (repository *AccountRepository) FindById(accountId string) (accountDomain.Account, error) {
	log.Println(loggerPrefix, "Finding account by id", accountId)
	var account accountDomain.Account

	stmt := `SELECT id, name, last_name, email, phone, created_at, updated_at, active, full_name, COALESCE(hashed_pwd, '')
	         FROM accounts
	         WHERE id = $1`

	row := repository.Db.QueryRow(stmt, accountId)

	if err := scanRow(row, &account); err != nil {
		log.Println(loggerPrefix, "Error when finding account by id.", err.Error())
		if err == sql.ErrNoRows {
			return account, errorDomain.NewError(accountDomain.AccountNotFound, "Account not found", nil)
		}
		return account, errorDomain.NewError(errorDomain.DbCommandError, "Error when finding account by id.", err)
	}

	log.Println(loggerPrefix, "Account found", account.Id)
	return account, nil
}

func scanRow(row interface{}, account *accountDomain.Account) error {
	switch r := row.(type) {
	case *sql.Row:
		return r.Scan(
			&account.Id,
			&account.Name,
			&account.LastName,
			&account.Email,
			&account.Phone,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.Active,
			&account.FullName,
			&account.HashedPwd,
		)
	case *sql.Rows:
		return r.Scan(
			&account.Id,
			&account.Name,
			&account.LastName,
			&account.Email,
			&account.Phone,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.Active,
			&account.FullName,
			&account.HashedPwd,
		)
	}
	return errorDomain.NewError(errorDomain.UnknownError, "An Unknown error happens", errors.New("ScanRow error is not sql.Row or sql.Rows"))
}
