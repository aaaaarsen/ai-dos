package handlers

import (
	
	"strings"

	"github.com/aaaaarsen/ai-dos/internal/ai"
	"github.com/aaaaarsen/ai-dos/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)


func GetInsightsHandler(pool *pgxpool.Pool, groqKey string, groqModel string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return 
		}
		userID := value.(int64)

		summaries, err := db.GetAllSummariesByUserID(pool, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		if len(summaries) == 0 {
			c.JSON(200, gin.H{"insight": "Пока недостаточно данных для анализа. Продолжайте вести журнал мыслей."})
			return 
		}

		var insights strings.Builder
		for _, summary := range summaries {
			insights.WriteString(summary.Content + "\n")
		}

		combinedText := insights.String()

		generatedInsight, err := ai.GenerateReply(groqKey, groqModel, []ai.ChatMessage{
			{Role: "system", Content: "Ты — аналитик эмоциональных паттернов. На основе следующих наблюдений о человеке составь краткий, тёплый и честный анализ: какие темы повторяются, какие эмоциональные паттерны прослеживаются, что может быть полезно осознать этому человеку. Пиши от второго лица (ты), 3-5 предложений, без диагнозов и медицинских терминов."},
			{Role: "user", Content: combinedText},
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return 
		}
		c.JSON(200, gin.H{"insight": generatedInsight})
	}
}