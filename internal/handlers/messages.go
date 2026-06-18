package handlers

import (

	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/aaaaarsen/ai-dos/internal/db"
)

type CreateMessageRequest struct {
	Content string `json:"content"`
}

func CreateMessageHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateMessageRequest

		query := c.Param("id")
		chatID, err := strconv.ParseInt(query, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return 
		}

		err = c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return 
		}

		_, err = db.CreateMessage(pool, chatID, "user", req.Content)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}

		assistantMessage, err := db.CreateMessage(pool, chatID, "assistant", "Это тестовый ответ ассистента")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(201, assistantMessage)

		
	}
}

func GetMessagesHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Param("id")
		chatID, err := strconv.ParseInt(query, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return 
		}

		messages, err := db.GetMessagesByChatID(pool, chatID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		
		c.JSON(200, messages)
	}
}