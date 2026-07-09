package handlers

import (
	"github.com/aaaaarsen/ai-dos/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SaveMoodRequest struct{
	Emoji string `json:"emoji"`
}

func SaveMoodHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {
		var req SaveMoodRequest

		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return 
		}
		if req.Emoji == "" {
			c.JSON(400, gin.H{"error": "emoji is required"})
			return
		}

		value, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error":"unauthorized"})
			return 
		}
		userID := value.(int64)

		mood, err := db.SaveMood(pool, userID, req.Emoji)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(200, mood)
	}
}

func GetTodayMoodHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return 
		}
		userID := value.(int64)

		todayMood, err := db.GetTodayMood(pool, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		if todayMood == nil {
			c.JSON(200, gin.H{"mood": nil})
			return 
		}
		c.JSON(200, todayMood)
		
	}
}