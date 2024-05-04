package request

type SignupRequest struct {
	Email    string `binding:"required,email"` // Validation tags
	Name     string `binding:"required,min=5,max=50"`
	Password string `binding:"required,min=5,max=15"`
}
