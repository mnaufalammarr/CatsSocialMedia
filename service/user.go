package service

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(signupRequest request.SignupRequest) (model.User, error)
	Login(loginRequest request.SignInRequest) (string, error)
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *userService {
	return &userService{repository}
}

func (s *userService) Create(signupRequest request.SignupRequest) (model.User, error) {
	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(signupRequest.Password), 10)

	if err != nil {
		return model.User{}, err
	}

	var isEmailExist bool = s.repository.EmailIsExist(signupRequest.Email)
	fmt.Println(isEmailExist)
	if isEmailExist {
		fmt.Println("hitted error email redudant")
		return model.User{}, errors.New("EMAIL ALREADY EXIST")
	}

	//save user
	user := model.User{
		Email:    signupRequest.Email,
		Name:     signupRequest.Name,
		Password: string(hash),
	}

	newUser, err := s.repository.Create(user)
	return newUser, err
}

func (s *userService) Login(loginRequest request.SignInRequest) (string, error) {
	//get user
	user, err := s.repository.FindByEmail(loginRequest.Email)

	if err != nil {
		return "", err
	} else if user.ID == 0 {
		return "", errors.New("Invalid email or password")
	}

	//compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

	if err != nil {
		return "", err
	}

	//sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
