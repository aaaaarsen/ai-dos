package handlers

import (
	"strconv"

	"github.com/aaaaarsen/ai-dos/internal/db"
	"github.com/aaaaarsen/ai-dos/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateChatRequest struct {
	Title *string `json:"title"`
} 

func CreateChatHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context){
		var req CreateChatRequest
		value, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return 
		}
		userID := value.(int64)
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return 
		}

		chat, err := db.CreateChat(pool, userID, req.Title)
		if err != nil{
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(201, chat)
	}
}


func GetChatsHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {
		value, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return 
		}
		userID := value.(int64)

		chats, err := db.GetChatsByUserID(pool, userID)
		if err != nil{
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		if chats == nil {
			chats = []models.Chat{}
		}
		c.JSON(200, chats)

	}
}

func DeleteChatHandler(pool *pgxpool.Pool)gin.HandlerFunc{
	return func(c *gin.Context) {
		value1 := c.Param("id")
		chatID, err := strconv.ParseInt(value1, 10, 64)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}

		value2, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return 
		}
		userID := value2.(int64)

		err = db.DeleteChat(pool, chatID, userID)
		if err != nil{
			if err.Error() == "chat not found" {
				c.JSON(404, gin.H{"error": err.Error()})
				return 
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(204, nil)
	}
}

func GetSummariesHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.Param("id")
		chatID, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return 
		}

		summaries, err := db.GetRecentSummaries(pool, chatID, 10)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(200, summaries)
	}
}