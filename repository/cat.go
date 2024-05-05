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
	"time"

	"github.com/jackc/pgx/v5"
)

type CatRepository interface {
	FindAll(filterParams map[string]interface{}) ([]response.CatResponse, error)
	FindByUserID(i int) (model.Cat, error)
	FindByID(catID string) (model.Cat, error)
	FindByIDAndUserID(catID string, userID int) (model.Cat, error)
	Create(cat model.Cat) (response.CreateCatResponse, error)
	Update(cat model.Cat) (model.Cat, error)
	Delete(catID string, userID int) error
	UpdateHasMatch(id string, isHasMatch bool) (string, error)
}
type catRepository struct {
	db *pgx.Conn
}

func NewCatRepository(db *pgx.Conn) *catRepository {
	return &catRepository{db}
}

func (r *catRepository) FindAll(filterParams map[string]interface{}) ([]response.CatResponse, error) {
	query := "SELECT id, name, race, sex, age_in_month, description, image_urls, created_at FROM cats WHERE 1=1"
	var args []interface{}
	fmt.Println(args...)
	num := 1

	if catID, ok := filterParams["id"].(string); ok && catID != "" {
		query += " AND id = $" + strconv.Itoa(num)
		num++
		args = append(args, catID)
	}

	if race, ok := filterParams["race"].(string); ok && race != "" {
		query += " AND race = $" + strconv.Itoa(num)
		args = append(args, race)
		num++
	}

	if sexStr, ok := filterParams["sex"].(string); ok && sexStr != "" {
		sex := enum.Sex(sexStr)
		query += " AND sex = $" + strconv.Itoa(num)
		args = append(args, sex)
		num++
	}

	if hasMatched, ok := filterParams["hasMatched"].(bool); ok {
		query += " AND has_match = $" + strconv.Itoa(num)
		args = append(args, hasMatched)
		num++
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
		query += fmt.Sprintf(" AND age_in_month %s %d", comparison, age)
	}

	if owned, ok := filterParams["owned"].(bool); ok {
		// userId := 0
		userId := filterParams["userID"]
		if owned {
			query += " AND user_id = $" + strconv.Itoa(num)
		} else {
			query += " AND user_id  = $" + strconv.Itoa(num)
		}
		args = append(args, userId)
		num++
	}

	if search, ok := filterParams["search"].(string); ok && search != "" {
		query += " AND name ILIKE '%%$" + strconv.Itoa(num) + "%%'"
		num++
	}

	query += (" ORDER BY id DESC ")

	if limit, ok := filterParams["limit"].(int); ok && limit > 0 {
		query += " LIMIT " + strconv.Itoa(limit)
	}

	if offset, ok := filterParams["offset"].(int); ok && offset >= 0 {
		query += " OFFSET  " + strconv.Itoa(offset)
	}
	fmt.Println(query)
	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var cats []response.CatResponse
	for rows.Next() {
		var cat model.Cat
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.Description, &cat.ImageUrls, &cat.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		// createdAtISO8601 := cat.CreatedAt.Format(time.RFC3339)
		catResponse := response.CatResponse{
			ID:          strconv.Itoa(cat.ID),
			Name:        cat.Name,
			Race:        cat.Race,
			Sex:         cat.Sex,
			AgeInMonth:  cat.AgeInMonth,
			ImageURLs:   cat.ImageUrls,
			Description: cat.Description,
			HasMatched:  cat.HasMatch,
			CreatedAt:   cat.CreatedAt,
		}
		cats = append(cats, catResponse)
	}
	if cats == nil {
		cats = make([]response.CatResponse, 0)
	}
	return cats, nil
}

func (r *catRepository) FindByID(catID string) (model.Cat, error) {
	var cat model.Cat
	err := r.db.QueryRow(context.Background(), "SELECT id, name, race, sex, age_in_month, has_match, description, image_urls, user_id FROM cats WHERE id = $1", catID).Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.HasMatch, &cat.Description, &cat.ImageUrls, &cat.UserID)
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
	err := r.db.QueryRow(context.Background(), "SELECT id, name, race, sex, age_in_month, description, image_urls FROM cats WHERE user_id = $1", i).Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.Description, &cat.ImageUrls)
	fmt.Println(err)
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
	err := r.db.QueryRow(context.Background(), "SELECT id, name, race, sex, age_in_month, description, image_urls FROM cats WHERE user_id = $1 and id = $2", userID, catID).Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.Description, &cat.ImageUrls)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Cat{}, nil // Kucing tidak ditemukan, tidak ada error
		}
		return model.Cat{}, err // Error lainnya
	}
	return cat, nil
}

func (r *catRepository) Create(cat model.Cat) (response.CreateCatResponse, error) {
	var id string
	var createdAt time.Time
	err := r.db.QueryRow(context.Background(), "INSERT INTO cats (name, race, sex, age_in_month, description, image_urls, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at", cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, cat.ImageUrls, cat.UserID).Scan(&id, &createdAt)
	if err != nil {
		return response.CreateCatResponse{}, err
	}

	// Konversi waktu pembuatan ke format ISO 8601
	// createdAtISO8601 := createdAt.Format(time.RFC3339)

	// Buat respons yang akan dikirimkan kembali
	response := response.CreateCatResponse{
		ID:        id,
		CreatedAt: createdAt,
	}

	return response, nil
}

func (r *catRepository) Update(cat model.Cat) (model.Cat, error) {
	_, err := r.db.Exec(context.Background(), "UPDATE cats SET name = $1, race = $2, sex = $3, age_in_month = $4, description = $5, image_urls = $6 WHERE id = $7", cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, cat.ImageUrls, cat.ID)
	if err != nil {
		return model.Cat{}, err
	}
	fmt.Println("cat updeted")
	return cat, nil
}

func (r *catRepository) UpdateHasMatch(id string, isHasMatch bool) (string, error) {
	_, err := r.db.Exec(context.Background(), "UPDATE cats SET has_match = $1 WHERE id = $2", isHasMatch, id)
	if err != nil {
		return id, err
	}
	fmt.Println("cat updated")
	return id, nil
}

func (r *catRepository) Delete(catID string, userID int) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM cats WHERE id = $1 and user_id = $2", catID, userID)
	if err != nil {
		return err
	}
	return nil
}
