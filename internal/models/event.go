package models

type Event struct {
	Id int

	OwnerUserId int

	Label          string
	CoverPhotoPath *string

	// manually added fields
	CoverPhotoURL *string `orm:"ignore"`

	CreatedAt string
	UpdatedAt string

	Owner *User
}
