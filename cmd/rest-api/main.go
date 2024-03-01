package main

import (
	"github.com/silasstoffel/account-service/internal/infra/database"
	"github.com/silasstoffel/account-service/internal/usecase"
)

func main() {

	accountRepository := database.NewAccountRepository()

	createAccount := usecase.CreateAccount{
		AccountRepository: accountRepository,
	}

	createAccount.CreateAccountUseCase(usecase.CreateAccountInput{
		Name:      "Silas",
		LastName:  "Stoffel",
		Email:     "silasstofel@gmail.com",
		Phone:     "123",
		HashedPwd: "XPTOABC",
	})
}
