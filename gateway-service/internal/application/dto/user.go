package dto

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Name     string
	Surname  string
	Phone    int
	Email    string
	JWTToken string
	Password string
}
