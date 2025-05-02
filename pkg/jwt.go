package pkg

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	//"github.com/HunCoding/meu-primeiro-crud-go/src/configuration/rest_err"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("chave_do_pernam")

type userDomain struct {
	id       string
	email    string
	password string
}

func GenerateToken(username string) (string, error) {
	tknID, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("erro token")
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Issuer:    username,
		IssuedAt:  time.Now().UTC().Unix(),
		Id:        tknID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func AuthMiddleware(tokenString string) string {
	// Valida o token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		fmt.Println("erro token")
		return ""
	}
	return "ok"
}
