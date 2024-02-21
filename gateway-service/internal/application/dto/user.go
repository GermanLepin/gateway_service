package dto

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Password     string
	Email        string
	Phone        int
	UserType     string
	AccessToken  string
	RefreshToken string
}
