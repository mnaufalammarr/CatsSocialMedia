package CatsSocialMedia

import (
	"CatsSocialMedia/controller"
	"CatsSocialMedia/db"
	"CatsSocialMedia/repository"
	"CatsSocialMedia/service"
	"CatsSocialMedia/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"

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
	dbParams := os.Getenv("DB_PARAMS")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", dbUsername, dbPassword, dbHost, dbPort, dbName, dbParams)

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
	routerV1.POST("/signup", userController.Signup)
	routerV1.POST("/login", userController.SignIn)
}
