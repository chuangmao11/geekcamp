package domain

type User struct {
	Id              int64
	Email           string
	Password        string
	ConfirmPassword string
	NickName        string
	Birthday        string
	Info            string
}
