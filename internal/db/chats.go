package db

import (
	"context"
	"errors"

	"github.com/aaaaarsen/ai-dos/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateChat (pool *pgxpool.Pool, userID int64, title *string) (*models.Chat, error) {
	query := `INSERT INTO chats (user_id, title) VALUES ($1, $2) RETURNING id, created_at`

	chat := &models.Chat{UserID: userID, Title: title}

	err := pool.QueryRow(context.Background(), query, userID, title).Scan(&chat.ID, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func GetChatsByUserID (pool *pgxpool.Pool, userID int64) ([]models.Chat, error) {
	query := 	`SELECT 
					id, user_id, title, created_at,
					(SELECT content FROM messages 
					WHERE chat_id = chats.id 
					ORDER BY created_at DESC 
					LIMIT 1) as last_message
				FROM chats 
				WHERE user_id = $1 
				ORDER BY created_at DESC`

	var chats []models.Chat

	rows, err := pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var chat models.Chat
		err = rows.Scan(&chat.ID, &chat.UserID, &chat.Title, &chat.CreatedAt, &chat.LastMessage)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return chats, nil
}

func DeleteChat(pool *pgxpool.Pool, chatID int64, userID int64) error {
	query := `DELETE FROM chats WHERE id = $1 AND user_id = $2`

	tag, err := pool.Exec(context.Background(), query, chatID, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("chat not found")
	}
	return nil
}