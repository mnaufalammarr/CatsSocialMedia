package request

type MatchRequest struct {
	MatchCatID int    `json:"matchCatId" binding:"required"`
	UserCatID  int    `json:"userCatId" binding:"required"`
	Message    string `json:"message" binding:"required"`
}
