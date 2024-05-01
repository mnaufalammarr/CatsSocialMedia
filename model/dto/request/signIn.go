package request

type SignInRequest struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}
