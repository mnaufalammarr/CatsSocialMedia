package response

import (
	"CatsSocialMedia/model/enum"
	"time"
)

type CatResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Race        enum.Race `json:"race"`
	Sex         enum.Sex  `json:"sex"`
	AgeInMonth  int       `json:"ageInMonth"`
	ImageURLs   []string  `json:"imageUrls"`
	Description string    `json:"description"`
	HasMatched  bool      `json:"hasMatched"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CatResponseString struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Race        string `json:"race"`
	Sex         string `json:"sex"`
	AgeInMonth  string `json:"ageInMonth"`
	ImageURLs   string `json:"imageUrls"`
	Description string `json:"description"`
	HasMatched  string `json:"hasMatched"`
	CreatedAt   string `json:"createdAt"`
}

type CreateCatResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
