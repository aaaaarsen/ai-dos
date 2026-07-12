package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aaaaarsen/ai-dos/internal/db"
	"github.com/aaaaarsen/ai-dos/internal/handlers"
	"github.com/aaaaarsen/ai-dos/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	_ = godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSslmode := os.Getenv("DB_SSLMODE")
	serverPort := os.Getenv("SERVER_PORT")

	jwtSecret := os.Getenv("JWT_SECRET")

	groqKey := os.Getenv("GROQ_API_KEY")
	groqModel := os.Getenv("GROQ_MODEL")

	

	
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSslmode) 

	err := db.RunMigrations(dsn)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations applied succesfully")

	pool, err := db.NewPool(dsn)
	if err != nil {
		log.Fatalf("Connect failed: %v", err)
	}
	defer pool.Close()
	err = pool.Ping(context.Background())
	
	if err != nil {
		log.Fatalf("Connect failed: %v", err)
	}
	log.Println("Database connected successfully")


	router := gin.Default()

	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api.POST("/auth/register", handlers.RegisterHandler(pool, jwtSecret))
	api.POST("/auth/login", handlers.LoginHandler(pool, jwtSecret))

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	protected.GET("/users/me", handlers.GetMeHandler(pool))
	protected.DELETE("/users/me", handlers.DeleteMeHandler(pool))
	protected.GET("/users/me/insights", handlers.GetInsightsHandler(pool, groqKey, groqModel))
	protected.GET("/users/me/stats", handlers.GetStatsHandler(pool))
	protected.POST("/mood", handlers.SaveMoodHandler(pool))
	protected.GET("/mood/today", handlers.GetTodayMoodHandler(pool))
	protected.POST("/chats", handlers.CreateChatHandler(pool))
	protected.GET("/chats", handlers.GetChatsHandler(pool))
	protected.DELETE("/chats/:id", handlers.DeleteChatHandler(pool))
	protected.POST("/chats/:id/messages", handlers.CreateMessageHandler(pool, groqKey, groqModel))
	protected.GET("/chats/:id/messages", handlers.GetMessagesHandler(pool))
	protected.GET("/chats/:id/summaries", handlers.GetSummariesHandler(pool))

	router.Run(":" + serverPort)
}
