package tokenRepo

import (
	"database/sql"
	storage "messanger/internal/database"
	"messanger/internal/models"
)

type TokenRepository struct {
	db *storage.Storage
}

func InitTokenRepository(db *storage.Storage) *TokenRepository {
	tokenRepository := TokenRepository{
		db: db,
	}

	return &tokenRepository
}

func (repo *TokenRepository) FindToken(refreshToken string) (*models.Token, error) {
	var response models.Token

	err := repo.db.Db.QueryRow(
		"SELECT token_id, user_id, refresh_token FROM tokens WHERE refresh_token = $1",
		refreshToken,
	).Scan(&response.Id, &response.UserId, &response.RefreshToken)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (repo *TokenRepository) SaveToken(refreshToken string, userID int) (*models.Token, error) {
	var response models.Token

	_, err := repo.FindToken(refreshToken)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			err := repo.db.Db.QueryRow(
				"INSERT INTO tokens(refresh_token, user_id) VALUES($1, $2) RETURNING user_id as UserId, refresh_token as RefreshToken",
				refreshToken,
				userID,
			).Scan(&response.UserId, &response.RefreshToken)

			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	err = repo.db.Db.QueryRow(
		"UPDATE tokens SET refresh_token = $1 WHERE user_id = $2 RETURNING token_id, user_id, refresh_token",
		refreshToken,
		userID,
	).Scan(&response.Id, &response.UserId, &response.RefreshToken)

	return &response, nil
}
