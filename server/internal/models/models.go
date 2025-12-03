package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Name     string `json:"name"`
}

type Token struct {
	Id           int
	RefreshToken string
	UserId       int
}
