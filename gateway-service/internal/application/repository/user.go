package repository

import (
	"context"
	"database/sql"

	"gateway-service/internal/application/dto"

	"github.com/google/uuid"
)

func (u *userRepository) CreateUser(ctx context.Context, user *dto.User) error {
	queryString := `
		insert into service.user(id, first_name, last_name, password, email, phone, user_type) 
		values ($1,$2,$3,$4,$5,$6, $7)
	;`

	err := u.db.QueryRow(queryString,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Email,
		user.Phone,
		user.UserType,
	)
	if err != nil {
		return err.Err()
	}

	return nil
}

func (u *userRepository) FetchUserByEmail(ctx context.Context, email string) (user dto.User, err error) {
	queryString := `select id, first_name, last_name, password, email, phone, user_type from service.user where email = $1;`

	err = u.db.QueryRow(queryString, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Email,
		&user.Phone,
		&user.UserType,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *userRepository) FetchUserById(ctx context.Context, userID uuid.UUID) (user dto.User, err error) {
	queryString := `select id, first_name, last_name, password, email, phone, user_type from service.user where id = $1;`

	err = u.db.QueryRow(queryString, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Email,
		&user.Phone,
		&user.UserType,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *userRepository) DeleteUserByEmail(ctx context.Context, email string) error {
	err := u.db.QueryRow("delete from service.user where email = $1;", email)
	if err != nil {
		return err.Err()
	}

	return nil
}

func (u *userRepository) DeleteUserById(ctx context.Context, userID uuid.UUID) error {
	err := u.db.QueryRow("delete from service.user where id = $1;", userID)
	if err != nil {
		return err.Err()
	}

	return nil
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}
