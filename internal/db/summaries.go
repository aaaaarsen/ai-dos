package db

import (
	"context"

	"github.com/aaaaarsen/ai-dos/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateSummary(pool *pgxpool.Pool, chatID int64, content string)(*models.Summary, error){
	query := `INSERT INTO summaries (chat_id, content) VALUES ($1, $2) RETURNING id, created_at`
	
	summary := &models.Summary{ChatID: chatID, Content: content}

	err := pool.QueryRow(context.Background(), query, chatID, content).Scan(&summary.ID, &summary.CreatedAt)
	if err != nil {
		return nil, err 
	}
	return summary, nil
}

func GetRecentSummaries(pool *pgxpool.Pool, chatID int64, limit int) ([]models.Summary, error){
	query := `SELECT id, chat_id, content, created_at FROM summaries WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2`

	var summaries []models.Summary

	rows, err := pool.Query(context.Background(), query, chatID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var summary models.Summary
		err = rows.Scan(&summary.ID, &summary.ChatID, &summary.Content, &summary.CreatedAt)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return summaries, nil
}

func GetAllSummariesByUserID(pool *pgxpool.Pool, userID int64) ([]models.Summary, error){
	query:=`SELECT s.id, s.chat_id, s.content, s.created_at 
			FROM summaries s 
			JOIN chats c ON s.chat_id = c.id 
			WHERE c.user_id = $1 
			ORDER BY s.created_at DESC 
			LIMIT 20`

	var summaries []models.Summary

	rows, err := pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var summary models.Summary
		err = rows.Scan(&summary.ID, &summary.ChatID, &summary.Content, &summary.CreatedAt)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return summaries, nil
}