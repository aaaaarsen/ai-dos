package db

import (
	"context"
	"errors"

	"github.com/aaaaarsen/ai-dos/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SaveMood(pool *pgxpool.Pool, userID int64, emoji string) (*models.Mood, error){
	query := `INSERT INTO moods (user_id, emoji, date)
		VALUES ($1, $2, CURRENT_DATE)
		ON CONFLICT (user_id, date)
		DO UPDATE SET emoji = EXCLUDED.emoji
		RETURNING id, user_id, emoji, date, created_at`

	mood := &models.Mood{UserID: userID, Emoji: emoji}

	err := pool.QueryRow(context.Background(), query, userID, emoji).Scan(&mood.ID,&mood.UserID, &mood.Emoji, &mood.Date, &mood.CreatedAt)
	if err != nil {
		return nil, err
	}
	return mood, nil
}

func GetTodayMood(pool *pgxpool.Pool, userID int64)(*models.Mood, error){
	query := `SELECT id, user_id, emoji, date, created_at
			FROM moods
			WHERE user_id = $1 AND date = CURRENT_DATE`
	
	moodToday := &models.Mood{}

	err := pool.QueryRow(context.Background(), query, userID).Scan(&moodToday.ID, &moodToday.UserID, &moodToday.Emoji, &moodToday.Date, &moodToday.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return nil, nil
		}
		return nil, err
	}

	return moodToday, nil
}