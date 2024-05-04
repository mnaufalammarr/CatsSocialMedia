package repository

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/response"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type MatchRepository interface {
	GetMatches(userId int) ([]response.MatchResponse, error)
	Create(match model.Match) (model.Match, error)
	MatchIsExist(matchId int) (model.Match, error)
	MatchApproval(matchId int, isApprove bool) (int, error)
	Delete(matchId int) error
	LatestMatch(catID string) (model.Match, error)
}

type matchRepository struct {
	db *pgx.Conn
}

func NewMatchRepository(db *pgx.Conn) *matchRepository {
	return &matchRepository{db}
}

func (r *matchRepository) Create(match model.Match) (model.Match, error) {
	fmt.Println(match)
	_, err := r.db.Exec(context.Background(), "INSERT INTO matchs (match_cat_id, user_cat_id, message, is_approved, issued_by, accepted_by) VALUES ($1,$2,$3,$4,$5,$6)", match.MatchCatID, match.UserCatID, match.Message, nil, match.IssuedBy, match.AcceptedBy)
	if err != nil {
		return model.Match{}, err
	}
	return match, nil
}

func (r *matchRepository) GetMatches(userId int) ([]response.MatchResponse, error) {
	query := `SELECT m.id AS match_id, 
				u_issuer.name AS issuer_name,
				u_issuer.email AS issuer_email,
				u_issuer.created_at AS issuer_created_at,
				c_match.id AS match_cat_id,
				c_match.name AS match_cat_name,
				c_match.race AS match_cat_race,
				c_match.sex AS match_cat_sex,
				c_match.age_in_months AS match_cat_age_in_month,
				c_match.description AS match_cat_description,
				c_match.image_urls AS match_cat_image_urls,
				c_match.has_match AS match_cat_has_match,
				c_match.created_at AS match_cat_created_at,
				c_user.id AS user_cat_id,
				c_user.name AS user_cat_name,
				c_user.race AS user_cat_race,
				c_user.sex AS user_cat_sex,
				c_user.age_in_months AS user_cat_age_in_month,
				c_user.description AS user_cat_description,
				c_user.image_urls AS user_cat_image_urls,
				c_user.has_match AS user_cat_has_match,
				c_user.created_at AS user_cat_created_at,
				m.message AS match_message,
				m.created_at AS match_created_at
			FROM matchs AS m
			JOIN users AS u_issuer ON m.issued_by = u_issuer.id 
			JOIN cats AS c_user ON m.user_cat_id = c_user.id 
			JOIN cats AS c_match ON m.match_cat_id = c_match.id
			WHERE (m.issued_by = $1 OR c_match.user_id = $1) AND m.is_matched = FALSE
	`
	rows, err := r.db.Query(context.Background(), query, userId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var matches []response.MatchResponse

	for rows.Next() {
		var match response.MatchResponse
		var issuer response.UserDetail
		var matchCat response.CatResponse
		var userCat response.CatResponse

		err := rows.Scan(
			&match.ID,
			&issuer.Name,
			&issuer.Email,
			&issuer.CreatedAt,
			&matchCat.ID,
			&matchCat.Name,
			&matchCat.Race,
			&matchCat.Sex,
			&matchCat.AgeInMonth,
			&matchCat.Description,
			&matchCat.ImageURLs,
			&matchCat.HasMatched,
			&matchCat.CreatedAt,
			&userCat.ID,
			&userCat.Name,
			&userCat.Race,
			&userCat.Sex,
			&userCat.AgeInMonth,
			&userCat.Description,
			&userCat.ImageURLs,
			&userCat.HasMatched,
			&userCat.CreatedAt,
			&match.Message,
			&match.CreatedAt,
		)

		if err != nil {
			fmt.Println(err)
			return []response.MatchResponse{}, err
		}

		fmt.Println(issuer)
		match.IssuedBy = issuer
		match.MatchCatDetail = matchCat
		match.UserCatDetail = userCat
		matches = append(matches, match)
	}

	return matches, nil
}

func (r *matchRepository) MatchIsExist(matchId int) (model.Match, error) {
	var match model.Match
	err := r.db.QueryRow(context.Background(), "SELECT id, match_cat_id, user_cat_id, message, issued_by, is_matched, created_at, updated_at, is_approved FROM matchs WHERE id = $1 LIMIT 1", matchId).Scan(&match.ID, &match.MatchCatID, &match.UserCatID, &match.Message, &match.IssuedBy, &match.IsMatched, &match.CreatedAt, &match.UpdatedAt, &match.IsAproved)

	fmt.Println(match)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Match{}, errors.New("MATCH IS NOT EXIST")
		}
	}
	return match, nil
}

func (r *matchRepository) MatchApproval(matchId int, isApprove bool) (int, error) {
	_, err := r.db.Exec(context.Background(), "UPDATE matchs SET is_approved = $1, is_matched = TRUE WHERE id = $2", isApprove, matchId)

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

func (r *matchRepository) LatestMatch(catID string) (model.Match, error) {
	var match model.Match
	err := r.db.QueryRow(context.Background(), "SELECT id, match_cat_id, user_cat_id, is_approved, message, issued_by, created_at, updated_at FROM matchs WHERE match_cat_id = $1 OR user_cat_id = $1 ORDER BY id DESC LIMIT 1", catID).Scan(&match.ID, &match.MatchCatID, &match.UserCatID, &match.IsAproved, &match.Message, &match.IssuedBy, &match.CreatedAt, &match.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Match{}, errors.New("MATCH IS NOT EXIST")
		}
	}
	return match, nil
}
