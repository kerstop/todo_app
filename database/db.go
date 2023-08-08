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

	DbConnPool, err = pgxpool.New(c, "postgres://postgres:postgrespw@localhost:32769/todo")
	if err != nil {
		fmt.Printf("Encountered an error connecting to the database: %v\n", err)
		os.Exit(1)
	}

	_, err = DbConnPool.Exec(c,
		`CREATE TABLE IF NOT EXISTS todo_entries (
	    id serial,
    	descript text,
    	complete boolean,
    	PRIMARY KEY(id)
	)`)
	if err != nil {
		fmt.Printf("Encountered an error creating database: %v\n", err)
		os.Exit(1)
	}

}
