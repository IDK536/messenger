package tokenService

import (
	"encoding/json"
	"errors"
	"fmt"
	"messanger/internal/config"
	tokenRepo "messanger/internal/database/repositories/token"
	"messanger/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	tokenRepository *tokenRepo.TokenRepository
	cfg             *Config
}

type Config struct {
	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string
}

func InitTokenService(tokenRepo *tokenRepo.TokenRepository, cfg *config.Config) *TokenService {
	serviceCfg := Config{
		JWT_ACCESS_SECRET:  cfg.JWT_ACCESS_SECRET,
		JWT_REFRESH_SECRET: cfg.JWT_REFRESH_SECRET,
	}

	tokenService := TokenService{
		tokenRepository: tokenRepo,
		cfg:             &serviceCfg,
	}

	return &tokenService
}

func (service *TokenService) SaveToken(refreshToken string, userID int) (*models.Token, error) {
	response, err := service.tokenRepository.SaveToken(refreshToken, userID)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *TokenService) ValidateAccessToken(accessToken string) (*TokenPayload, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(service.cfg.JWT_ACCESS_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	payload := TokenPayload{
		Id:       claims.Id,
		Name:     claims.Name,
		Username: claims.Username,
	}

	return &payload, nil
}

func (service *TokenService) ValidateRefreshToken(refreshToken string) (*TokenPayload, error) {
	var cliams jwt.MapClaims

	_, err := jwt.ParseWithClaims(refreshToken, cliams, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(service.cfg.JWT_REFRESH_SECRET), nil
	})
	if err != nil {
		return nil, err
	}

	var payload TokenPayload

	subject, err := cliams.GetSubject()

	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(subject), &payload)

	return &payload, nil
}

func (service *TokenService) GenerateToken(payload *TokenPayload) (*GenerateTokensResponse, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":      time.Now().Add(time.Minute * 30).Unix(),
		"id":       payload.Id,
		"name":     payload.Name,
		"username": payload.Username,
	}).SignedString([]byte(service.cfg.JWT_ACCESS_SECRET))

	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"id":       payload.Id,
		"name":     payload.Name,
		"username": payload.Username,
	}).SignedString([]byte(service.cfg.JWT_REFRESH_SECRET))

	if err != nil {
		return nil, err
	}

	return &GenerateTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
