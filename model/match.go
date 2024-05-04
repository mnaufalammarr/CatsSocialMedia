package model

type Match struct {
	ID         int    `json:"id"`
	MatchCatID int    `json:"match_cat_id"`
	UserCatID  int    `json:"user_cat_id"`
	IssuedBy   int    `json:"issued_by"`
	IsAproved  bool   `json:"is_approved"`
	Message    string `json:"message"`
	IsMatched  bool   `json:"is_matched"`
	BaseModel
}
