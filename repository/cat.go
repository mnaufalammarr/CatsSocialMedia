package repository

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/response"
	"CatsSocialMedia/model/enum"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type CatRepository interface {
	FindAll(filterParams map[string]interface{}) ([]response.CatResponse, error)
	FindByUserID(i int) (model.Cat, error)
	FindByID(catID string) (model.Cat, error)
	FindByIDAndUserID(catID string, userID int) (model.Cat, error)
	Create(cat model.Cat) (model.Cat, error)
	Update(cat model.Cat) (model.Cat, error)
	Delete(catID string, userID int) error
}
type catRepository struct {
	db *pgx.Conn
}

func NewCatRepository(db *pgx.Conn) *catRepository {
	return &catRepository{db}
}

func (r *catRepository) FindAll(filterParams map[string]interface{}) ([]response.CatResponse, error) {
	query := "SELECT id, name, race, sex, age_in_months, description, image_urls FROM cats WHERE 1=1"
	var args []interface{}

	if catID, ok := filterParams["id"].(string); ok && catID != "" {
		query += " AND id = $1"
		args = append(args, catID)
	}

	if race, ok := filterParams["race"].(enum.Race); ok && race != "" {
		query += fmt.Sprintf(" AND race = '%s'", race)
	}

	if sexStr, ok := filterParams["sex"].(string); ok && sexStr != "" {
		sex := enum.Sex(sexStr)
		query += fmt.Sprintf(" AND sex = '%s'", sex)
	}

	if hasMatched, ok := filterParams["hasMatched"].(bool); ok {
		query += fmt.Sprintf(" AND has_match = %t", hasMatched)
	}

	if ageInMonth, ok := filterParams["ageInMonth"].(string); ok && ageInMonth != "" {
		var comparison string
		var age int
		if strings.HasPrefix(ageInMonth, ">") {
			comparison = ">"
			age, _ = strconv.Atoi(strings.TrimPrefix(ageInMonth, ">"))
		} else if strings.HasPrefix(ageInMonth, "<") {
			comparison = "<"
			age, _ = strconv.Atoi(strings.TrimPrefix(ageInMonth, "<"))
		} else {
			comparison = "="
			age, _ = strconv.Atoi(ageInMonth)
		}
		query += fmt.Sprintf(" AND age_in_months %s %d", comparison, age)
	}

	if owned, ok := filterParams["owned"].(bool); ok {
		// userId := 0
		userId := filterParams["userID"]
		if owned {
			query += fmt.Sprintf(" AND user_id = %d", userId)
		} else {
			query += fmt.Sprintf(" AND user_id != %d", userId)
		}
	}

	if search, ok := filterParams["search"].(string); ok && search != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", search)
	}

	if limit, ok := filterParams["limit"].(int); ok && limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
	}

	if offset, ok := filterParams["offset"].(int); ok && offset >= 0 {
		query += " OFFSET $3"
		args = append(args, offset)
	}
	fmt.Println(query)
	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []response.CatResponse
	for rows.Next() {
		var cat model.Cat
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonths, &cat.Description, &cat.ImageUrls)
		if err != nil {
			return nil, err
		}
		catResponse := response.CatResponse{
			ID:          cat.ID,
			Name:        cat.Name,
			Race:        cat.Race,
			Sex:         cat.Sex,
			AgeInMonth:  cat.AgeInMonths,
			ImageURLs:   cat.ImageUrls,
			Description: cat.Description,
			HasMatched:  cat.HasMatch,
		}
		cats = append(cats, catResponse)
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

func (r *catRepository) FindByUserID(i int) (model.Cat, error) {
	var cat model.Cat
	err := r.db.QueryRow(context.Background(), "SELECT id, name, race, sex, age_in_months, description, image_urls FROM cats WHERE user_id = $1", i).Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonths, &cat.Description, &cat.ImageUrls)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Cat{}, nil // Kucing tidak ditemukan, tidak ada error
		}
		return model.Cat{}, err // Error lainnya
	}
	return cat, nil
}

func (r *catRepository) FindByIDAndUserID(catID string, userID int) (model.Cat, error) {
	var cat model.Cat
	err := r.db.QueryRow(context.Background(), "SELECT id, name, race, sex, age_in_months, description, image_urls FROM cats WHERE user_id = $1 and id = $2", userID, catID).Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonths, &cat.Description, &cat.ImageUrls)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Cat{}, nil // Kucing tidak ditemukan, tidak ada error
		}
		return model.Cat{}, err // Error lainnya
	}
	return cat, nil
}

func (r *catRepository) Create(cat model.Cat) (model.Cat, error) {
	_, err := r.db.Exec(context.Background(), "INSERT INTO cats (name, race, sex, age_in_months, description, image_urls, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", cat.Name, cat.Race, cat.Sex, cat.AgeInMonths, cat.Description, cat.ImageUrls, cat.UserID)
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

func (r *catRepository) Delete(catID string, userID int) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM cats WHERE id = $1 and user_id = $2", catID, userID)
	if err != nil {
		return err
	}
	return nil
}
