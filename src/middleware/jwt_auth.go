package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ausente"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato inválido do token"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Extrair account_id do token e armazenar no contexto
		accountID, ok := claims["account_id"].(string)
		if !ok || accountID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "account_id ausente ou inválido no token"})
			c.Abort()
			return
		}

		// (Opcional) Extrair outros dados úteis
		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)

		// Armazenar no contexto
		c.Set("account_id", accountID)
		c.Set("user_id", userID)
		c.Set("email", email)

		c.Next()
	}
}
