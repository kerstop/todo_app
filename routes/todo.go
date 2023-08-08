package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todolist/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoEntry struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Complete    bool   `json:"complete"`
}

func Todo(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		todo_list, err := get_todo(database.DbConnPool, r)
		if err != nil {
			fmt.Fprintf(w, "Encountered Error: %v\n", err)
			return
		}

		json.NewEncoder(w).Encode(todo_list)

	}

	if r.Method == "POST" {
		id, err := post_todo(w, database.DbConnPool, r)

		if err != nil {
			w.WriteHeader(500)
			return
		}

		json.NewEncoder(w).Encode(id)
	}

}

func get_todo(conn *pgxpool.Pool, r *http.Request) ([]TodoEntry, error) {

	todo := make([]TodoEntry, 0)

	rows, err := conn.Query(r.Context(), "SELECT id, descript, complete FROM todo_entries")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := TodoEntry{}
		err = rows.Scan(&t.Id, &t.Description, &t.Complete)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
		todo = append(todo, t)
	}

	if rows.Err() != nil {
		fmt.Printf("ERROR: %v\n", rows.Err())
	}

	rows.Close()

	return todo, nil
}

func post_todo(w http.ResponseWriter, conn *pgxpool.Pool, r *http.Request) (id int, err error) {
	var request struct {
		Description string `json:"description"`
	}
	json.NewDecoder(r.Body).Decode(&request)

	fmt.Printf("INFO: description is `%v`\n", request.Description)

	err = conn.QueryRow(r.Context(), "INSERT INTO todo_entries (descript, complete) values ($1, false) RETURNING id", request.Description).Scan(&id)

	if err != nil {
		return
	}

	return id, err
}

