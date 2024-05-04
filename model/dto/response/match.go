package response

import (
	"time"
)

type MatchResponse struct {
	ID             int         `json:"id"`
	IssuedBy       UserDetail  `json:"issuedBy"`
	MatchCatDetail CatResponse `json:"matchCatDetail"`
	UserCatDetail  CatResponse `json:"userCatDetail"`
	Message        string      `json:"message"`
	CreatedAt      time.Time   `json:"createdAt"`
}

type UserDetail struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
