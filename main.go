package main

import (
	"context"
	"todolist/database"
	"todolist/routes"

	"fmt"
	"net/http"
)

func main() {

	database.Connect(context.Background())
	defer database.DbConnPool.Close()

	http.HandleFunc("/api/todo", routes.Todo)
	http.HandleFunc("/api/todo/toggleDone", routes.ToggleDone)
	http.HandleFunc("/api/auth", routes.Auth)
	http.HandleFunc("/api/auth/newUser", routes.NewUser)

	err := http.ListenAndServe("127.0.0.1:80", nil)

	fmt.Printf("Todo app closed: %v", err)
}

