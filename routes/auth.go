package routes

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"todolist/authentication"
	"todolist/database"

	"github.com/jackc/pgx/v5"
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

func Auth(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		user, _, err := authentication.AuthenticateRequest(r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("%v\n", err)))
		}

		w.Write([]byte(fmt.Sprintf("attempting to login as %s\n", user)))

	case "POST":

		var err error

		var loginRequest loginRequestForm

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
		err = row.Scan(&passwordHash, &hashedSalt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				w.WriteHeader(400)
				w.Write([]byte("Invalid username or password"))
				return
			}
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %v\n", err)))
			return
		}

		salt, err := base64.RawStdEncoding.DecodeString(hashedSalt)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("%v\n", err)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), append([]byte(loginRequest.Password), salt...))
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				w.WriteHeader(400)
				w.Write([]byte("Invalid username or password"))
				return
			}
			w.WriteHeader(500)
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
			Name:  "user_session",
			Value: tokenString,
		}

		http.SetCookie(w, cookie)
	}
}

func NewUser(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":

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

		_, err = database.DbConnPool.Exec(r.Context(), `
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
