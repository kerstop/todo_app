package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todolist/authentication"
	"todolist/database"
)

type TodoEntry struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Complete    bool   `json:"complete"`
}

func Todo(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		user, _, err := authentication.AuthenticateRequest(r)
		if err != nil {
			w.WriteHeader(401)
			http.SetCookie(w, &http.Cookie{
				Name:   "user_session",
				Value:  "",
				MaxAge: -1,
			})
			return
		}

		todo := make([]TodoEntry, 0)

		query := `
		SELECT todo_entries.id, descript, complete 
		FROM todo_entries 
		inner join users on users.id=todo_entries.user_id 
		where users.username=$1
		`
		rows, err := database.DbConnPool.Query(r.Context(), query, user)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Encountered Error: %v\n", err)
			return
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

		if err != nil {
			return
		}

		response := struct {
			TodoEntries []TodoEntry `json:"todoEntries"`
			Username    string      `json:"username"`
		}{
			TodoEntries: todo,
			Username:    user,
		}

		json.NewEncoder(w).Encode(response)
	}

	if r.Method == "POST" {

		var err error

		_, user_id, err := authentication.AuthenticateRequest(r)
		if err != nil {
			w.WriteHeader(401)
			fmt.Fprintf(w, "Must be logged in")
			return
		}

		var request struct {
			Description string `json:"description"`
		}
		json.NewDecoder(r.Body).Decode(&request)

		var todo_id int
		query := "INSERT INTO todo_entries (descript, complete, user_id) values ($1, false, $2) RETURNING id"
		err = database.DbConnPool.QueryRow(r.Context(), query, request.Description, user_id).Scan(&todo_id)
		if err != nil {
			w.WriteHeader(500)
			fmt.Printf("Error: %v\n", err)
			return
		}

		json.NewEncoder(w).Encode(todo_id)
	}

}

func ToggleDone(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		_, user_id, err := authentication.AuthenticateRequest(r)
		if err != nil {
			w.WriteHeader(401)
			return
		}

		var entry_id int
		err = json.NewDecoder(r.Body).Decode(&entry_id)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		var entry_owner_id int
		query := `SELECT user_id FROM todo_entries where id = $1`
		err = database.DbConnPool.QueryRow(r.Context(), query, entry_id).Scan(&entry_owner_id)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		if entry_owner_id != user_id {
			w.WriteHeader(403)
			return
		}

		query = `UPDATE todo_entries SET "complete" = NOT "complete" WHERE id = $1`
		_, err = database.DbConnPool.Exec(r.Context(), query, entry_id)
		if err != nil {
			w.WriteHeader(500)
			return
		}

	}

}
