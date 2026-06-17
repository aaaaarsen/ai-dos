package db

import (
	"context"

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