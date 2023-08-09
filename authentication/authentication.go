package authentication

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"todolist/database"

	"github.com/golang-jwt/jwt"
)

var secret [1024]byte

var ErrTooOld = errors.New("this token has expired")
var ErrNoToken = errors.New("no token was provided")
var ErrInvalid = errors.New("the provided token is not valid")
var ErrMalformed = errors.New("an error occured while parsing the token")

type CustomClaims struct {
	*jwt.StandardClaims
	Username string `json:"username"`
}

func init() {
	_, err := rand.Read(secret[:])
	if err != nil {
		fmt.Printf("Error occurred instantiating random secret: %v\n", err)
		os.Exit(-1)
	}
}

func AuthenticateRequest(r *http.Request) (string, int, error) {
	cookie, err := r.Cookie("user_session")
	if err != nil {
		return "", 0, ErrNoToken
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &CustomClaims{},  func(_ *jwt.Token) (interface{}, error) {return secret[:], nil})
	if err != nil {
		return "", 0, errors.Join(ErrInvalid, err)
	}

	if !token.Valid {
		return "", 0, ErrInvalid
	}
	username := token.Claims.(*CustomClaims).Username
	var id int

	database.DbConnPool.QueryRow(r.Context(), "SELECT id FROM users WHERE username=$1", username).Scan(&id)

	return username, id, nil

}

func GenerateToken(user string) (string, error) {

		token := jwt.New(jwt.GetSigningMethod("HS512"))

		expirationTime := time.Now().Add(time.Hour * 12)
		token.Claims = &CustomClaims{
			&jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
			user,
		}

		tokenString, err := token.SignedString(secret[:])
		
		return tokenString, err
}