package CatsSocialMedia

import (
	"CatsSocialMedia/db"
	"CatsSocialMedia/utils"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

var conn *pgx.Conn

func init() {
	var err error

	utils.LoadEnvVariables()
	urlDb := os.Getenv("DATABASE_URL")
	conn, err = db.ConnectToDatabase(urlDb)

	if err != nil {
		log.Fatal("db connection failed")
	}
}

func main() {

}
