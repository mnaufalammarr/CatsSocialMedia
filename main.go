package main

import (
	"CatsSocialMedia/controller"
	"CatsSocialMedia/db"
<<<<<<< HEAD
	"CatsSocialMedia/routes"
=======
	"CatsSocialMedia/repository"
	"CatsSocialMedia/service"
>>>>>>> master
	"CatsSocialMedia/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func init() {
	var err error

	utils.LoadEnvVariables()
	// urlDb := os.Getenv("DATABASE_URL")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	// dbParams := os.Getenv("DB_PARAMS")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)

	conn, err = db.ConnectToDatabase(dbURL)
	if err != nil {
		log.Fatal("db connection failed")
	}
}

func main() {
	userRepository := repository.NewUserRepository(conn)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	router := gin.Default()
	routerV1 := router.Group("/v1")
	routerV1.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
	routerV1.POST("/signup", userController.Signup)
	routerV1.POST("/login", userController.SignIn)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
