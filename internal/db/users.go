package db

import (
	"context"

	"github.com/aaaaarsen/ai-dos/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	
)



func CreateUser (pool *pgxpool.Pool, email string, passwordHash string) (*models.User, error) {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`

	user := &models.User{Email: email, PasswordHash: passwordHash}

	err := pool.QueryRow(context.Background(), query, email, passwordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByEmail (pool *pgxpool.Pool, email string) (*models.User, error){
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`

	user := &models.User{}

	err := pool.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}