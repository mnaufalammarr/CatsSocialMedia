package request

type MatchApprovalRequest struct {
	MatchID int `json:"matchId" binding:"required"`
}
