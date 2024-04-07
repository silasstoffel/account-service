package domain

type ListAccountInput struct {
	Page  int
	Limit int
}

type ReadOneAccountRepository interface {
	FindById(accountId string) (Account, error)
}

type ReadAccountRepository interface {
	ReadOneAccountRepository
	FindByEmail(email string) (Account, error)
	FindByPhone(phone string) (Account, error)
	List(input ListAccountInput) ([]Account, error)
}

type WriteAccountRepository interface {
	Create(account Account) (Account, error)
	Update(id string, data Account) (Account, error)
}

type AccountRepository interface {
	WriteAccountRepository
	ReadAccountRepository
}
