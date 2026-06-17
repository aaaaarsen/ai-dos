package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"fmt"
	"github.com/aaaaarsen/ai-dos/internal/db"
	"github.com/aaaaarsen/ai-dos/internal/handlers"
	"context"
	"github.com/gin-gonic/gin"

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
	router.POST("/chats", handlers.CreateChatHandler(pool))
	router.GET("/chats", handlers.GetChatsHandler(pool))

	router.GET("/health", func(c *gin.Context) {c.JSON(200, gin.H{"status": "ok"})})
	router.Run(":"+serverPort)
}
