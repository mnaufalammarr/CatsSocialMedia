package repository

import (
	"CatsSocialMedia/model"
	"context"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(user model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
	//IsEmailUsed(email string) (bool, error)
}
type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user model.User) (model.User, error) {
	_, err := r.db.Exec(context.Background(), "INSERT INTO users (email, name, password) VALUES ($1, $2, $3)", user.Email, user.Name, user.Password)
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

//
//func (r *userRepository) IsEmailUnique(email string) (bool, error) {
//	// Simulate database check (replace with actual database interaction)
//	var count int
//	err := r.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
//	if err != nil {
//		return false, err
//	}
//	return count == 0, nil // True if no matching email found
//}
