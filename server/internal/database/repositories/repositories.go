package repositories

import (
	storage "messanger/internal/database"
	tokenRepo "messanger/internal/database/repositories/token"
	userRepo "messanger/internal/database/repositories/user"
)

type Repositories struct {
	UserRepository  *userRepo.UserRepository
	TokenRepository *tokenRepo.TokenRepository
}

func InitRepositories(db *storage.Storage) *Repositories {
	userRepository := userRepo.InitUserRepository(db)
	tokenRepository := tokenRepo.InitTokenRepository(db)

	Repositories := Repositories{
		UserRepository:  userRepository,
		TokenRepository: tokenRepository,
	}

	return &Repositories
}
