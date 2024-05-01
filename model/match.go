package model

type Match struct {
	ID         int  `json:"id"`
	MatchCatID Cat  `json:"match_cat_id"`
	UserCatID  Cat  `json:"user_cat_id"`
	HasMatched bool `json:"has_matched"`
	BaseModel
}
