package request

type SignInRequest struct {
	Email    string `binding:"required,email"` // Validation tags
	Password string `binding:"required"`
}
