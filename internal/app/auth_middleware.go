package server

import (
	"net/http"
	"strings"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, errMessage := extractToken(c)
		if errMessage != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMessage})
			c.Abort()
			return
		}

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserId)

		c.Next()
	}
}

func extractToken(c *gin.Context) (string, string) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return "", "Invalid authorization header format"
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			return "", "Invalid authorization header format"
		}

		return token, ""
	}

	if c.Request.Method == http.MethodGet {
		token := strings.TrimSpace(c.Query("token"))
		if token == "" {
			return "", "Token missing"
		}

		return token, ""
	}

	return "", "Authorization header missing"
}
