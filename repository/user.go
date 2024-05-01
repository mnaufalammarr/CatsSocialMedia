package repository

import (
	"CatsSocialMedia/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(user model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
}
type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user model.User) (model.User, error) {
	fmt.Println("repo")
	fmt.Println(user)
	_, err := r.db.Exec(context.Background(), "INSERT INTO users (email, name, password) VALUES ($1, $2, $3)", user.Email, user.Name, user.Password)
	fmt.Println(err)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.QueryRow(context.Background(), "SELECT id, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
