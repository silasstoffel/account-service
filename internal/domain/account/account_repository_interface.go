package domain

type ListAccountInput struct {
	Page  int
	Limit int
}

type AccountRepository interface {
	List(input ListAccountInput) ([]Account, error)
	FindById(accountId string) (Account, error)
	Create(account Account) (Account, error)
	FindByEmail(email string) (Account, error)
	FindByPhone(phone string) (Account, error)
	Update(id string, data Account) (Account, error)
}
