package db

import (
	"context"
	"errors"

	"github.com/aaaaarsen/ai-dos/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)



func CreateUser (pool *pgxpool.Pool, name *string, email string, passwordHash string) (*models.User, error) {
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, name, created_at`

	user := &models.User{}

	err := pool.QueryRow(context.Background(), query, name, email, passwordHash).Scan(&user.ID, &user.Name, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505"{
			return nil, errors.New("email already exists")
		}

		return nil, err
	}

	user.Email = email
	user.PasswordHash = passwordHash
	return user, nil
}

func GetUserByEmail (pool *pgxpool.Pool, email string) (*models.User, error){
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`

	user := &models.User{}

	err := pool.QueryRow(context.Background(), query, email).Scan(&user.ID,&user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID (pool *pgxpool.Pool, id int64) (*models.User, error) {
	query := `SELECT id, name, email, created_at FROM users where id = $1`

	user := &models.User{}

	err := pool.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteUser(pool *pgxpool.Pool, userID int64) error {
	query := `DELETE FROM users WHERE id = $1`

	tag, err := pool.Exec(context.Background(), query, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
