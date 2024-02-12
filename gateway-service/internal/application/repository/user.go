package repository

import (
	"context"
	"database/sql"

	"gateway-service/internal/application/dto"

	"github.com/google/uuid"
)

func (u *userRepository) CreateUserById(ctx context.Context, user dto.CretaeUserRequest) error {
	err := u.db.QueryRow("insert into service.users (id, name) values ($1,$2);", user.ID, user.Name)
	if err != nil {
		return err.Err()
	}

	return nil
}

func (u *userRepository) DeleteUserById(ctx context.Context, userID uuid.UUID) error {
	err := u.db.QueryRow("select * from service.users where id = $1;", userID)
	if err != nil {
		return err.Err()
	}

	err = u.db.QueryRow("delete from service.users where id = $1;", userID)
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
