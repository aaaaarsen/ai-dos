package handlers

import (
	"strconv"

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


func GetChatsHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {
		query := c.Query("user_id")
		userID, err := strconv.ParseInt(query, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return 
		}

		chats, err := db.GetChatsByUserID(pool, userID)
		if err != nil{
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(200, chats)

	}
}