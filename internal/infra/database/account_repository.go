package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	accountDomain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/exception"
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
	now := time.Now().UTC()

	account.Id = helper.NewULID()
	account.CreatedAt = now
	account.UpdatedAt = now

	stmt := `INSERT INTO accounts (id, name, last_name, email, phone, created_at, updated_at, active, full_name, hashed_pwd)
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := repository.Db.Exec(
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
		log.Println(loggerPrefix, "Error when creating account", err.Error())
		return account, exception.New(exception.DbCommandError, &err)
	}

	return account, nil
}

func (repository *AccountRepository) FindByEmail(email string) (accountDomain.Account, error) {
	var account accountDomain.Account

	stmt := `SELECT id, name, last_name, email, phone, created_at, updated_at, active, full_name, COALESCE(hashed_pwd, '')
	         FROM accounts
	         WHERE email = $1`

	row := repository.Db.QueryRow(stmt, email)
	err := scanRow(row, &account)
	if err != nil {
		log.Println(loggerPrefix, "Error find account by email. Detail", err.Error())
		if err == sql.ErrNoRows {
			return account, exception.New(exception.AccountNotFound, &err)
		}
		return account, exception.New(exception.DbCommandError, &err)
	}
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
			return account, exception.New(exception.AccountNotFound, &err)
		}
		return account, exception.New(exception.DbCommandError, &err)
	}

	log.Println(loggerPrefix, "Account found with id", account.Id)
	return account, nil
}

func (repository *AccountRepository) List(input accountDomain.ListAccountInput) ([]accountDomain.Account, error) {
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
		return accounts, exception.New(exception.DbCommandError, &err)
	}

	defer rows.Close()
	for rows.Next() {
		var account accountDomain.Account
		if err := scanRow(rows, &account); err != nil {
			log.Println(loggerPrefix, "error when scan result", err.Error())
			return accounts, exception.New(exception.DbCommandError, &err)
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (repository *AccountRepository) FindById(accountId string) (accountDomain.Account, error) {
	var account accountDomain.Account

	stmt := `SELECT id, name, last_name, email, phone, created_at, updated_at, active, full_name, COALESCE(hashed_pwd, '')
	         FROM accounts
	         WHERE id = $1`

	row := repository.Db.QueryRow(stmt, accountId)

	if err := scanRow(row, &account); err != nil {
		log.Println(loggerPrefix, "Error when finding account by id.", err.Error())
		if err == sql.ErrNoRows {
			return account, exception.New(exception.AccountNotFound, &err)
		}
		return account, exception.New(exception.DbCommandError, &err)
	}
	return account, nil
}

func (repository *AccountRepository) Update(id string, data accountDomain.Account) (accountDomain.Account, error) {
	account, err := (repository).FindById(id)
	if err != nil {
		log.Println(loggerPrefix, "Error when finding account", id, "Detail", err.Error())
		return account, err
	}

	var args []interface{}
	var updateFields []string
	argCount := 1

	if data.Name != "" {
		updateFields = append(updateFields, "name = $"+strconv.Itoa(argCount))
		args = append(args, data.Name)
		argCount++
		account.Name = data.Name
	}

	if data.LastName != "" {
		updateFields = append(updateFields, "last_name = $"+strconv.Itoa(argCount))
		args = append(args, data.LastName)
		argCount++
		account.LastName = data.LastName
	}

	if data.Email != "" {
		updateFields = append(updateFields, "email = $"+strconv.Itoa(argCount))
		args = append(args, data.Email)
		argCount++
		account.Email = data.Email
	}

	if data.Phone != "" {
		updateFields = append(updateFields, "phone = $"+strconv.Itoa(argCount))
		args = append(args, data.Phone)
		argCount++
		account.Phone = data.Phone
	}

	if data.HashedPwd != "" {
		updateFields = append(updateFields, "hashed_pwd = $"+strconv.Itoa(argCount))
		args = append(args, data.HashedPwd)
		argCount++
		account.HashedPwd = data.HashedPwd
	}

	updateFields = append(updateFields, "updated_at = $"+strconv.Itoa(argCount))
	now := time.Now().UTC()
	args = append(args, now)
	argCount++

	args = append(args, id)
	cols := strings.Join(updateFields, ", ")
	query := fmt.Sprintf("UPDATE accounts SET %s WHERE id = $%s", cols, strconv.Itoa(argCount))
	account.UpdatedAt = now

	log.Println(loggerPrefix, "Updating account", id, "query:", query)
	_, err = repository.Db.Exec(query, args...)
	if err != nil {
		log.Println(loggerPrefix, "Error when updating account", id, err.Error())
		return account, exception.New(exception.DbCommandError, &err)
	}

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
	m := errors.New("ScanRow error is not sql.Row or sql.Rows")
	return exception.New(exception.UnknownError, &m)
}
