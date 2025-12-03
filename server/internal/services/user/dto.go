package userService

type Register struct {
	Username string
	Name     string
	Password string
}

type FindUserResponse struct {
	Username string
	Name     string
	ID       int
}
