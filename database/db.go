package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func ConnectToDB() *pgx.Conn {
	databaseURL := "database_url"

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}
