package handlers

import (
	"github.com/aaaaarsen/ai-dos/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateChatRequest struct {
	UserID int64 `json:"user_id"`
	Title *string `json:"title"`
} 

func CreateChatHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context){
		var req CreateChatRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return 
		}

		chat, err := db.CreateChat(pool, req.UserID, req.Title)
		if err != nil{
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(201, chat)
	}
}

