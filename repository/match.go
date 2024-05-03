package repository

import (
	"CatsSocialMedia/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type MatchRepository interface {
	Create(match model.Match) (model.Match, error)
	MatchIsExist(matchId int) (model.Match, error)
	MatchApproval(matchId int, isApprove bool) (int, error)
	Delete(matchId int) error
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

func (r *matchRepository) MatchIsExist(matchId int) (model.Match, error) {
	var match model.Match
	err := r.db.QueryRow(context.Background(), "SELECT id, match_cat_id, user_cat_id, is_approved, message, issued_by, created_at, updated_at FROM matchs WHERE id = $1 LIMIT 1", matchId).Scan(&match.ID, &match.MatchCatID, &match.UserCatID, &match.IsAproved, &match.Message, &match.IssuedBy, &match.CreatedAt, &match.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Match{}, errors.New("MATCH IS NOT EXIST")
		}
	}
	return match, nil
}

func (r *matchRepository) MatchApproval(matchId int, isApprove bool) (int, error) {
	_, err := r.db.Exec(context.Background(), "UPDATE matchs SET is_approved = $1 WHERE id = $2", isApprove, matchId)

	if err != nil {
		return matchId, err
	}
	return matchId, nil
}

func (r *matchRepository) Delete(matchId int) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM matchs WHERE id = $1", matchId)
	if err != nil {
		return err
	}
	return nil
}
