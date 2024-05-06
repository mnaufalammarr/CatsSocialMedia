package repository

import (
	"CatsSocialMedia/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(user model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
	EmailIsExist(email string) bool
	FindById(id int) (model.User, error)
}
type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, nil // Kucing tidak ditemukan, tidak ada error
		}
		return model.User{}, err // Error lainnya
	}
	return user, nil
}

func (r *userRepository) EmailIsExist(email string) bool {
	var exist string
	err := r.db.QueryRow(context.Background(), "SELECT email FROM users WHERE email = $1 LIMIT 1", email).Scan(&exist)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false
		}
	}
	return true
}

func (r *userRepository) FindById(id int) (model.User, error) {
	var user model.User
	err := r.db.QueryRow(context.Background(), "SELECT id, email, name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, nil // Kucing tidak ditemukan, tidak ada error
		}
		return model.User{}, err // Error lainnya
	}

	return user, nil
}
