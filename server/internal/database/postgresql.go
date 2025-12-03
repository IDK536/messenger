package storage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"messanger/internal/config"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db *sql.DB
}

func Connect(cfg *config.Config, logger *slog.Logger) (*Storage, error) {
	connstr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d sslmode=%s dbname=%s",
		cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.SslMode, cfg.DbName,
	)

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}

	logger.Info("Successfully connected to database!")

	sqlCreate := `
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		user_name VARCHAR(128) NOT NULL UNIQUE,
		password VARCHAR(60) NOT NULL,
    	name VARCHAR(128) NOT NULL
	);
	`
	_, err = db.Exec(sqlCreate)
	if err != nil {
		return nil, fmt.Errorf("Error creating users table: %v", err)
	}

	logger.Info("Users table is created!")

	sqlCreate = `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(user_id),
		text TEXT NOT NULL
	);
	`
	_, err = db.Exec(sqlCreate)
	if err != nil {
		return nil, fmt.Errorf("Error creating messages table: %v", err)
	}

	sqlCreate = `
	CREATE TABLE IF NOT EXISTS tokens (
		token_id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(user_id),
		refresh_token TEXT NOT NULL
	);
	`
	_, err = db.Exec(sqlCreate)
	if err != nil {
		return nil, fmt.Errorf("Error creating messages table: %v", err)
	}

	logger.Info("Messages table is created!")

	return &Storage{Db: db}, nil
}
