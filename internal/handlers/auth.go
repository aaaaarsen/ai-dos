package handlers

import (
	"github.com/aaaaarsen/ai-dos/internal/auth"
	"github.com/aaaaarsen/ai-dos/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(pool *pgxpool.Pool, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return 
		}

		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}

		user, err := db.CreateUser(pool, req.Email, hashedPassword)
		if err != nil {
			if err.Error() == "email already exists"  {
				c.JSON(409, gin.H{"error": err.Error()})
				return 
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}


		tokenString, err := auth.GenerateToken(user.ID, jwtSecret)
		if err != nil{
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(201, gin.H{"token": tokenString, "user": user})
	}
}

func LoginHandler(pool *pgxpool.Pool, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(400, gin.H{"error":err.Error()})
			return 
		}

		user, err := db.GetUserByEmail(pool, req.Email)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return 
		}

		err = auth.CheckPasswordHash(req.Password, user.PasswordHash)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return 
		}

		tokenString, err := auth.GenerateToken(user.ID, jwtSecret)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}

		c.JSON(200, gin.H{"token": tokenString, "user": user})

	}
}

func GetMeHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func (c *gin.Context){
		value, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return 
		}
		userID := value.(int64)

		me, err := db.GetUserByID(pool, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(200, me)
	}
}