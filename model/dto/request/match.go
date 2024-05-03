package request

type MatchRequest struct {
	MatchCatID int    `json:"match_cat_id" binding:"required"`
	UserCatID  int    `json:"user_cat_id" binding:"required"`
	Message    string `json:"message" binding:"required"`
}
