package userRepo

import (
	storage "messanger/internal/database"
	"messanger/internal/models"
)

type UserRepository struct {
	db *storage.Storage
}

func InitUserRepository(db *storage.Storage) *UserRepository {
	userRepositpry := UserRepository{
		db: db,
	}

	return &userRepositpry
}

func (repo *UserRepository) CreateUser(username string, name string, hasedPassword string) (*CreateUserResponse, error) {
	var result CreateUserResponse

	err := repo.db.Db.QueryRow(
		"INSERT INTO users (user_name, name, password) VALUES ($1, $2, $3) RETURNING user_id as id, user_name, name",
		username,
		name,
		hasedPassword,
	).Scan(&result.ID, &result.Username, &result.Name)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *UserRepository) FindUserByUserName(username string) (*models.User, error) {
	var result models.User

	err := repo.db.Db.QueryRow(
		"SELECT user_id, user_name, name, password FROM users WHERE user_name = $1",
		username,
	).Scan(&result.ID, &result.Username, &result.Name, &result.Password)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *UserRepository) FindUserByID(userId int) (*models.User, error) {
	var result models.User

	err := repo.db.Db.QueryRow(
		"SELECT user_id, user_name, name, password FROM users WHERE user_id = $1",
		userId,
	).Scan(&result.ID, &result.Username, &result.Name, &result.Password)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *UserRepository) UpdateUserData(userID int, username string, name string) (*models.User, error) {
	var result models.User

	err := repo.db.Db.QueryRow(
		"UPDATE users SET user_name = $1, name = $2 WHERE user_id = $3 RETURNING user_id, user_name, name, password",
		username,
		name,
		userID,
	).Scan(&result.ID, &result.Username, &result.Name, &result.Password)

	if err != nil {
		return nil, err
	}

	return &result, err
}
