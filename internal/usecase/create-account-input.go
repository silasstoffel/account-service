package usecase

type CreateAccountInput struct {
	Name      string
	LastName  string
	Email     string
	Phone     string
	HashedPwd string
}
