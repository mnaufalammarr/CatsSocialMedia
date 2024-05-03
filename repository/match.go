package repository

import (
	"CatsSocialMedia/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type MatchRepository interface {
	Create(match model.Match) (model.Match, error)
}

type matchRepository struct {
	db *pgx.Conn
}

func NewMatchRepository(db *pgx.Conn) *matchRepository {
	return &matchRepository{db}
}

func (r *matchRepository) Create(match model.Match) (model.Match, error) {
	fmt.Println(match)
	_, err := r.db.Exec(context.Background(), "INSERT INTO matchs (match_cat_id, user_cat_id, message, is_approved, issued_by) VALUES ($1,$2,$3,$4,$5)", match.MatchCatID, match.UserCatID, match.Message, nil, match.IssuedBy)
	if err != nil {
		return model.Match{}, err
	}
	return match, nil
}
