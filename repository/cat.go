package repository

import (
	"CatsSocialMedia/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type CatRepository interface {
	FindAll() ([]model.Cat, error)
	FindByID(catID string) (model.Cat, error)
	Create(cat model.Cat) (model.Cat, error)
	Update(cat model.Cat) (model.Cat, error)
	Delete(catID string) error
}
type catRepository struct {
	db *pgx.Conn
}

func NewCatRepository(db *pgx.Conn) *catRepository {
	return &catRepository{db}
}

func (r *catRepository) FindAll() ([]model.Cat, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, name, race, sex, age_in_months, description, image_urls FROM cats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.Cat
	for rows.Next() {
		var cat model.Cat
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonths, &cat.Description, &cat.ImageUrls)
		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

func (r *catRepository) FindByID(catID string) (model.Cat, error) {
	var cat model.Cat
	err := r.db.QueryRow(context.Background(), "SELECT id, name, race, sex, age_in_months, description, image_urls FROM cats WHERE id = $1", catID).Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonths, &cat.Description, &cat.ImageUrls)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Cat{}, nil // Kucing tidak ditemukan, tidak ada error
		}
		return model.Cat{}, err // Error lainnya
	}
	return cat, nil
}

func (r *catRepository) Create(cat model.Cat) (model.Cat, error) {
	_, err := r.db.Exec(context.Background(), "INSERT INTO cats (name, race, sex, age_in_months, description, image_urls, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", cat.Name, cat.Race, cat.Sex, cat.AgeInMonths, cat.Description, cat.ImageUrls, 1)
	if err != nil {
		return model.Cat{}, err
	}
	return cat, nil
}

func (r *catRepository) Update(cat model.Cat) (model.Cat, error) {
	_, err := r.db.Exec(context.Background(), "UPDATE cats SET name = $1, race = $2, sex = $3, age_in_months = $4, description = $5, image_urls = $6 WHERE id = $7", cat.Name, cat.Race, cat.Sex, cat.AgeInMonths, cat.Description, cat.ImageUrls, cat.ID)
	if err != nil {
		return model.Cat{}, err
	}
	return cat, nil
}

func (r *catRepository) Delete(catID string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM cats WHERE id = $1", catID)
	if err != nil {
		return err
	}
	return nil
}
