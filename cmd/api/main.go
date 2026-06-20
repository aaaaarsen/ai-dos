package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aaaaarsen/ai-dos/internal/ai"
	"github.com/aaaaarsen/ai-dos/internal/db"
	"github.com/aaaaarsen/ai-dos/internal/handlers"
	"github.com/aaaaarsen/ai-dos/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	protected.POST("/chats", handlers.CreateChatHandler(pool))
	protected.GET("/chats", handlers.GetChatsHandler(pool))
	protected.POST("/chats/:id/messages", handlers.CreateMessageHandler(pool))
	protected.GET("/chats/:id/messages", handlers.GetMessagesHandler(pool))
	router.POST("/auth/register", handlers.RegisterHandler(pool, jwtSecret))
	router.POST("/auth/login", handlers.LoginHandler(pool, jwtSecret))

	router.GET("/health", func(c *gin.Context) {c.JSON(200, gin.H{"status": "ok"})})
	router.Run(":"+serverPort)
}
