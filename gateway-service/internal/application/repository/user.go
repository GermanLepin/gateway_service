package repository

import (
	"context"
	"database/sql"

	"gateway-service/internal/application/dto"

	"github.com/google/uuid"
)

func (u *userRepository) CreateUser(ctx context.Context, user *dto.User) error {
	queryString := `
		insert into service.user 
		(id, name, surname, phone, email, password) 
		values ($1,$2,$3,$4,$5,$6)
	;`

	err := u.db.QueryRow(queryString, user.ID, user.Name, user.Surname, user.Phone, user.Email, user.Password)
	if err != nil {
		return err.Err()
	}

	return nil
}

func (u *userRepository) FetchUserByEmail(ctx context.Context, userEmail string) (user dto.User, err error) {
	queryString := `select * from service.user where email = $1;`

	err = u.db.QueryRow(queryString, userEmail).Scan(&user.ID, &user.Name, &user.Surname, &user.Phone, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *userRepository) FetchUserById(ctx context.Context, userID uuid.UUID) (user dto.User, err error) {
	queryString := `select * from service.user where id = $1;`

	err = u.db.QueryRow(queryString, userID).Scan(&user.ID, &user.Name, &user.Surname, &user.Phone, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *userRepository) DeleteUserByEmail(ctx context.Context, userEmail string) error {
	err := u.db.QueryRow("delete from service.user where email = $1;", userEmail)
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
