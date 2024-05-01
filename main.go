package CatsSocialMedia

import (
	"CatsSocialMedia/db"
	"CatsSocialMedia/utils"
	"fmt"
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

}
