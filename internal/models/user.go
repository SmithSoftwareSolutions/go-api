package models

type User struct {
	Id int

	Email        string
	PasswordHash string

	CreatedAt string
	UpdatedAt string

	Events *[]Event
}
