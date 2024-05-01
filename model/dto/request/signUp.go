package request

type SignupRequest struct {
	Email    string `binding:"required"`
	Name     string `binding:"required"`
	Password string `binding:"required"`
}
