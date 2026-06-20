package middleware

import (
	"strings"

	"github.com/aaaaarsen/ai-dos/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware (jwtString string) gin.HandlerFunc{
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "missing authorization header"})
			c.Abort()
			return 
		}

		if !strings.HasPrefix(authHeader, "Bearer "){
			c.JSON(401, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return 
		}

		tokenString := authHeader[7:]

		claims, err := auth.ValidateToken(tokenString, jwtString)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return 
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}