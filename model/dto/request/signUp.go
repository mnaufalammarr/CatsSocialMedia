package request

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email,unique=users.email"` // Validation tags
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}
