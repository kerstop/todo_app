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

	http.HandleFunc("/todo", handlerWraper(routes.Todo))
	http.HandleFunc("/auth", handlerWraper(routes.Auth))

	err := http.ListenAndServe("127.0.0.1:80", nil)

	fmt.Printf("Todo app closed: %v", err)
}

func handlerWraper(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")

		handler(w, r)

	}
}
