package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Récupère le header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("auth_request", zap.String("error", "token manquant"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token manquant"})
			return
		}

		// 2. Extrait le token (format "Bearer <token>")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("auth_request", zap.String("error", "format de token invalide"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "format de token invalide"})
			return
		}

		// 3. Parse et vérifie le token
		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Warn("auth_request", zap.String("error", "méthode de signature invalide"))
				return nil, errors.New("méthode de signature invalide")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			logger.Warn("auth_request", zap.String("error", "token invalide ou expiré"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalide ou expiré"})
			return
		}

		// 4. Injecte le user_id dans le contexte pour les handlers
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("user_role", claims["role"].(string))

		c.Next()
	}
}
