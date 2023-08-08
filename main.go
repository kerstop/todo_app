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

	http.HandleFunc("/api/todo", handlerWraper(routes.Todo))
	http.HandleFunc("/api/auth", handlerWraper(routes.Auth))
	http.HandleFunc("/api/auth/newUser", handlerWraper(routes.NewUser))

	err := http.ListenAndServe("127.0.0.1:80", nil)

	fmt.Printf("Todo app closed: %v", err)
}

func handlerWraper(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		handler(w, r)

	}
}
