package handlers

import (

	"github.com/aaaaarsen/ai-dos/internal/auth"
	"github.com/aaaaarsen/ai-dos/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRequest struct {
	Name *string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

// @Summary      Регистрация
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body AuthRequest true "Email и пароль"
// @Success      201  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /auth/register [post]
func RegisterHandler(pool *pgxpool.Pool, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return 
		}
		if req.Email == "" {
			c.JSON(400, gin.H{"error": "email is required"})
			return 
		}
		if len(req.Password) < 6 {
			c.JSON(400, gin.H{"error": "password must be at least 6 characters"})
			return 
		}

		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}

		user, err := db.CreateUser(pool,req.Name, req.Email, hashedPassword)
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

// @Summary      Вход
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body AuthRequest true "Email и пароль"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /auth/login [post]
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

// @Summary      Профиль пользователя
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /users/me [get]
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

// @Summary      Удалить аккаунт
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      204
// @Failure      401  {object}  map[string]interface{}
// @Router       /users/me [delete]
func DeleteMeHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {
		value, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error":"unauthorized"})
			return 
		}
		userID := value.(int64)

		err := db.DeleteUser(pool, userID)
		if err != nil {
			if err.Error() == "user not found" {
				c.JSON(404, gin.H{"error": err.Error()})
				return 
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(204, nil)
	}
}