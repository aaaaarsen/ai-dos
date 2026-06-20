package db

import (
	"context"
	"slices"

	"github.com/aaaaarsen/ai-dos/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateMessage (pool*pgxpool.Pool, chatID int64, role string, content string) (*models.Message, error){
	query := `INSERT INTO messages (chat_id, role, content) VALUES ($1, $2, $3) RETURNING id, created_at`
	
	message := &models.Message{ChatID: chatID, Role: role, Content: content}

	err := pool.QueryRow(context.Background(), query, chatID, role, content).Scan(&message.ID, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func GetMessagesByChatID (pool *pgxpool.Pool, chatID int64) ([]models.Message, error){
	query := `SELECT id, chat_id, role, content, created_at FROM messages WHERE chat_id = $1 ORDER BY created_at ASC`

	var messages []models.Message

	rows, err := pool.Query(context.Background(), query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	for rows.Next(){
		var message models.Message
		err = rows.Scan(&message.ID, &message.ChatID, &message.Role, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func CountMessagesByChatID(pool *pgxpool.Pool, chatID int64) (int, error){
	query := `SELECT COUNT(*) FROM messages WHERE chat_id = $1`

	var count int
	err := pool.QueryRow(context.Background(), query, chatID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetRecentMessages(pool *pgxpool.Pool, chatID int64, limit int)([]models.Message, error){
	query := `SELECT id, chat_id, role, content, created_at FROM messages WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2`

	var messages []models.Message

	rows, err := pool.Query(context.Background(), query, chatID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var message models.Message
		err = rows.Scan(&message.ID, &message.ChatID, &message.Role, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
		
	}

	slices.Reverse(messages)
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
