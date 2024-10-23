package ports

type AuthService interface {
	Register(username string, password string) error
	Login(username string, password string) (string, error)
}
