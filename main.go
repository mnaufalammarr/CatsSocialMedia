package main

import (
	"CatsSocialMedia/controller"
	"CatsSocialMedia/db"
	"CatsSocialMedia/middleware"
	"CatsSocialMedia/repository"
	"CatsSocialMedia/service"
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
	catRepository := repository.NewCatRepository(conn)
	matchRepository := repository.NewMatchRepository(conn)

	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	matchService := service.NewMatchService(matchRepository, catRepository)
	matchController := controller.NewMatchController(matchService)

	catService := service.NewCatService(catRepository, matchService)
	catController := controller.NewCatController(catService)

	router := gin.Default()
	routerV1 := router.Group("/v1")

	//user
	userRouter := routerV1.Group("/user")
	userRouter.POST("/register", userController.Signup)
	userRouter.POST("/login", userController.SignIn)

	//cat
	catRouter := routerV1.Group("/cat", middleware.RequireAuth)
	catRouter.GET("/", catController.FindAll)
	catRouter.POST("/", catController.Create)
	catRouter.PUT("/:id", catController.Update)
	catRouter.GET("/:id", catController.FindByID)
	// catRouter.GET("/mine", catController.FindByUserID)
	catRouter.DELETE("/:id", catController.Delete)

	matchRouter := catRouter.Group("/match", middleware.RequireAuth)
	matchRouter.POST("/", matchController.Create)
	matchRouter.POST("/approve", matchController.Approve)
	matchRouter.POST("/reject", matchController.Reject)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
