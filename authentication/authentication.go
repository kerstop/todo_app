package authentication

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

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

func AuthenticateRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("user_session")
	if err != nil {
		return "", ErrNoToken
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &CustomClaims{},  func(_ *jwt.Token) (interface{}, error) {return secret[:], nil})
	if err != nil {
		return "", errors.Join(ErrInvalid, err)
	}

	if !token.Valid {
		return "", ErrInvalid
	}

	username := token.Claims.(*CustomClaims).Username

	return username, nil

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