package model

//import "github.com/go-playground/validator/v10"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email,unique=users.email"` // Validation tags
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
	Cats     []Cat  `json:"cats"` // One-to-many relationship with Cat
	BaseModel
}
