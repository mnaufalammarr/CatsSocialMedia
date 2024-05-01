package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func ConnectToDatabase(urlDb string) (*pgx.Conn, error) {
	//urlDb := "postgres://postgres:postgres@localhost:5432/pustakaapi?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), urlDb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	return conn, err
}
