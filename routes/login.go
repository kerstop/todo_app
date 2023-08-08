package routes

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"todolist/authentication"
	"todolist/database"

	"golang.org/x/crypto/bcrypt"
)

var hashing_difficulty = 12

type loginRequestForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type createUserRequestForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {


	switch r.Method {

	case "GET":
		user, err := authentication.AuthenticateRequest(r)
		if err != nil {
		w.Write([]byte(fmt.Sprintf("%v\n", err)))
		}

		w.Write([]byte(fmt.Sprintf("attempting to login as %s\n", user)))

	case "POST":

		var err error

		var loginRequest loginRequestForm
		// loginHeader := []byte(r.Header.Get("todo-list-auth"))
		fmt.Printf("%+v\n", r.Header)
		// loginHeader, err := base64.StdEncoding.DecodeString(r.Header.Get("todo-list-auth"))
		// if err != nil {
		// 	w.WriteHeader(400)
		// 	w.Write([]byte(fmt.Sprintf("Error: %v\n", err)))
		// 	return
		// }
		err = json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %v\n", err)))
			return
		}

		row := database.DbConnPool.QueryRow(r.Context(), `
			SELECT passwd_hash, salt FROM users WHERE username = $1;
		`, loginRequest.Username)

		var passwordHash, hashedSalt string
		row.Scan(&passwordHash, &hashedSalt)

		salt, err := base64.RawStdEncoding.DecodeString(hashedSalt)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("%v\n", err)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), append([]byte(loginRequest.Password), salt...))

		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %v\n", err)))
			return
		}

		tokenString, err := authentication.GenerateToken(loginRequest.Username)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			w.WriteHeader(500)
			return
		}

		cookie := &http.Cookie{
			Name: "user_session",
			Value: tokenString,
		}

		http.SetCookie(w, cookie)

	case "PUT":

		var createUserRequest createUserRequestForm
		err := json.NewDecoder(r.Body).Decode(&createUserRequest)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %v\n", err)))
			fmt.Printf("Error: %v\n", err)
			return
		}

		var salt [16]byte
		rand.Read(salt[:])

		salted_password := append([]byte(createUserRequest.Password), salt[:]...)

		hash, err := bcrypt.GenerateFromPassword(salted_password, hashing_difficulty)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %v\n", err)))
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		fmt.Printf("%s\n", hash)

		_ , err = database.DbConnPool.Exec(r.Context(), `
		INSERT INTO users (username, passwd_hash, salt)
		VALUES ($1, $2, $3);
		`, createUserRequest.Username, string(hash[:]), base64.RawStdEncoding.EncodeToString(salt[:]))
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %v\n", err)))
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

}
