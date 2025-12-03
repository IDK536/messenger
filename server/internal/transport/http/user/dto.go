package userController

import tokenServiceDto "messanger/internal/services/token"

type RegisterUserBody struct {
	Username string `json:"username" xml:"username" form:"username"`
	Name     string `json:"name" xml:"name" form:"name"`
	Password string `json:"password" xml:"password" form:"password"`
}

type AuthorizationBody struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
}

type AuthorizationResponse struct {
	Tokens  *tokenServiceDto.GenerateTokensResponse `json:"tokens"`
	Payload *AuthorizationResponsePayload           `json:"payload"`
}

type AuthorizationResponsePayload struct {
	Username string `json:"username"`
	Name     string `json:"name" xml:"name" form:"name"`
}

type UpdateUserDataBody struct {
	Username string `json:"username"`
	Name     string `json:"name" xml:"name" form:"name"`
}
