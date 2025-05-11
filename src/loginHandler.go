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

	claims := jwt.MapClaims{} // aqui, NÃO é ponteiro
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token inválido"})
		return
	}

	// Atualizar expiração
	claims["exp"] = time.Now().Add(4 * time.Hour).Unix()
	claims["iat"] = time.Now().Unix()

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
	sqlQuery := `SELECT id, username, email, password, active, roles, departament, emailmanager, account_id FROM public."User" WHERE email = $1`
	err = db.QueryRow(sqlQuery, req.Email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.Roles,
		&user.Departament,
		&user.EmailManager,
		&user.AccountID,
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

	// Verificar status da conta (account)
	var accountStatus bool
	sqlAccQuery := `SELECT status FROM obra.account WHERE id = $1`
	err = db.QueryRow(sqlAccQuery, user.AccountID).Scan(&accountStatus)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Conta não encontrada"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar status da conta"})
		return
	}

	if !accountStatus {
		c.JSON(http.StatusForbidden, gin.H{"error": "Conta inativa! Por favor acionar o Administrador da página"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | A1 - Auth Login - User: " + user.Username + " Status - Usuário ou senha inválidos")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário ou senha inválidos"})
		return
	}

	expirationTime := time.Now().Add(4 * time.Hour)

	claims := jwt.MapClaims{
		"sub":        user.ID,
		"email":      user.Email,
		"account_id": user.AccountID,
		"exp":        expirationTime.Unix(),
		"iat":        time.Now().Unix(),
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
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"roles":      user.Roles,
			"account_id": user.AccountID,
		},
	})
}
