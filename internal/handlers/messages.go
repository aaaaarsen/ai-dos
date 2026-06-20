package handlers

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aaaaarsen/ai-dos/internal/ai"
	"github.com/aaaaarsen/ai-dos/internal/db"
)

type CreateMessageRequest struct {
	Content string `json:"content"`
}

func CreateMessageHandler(pool *pgxpool.Pool, groqKey string, groqModel string) gin.HandlerFunc {
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

		recentSummaries, err := db.GetRecentSummaries(pool, chatID, 3)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}

		recentMessages, err := db.GetRecentMessages(pool, chatID, 10)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}

		var chatMessage []ai.ChatMessage
		chatMessage = append(chatMessage, ai.ChatMessage{Role: "system", Content: ai.SystemPrompt})

		for _, summary := range recentSummaries {
			chatMessage = append(chatMessage, ai.ChatMessage{Role: "system", Content: "Ранее: "+summary.Content})
		}

		for _, msg := range recentMessages {
			chatMessage = append(chatMessage, ai.ChatMessage{Role: msg.Role, Content: msg.Content})
		}
		
		reply, err := ai.GenerateReply(groqKey, groqModel, chatMessage)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}


		assistantMessage, err := db.CreateMessage(pool, chatID, "assistant",reply)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}

		count, err := db.CountMessagesByChatID(pool, chatID)
		
		if err != nil {
			log.Println("count failed:", err)
		} else if err == nil && count%10 == 0 {
			recentTen, err := db.GetRecentMessages(pool, chatID, 10)
			if err != nil {
				log.Println("get recent ten failed:", err)
			} else {
				var sb strings.Builder
				for _, m := range recentTen {
					sb.WriteString(m.Role + ": " + m.Content + "\n")
				}

				generatedSummary, err := ai.GenerateReply(groqKey, groqModel, []ai.ChatMessage{
					{Role: "system", Content: ai.SummaryPrompt},
					{Role: "user", Content: sb.String()},
				})
				if err != nil {
					log.Println("summary generation failed:", err)
				} else {
					_, err = db.CreateSummary(pool, chatID, generatedSummary)
					if err != nil {
						log.Println("save summary failed:", err)
					}
				}
			}
			
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