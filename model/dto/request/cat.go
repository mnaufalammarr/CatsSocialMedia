package request

import "CatsSocialMedia/model/enum"

type CatRequest struct {
	UserId      int       `json:"user_id"`
	Name        string    `json:"name" binding:"required,min=1,max=30"`
	Race        enum.Race `json:"race" binding:"required"`
	Sex         enum.Sex  `json:"sex" binding:"required"`
	AgeInMonths int       `json:"age_in_months" binding:"required,min=1,max=120082"`
	Description string    `json:"description" binding:"required,min=1,max=200"`
	ImageUrls   []string  `json:"image_urls" binding:"required,min=1,dive,url"`
}
