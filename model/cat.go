package model

type Cat struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Race        string   `json:"race"`
	Sex         bool     `json:"sex"`
	AgeInMonths int      `json:"age_in_months"`
	Description string   `json:"description"`
	ImageUrls   []string `json:"image_urls"`
	UserID      int      `json:"user_id"` // Foreign key referencing User
	HasMatch    bool     `json:"hasMatch"`
	BaseModel
}
