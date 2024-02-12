package repository

import (
	"context"
	"database/sql"

	"gateway-service/internal/application/dto"

	"github.com/google/uuid"
)

func (u *userRepository) CreateUser(ctx context.Context, user *dto.User) error {
	err := u.db.QueryRow("insert into service.user (id, name, email, password) values ($1,$2,$3,$4);", user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		return err.Err()
	}

	return nil
}
func (u *userRepository) FetchUserById(ctx context.Context, userID uuid.UUID) (user dto.User, err error) {
	err = u.db.QueryRow("select * from service.user where id = $1;", userID).Scan(&user.ID, &user.Name)
	if err != nil {
		return user, err
	}

	return user, nil
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
