package src

import (
	"backendgestaoobra/pkg"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RefreshTokenHandler(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token ausente"})
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token inválido"})
		return
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(4 * time.Hour))
	claims.IssuedAt = jwt.NewNumericDate(time.Now().UTC())

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := newToken.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao gerar novo token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	currentTime := time.Now()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de login inválidos"})
		return
	}

	db, err := pkg.OpenConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro na conexão com banco de dados"})
		return
	}
	defer db.Close()

	var user User
	sqlQuery := `SELECT id, username, email, password, active, roles, departament, emailmanager FROM public."User" WHERE email = $1`
	err = db.QueryRow(sqlQuery, req.Email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.Roles,
		&user.Departament,
		&user.EmailManager,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário ou senha inválidos"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar usuário"})
		return
	}

	if !user.Active {
		c.JSON(http.StatusForbidden, gin.H{"error": "Usuário inativo"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | A1 - Auth Login - User: " + user.Username + " Status - Usuário ou senha inválidos")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário ou senha inválidos"})
		return
	}

	expirationTime := time.Now().Add(4 * time.Hour)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    user.Email,
		Subject:   user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | A1 - Auth Login - User: " + user.Username + " Status 200 OK!")
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"roles":    user.Roles,
		},
	})
}
