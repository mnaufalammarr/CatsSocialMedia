package model

import "CatsSocialMedia/model/enum"

type Cat struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Race        enum.Race `json:"race"`
	Sex         enum.Sex  `json:"sex"`
	AgeInMonth  int       `json:"age_in_month"`
	Description string    `json:"description"`
	ImageUrls   []string  `json:"image_urls"`
	UserID      int       `json:"user_id"` // Foreign key referencing User
	HasMatch    bool      `json:"hasMatch"`
	BaseModel
}
