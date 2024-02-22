package dto

import (
	"github.com/google/uuid"
)

type CreateUserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     int       `json:"phone"`
	UserType  string    `json:"user_type"`
	Message   string    `json:"message"`
}

type LoginResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone        int       `json:"phone"`
	UserType     string    `json:"user_type"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Message      string    `json:"message"`
}

type FetchUserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     int       `json:"phone"`
	UserType  string    `json:"user_type"`
	Message   string    `json:"message"`
}

type DeleteUserResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
