package services

import (
	"messanger/internal/config"
	"messanger/internal/database/repositories"
	tokenService "messanger/internal/services/token"
	userService "messanger/internal/services/user"
)

type Services struct {
	UserService  *userService.UserService
	TokenService *tokenService.TokenService
}

func InitServices(repositories *repositories.Repositories, cfg *config.Config) *Services {
	userService := userService.InitUserService(repositories.UserRepository)
	tokenService := tokenService.InitTokenService(repositories.TokenRepository, cfg)

	Services := Services{
		UserService:  userService,
		TokenService: tokenService,
	}

	return &Services
}
