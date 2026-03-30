package entity

import "time"

type User struct {
	ID          string
	FirstName   string
	LastName    string
	Username    string
	Email       string
	PhoneNumber string
	Status      string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
