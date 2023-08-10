package database

import (

	"github.com/jackc/pgx/v5/pgxpool"

	"context"
	"fmt"
	"os"
)

var DbConnPool *pgxpool.Pool

func Connect(c context.Context) {

	var err error

	port_number := os.Getenv("POSGRES_PORT")
	if len(port_number) == 0 {
		port_number = "32769"
	}

	DbConnPool, err = pgxpool.New(c, fmt.Sprintf("postgres://postgres:postgrespw@localhost:%s/todo", port_number))
	if err != nil {
		fmt.Printf("Encountered an error connecting to the database: %v\n", err)
		os.Exit(1)
	}

}
