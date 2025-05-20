package src

import (
	"backendgestaoobra/pkg"
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
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

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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
		"id":         user.ID,
		"username":   user.Username,
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

func ForgotPassword(c *gin.Context) {
	tokenMailJet := os.Getenv("TOKEN_MAIL_JET")
	var req ForgotPasswordRequest
	currentTime := time.Now()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados reset de senha inválido"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Erro a redefinir senha"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro a redefinir senha"})
		return
	}

	if !user.Active {
		c.JSON(http.StatusForbidden, gin.H{"error": "Erro a redefinir senha"})
		return
	}

	// Verificar status da conta (account)
	var accountStatus bool
	sqlAccQuery := `SELECT status FROM obra.account WHERE id = $1`
	err = db.QueryRow(sqlAccQuery, user.AccountID).Scan(&accountStatus)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Erro a redefinir senha"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro a redefinir senha"})
		return
	}

	if !accountStatus {
		c.JSON(http.StatusForbidden, gin.H{"error": "Conta inativa! Por favor acionar o Administrador da página"})
		return
	}

	var existingToken string
	sqlTokenQuery := `
		SELECT token FROM public.password_reset_tokens 
		WHERE user_id = $1 AND used = false AND expires_at > NOW()
		LIMIT 1
	`
	err = db.QueryRow(sqlTokenQuery, user.ID).Scan(&existingToken)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar token"})
		return
	}

	if err == sql.ErrNoRows {
		// Não existe token válido, então cria um novo
		newToken, err := generateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
			return
		}

		expiration := time.Now().Add(20 * time.Minute)

		insertQuery := `
			INSERT INTO public.password_reset_tokens (user_id, token, created_at, expires_at, used) 
			VALUES ($1, $2, $3, $4, false)
		`
		_, err = db.Exec(insertQuery, user.ID, newToken, time.Now(), expiration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar token"})
			return
		}

		mailPayload := map[string]interface{}{
			"Messages": []map[string]interface{}{
				{
					"From": map[string]string{
						"Email": "noreply@gestaoobrafacil.com.br",
						"Name":  "Contato Obra Fácil",
					},
					"To": []map[string]string{
						{
							"Email": user.Email,
							"Name":  user.Username, // ou o nome real, se tiver
						},
					},
					"TemplateID":       7001820,
					"TemplateLanguage": true,
					"Subject":          "Redefina sua senha no Gestor Obra Fácil",
					"Variables": map[string]interface{}{
						"confirmation_link": "https://www.gestaoobrafacil.com.br/reset-password?token=" + newToken,
						"name":              user.Username,
						"defaultValue":      "",
						"isDate":            false,
						"type":              "data",
					},
				},
			},
		}

		jsonData, err := json.Marshal(mailPayload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao preparar o e-mail"})
			return
		}

		reqEmail, err := http.NewRequest("POST", "https://api.mailjet.com/v3.1/send", bytes.NewBuffer(jsonData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar requisição de e-mail"})
			return
		}

		reqEmail.Header.Add("Content-Type", "application/json")
		reqEmail.Header.Add("Authorization", "Basic "+tokenMailJet)

		client := &http.Client{}
		respEmail, err := client.Do(reqEmail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar e-mail"})
			return
		}
		defer respEmail.Body.Close()

		if respEmail.StatusCode >= 400 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro no envio de e-mail"})
			return
		}
	}

	fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | A1 - Redefinir Senha - User: " + user.Username + " Status 200 OK!")
	c.JSON(http.StatusOK, gin.H{
		"message": "Se o e-mail informado estiver cadastrado, você receberá um link para redefinir sua senha. Acesse sua caixa de e-mail e click no link recebido",
	})
}

func ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos."})
		return
	}

	db, err := pkg.OpenConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro na conexão com o banco de dados."})
		return
	}
	defer db.Close()

	var user User
	var expiresAt time.Time
	var used bool
	currentTime := time.Now()

	err = db.QueryRow(`SELECT user_id, expires_at, used FROM public.password_reset_tokens WHERE token = $1`, req.Token).Scan(&user.ID, &expiresAt, &used)
	if err == sql.ErrNoRows || used || time.Now().After(expiresAt) {
		fmt.Println("[GIN] " + currentTime.Format("2006/01/02 - 15:04:05") + " | A1 - Erro ao redefinir Senha - User: " + user.ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token inválido ou expirado." + expiresAt.String()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criptografar a senha."})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar a transação."})
		return
	}

	_, err = tx.Exec(`UPDATE public."User" SET password = $1 WHERE id = $2`, string(hashedPassword), user.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar a senha."})
		return
	}

	_, err = tx.Exec(`UPDATE public.password_reset_tokens SET used = true WHERE token = $1`, req.Token)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar o token."})
		return
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao confirmar a transação."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Senha redefinida com sucesso."})
}
