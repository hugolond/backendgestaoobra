package src

import (
	"backendgestaoobra/pkg"
	"database/sql"
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

func LoginHandler(c *gin.Context) {
	var req LoginRequest
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário ou senha inválidos"})
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"roles": user.Roles,
		"exp":   expirationTime.Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

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
