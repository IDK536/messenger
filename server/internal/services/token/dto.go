package tokenService

import "github.com/golang-jwt/jwt/v5"

type GenerateTokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenPayload struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Claims struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
