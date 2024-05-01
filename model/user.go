package model

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Cats     []Cat  `json:"cats"` // One-to-many relationship with Cat
	BaseModel
}
