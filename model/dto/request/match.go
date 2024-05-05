package request

type MatchRequest struct {
	MatchCatID string `json:"matchCatId" binding:"required"`
	UserCatID  string `json:"userCatId" binding:"required"`
	Message    string `json:"message" binding:"required"`
}
