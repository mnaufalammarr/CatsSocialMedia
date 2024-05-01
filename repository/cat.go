package repository

import (
	"CatsSocialMedia/model"
	"context"

	"github.com/jackc/pgx/v5"
)

type CatRepository interface {
	Create(cat model.Cat) (model.Cat, error)
}
type catRepository struct {
	db *pgx.Conn
}

func NewCatRepository(db *pgx.Conn) *catRepository {
	return &catRepository{db}
}

func (r *catRepository) Create(cat model.Cat) (model.Cat, error) {
	_, err := r.db.Exec(context.Background(), "INSERT INTO cats (name, race, sex, age_in_month, description, image_urls, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", cat.Name, cat.Race, cat.Sex, cat.AgeInMonths, cat.Description, cat.ImageUrls, 1)
	if err != nil {
		return model.Cat{}, err
	}
	return cat, nil
}
