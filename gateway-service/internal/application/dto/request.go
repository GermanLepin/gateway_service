package dto

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     int    `json:"phone"`
	UserType  string `json:"user_type"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
