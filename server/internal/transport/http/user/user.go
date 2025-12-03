package userController

import (
	"database/sql"
	"messanger/internal/config"
	"messanger/internal/middleware"
	tokenService "messanger/internal/services/token"
	userService "messanger/internal/services/user"

	"time"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userService  *userService.UserService
	tokenService *tokenService.TokenService
	router       fiber.Router
}

type UserControllerInterface interface {
	RegisterUser(c fiber.Ctx) error
	Authorization(c fiber.Ctx) error
	UpdateUserData(c fiber.Ctx) error
	Refresh(c fiber.Ctx) error
}

func InitUserController(userServise *userService.UserService, tokenService *tokenService.TokenService, app *fiber.App) (controller *UserController) {
	userRouter := app.Group("/user")

	UserController := UserController{
		userService:  userServise,
		tokenService: tokenService,
		router:       userRouter,
	}

	return &UserController
}

func (controller *UserController) Start(cfg *config.Config) {
	withoutAuth := controller.router.Group("/")

	withoutAuth.Post("/register", controller.RegisterUser)
	withoutAuth.Post("/auth", controller.Authorization)

	withAuth := controller.router.Group("/", middleware.NewAuth(cfg, controller.tokenService))

	withAuth.Patch("/change/info", controller.UpdateUserData)
}

func (controller *UserController) RegisterUser(c fiber.Ctx) error {
	body := new(RegisterUserBody)

	if err := c.Bind().Body(body); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	user, err := controller.userService.RegisterUser(body.Username, body.Name, body.Password)

	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (controller *UserController) Refresh(c fiber.Ctx) error {
	refreshToken := c.Cookies("refreshToken")

	if refreshToken == "" {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "You are not authorized")
	}

	payload, err := controller.tokenService.ValidateRefreshToken(refreshToken)

	if err != nil {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "You are not authorized(token is not valid)")
	}

	user, err := controller.userService.FindUserByID(payload.Id)

	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	tokens, err := controller.tokenService.GenerateToken(&tokenService.TokenPayload{
		Id:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	})

	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	_, err = controller.tokenService.SaveToken(tokens.RefreshToken, user.ID)

	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	c.Cookie(&fiber.Cookie{
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
	})

	return c.Status(fiber.StatusOK).JSON(tokens)
}

func (controller *UserController) Authorization(c fiber.Ctx) error {
	body := new(AuthorizationBody)

	if err := c.Bind().Body(body); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	user, err := controller.userService.Authorization(body.Username, body.Password)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return fiber.NewError(fiber.ErrBadRequest.Code, "User with this username does not exist!")
		case bcrypt.ErrMismatchedHashAndPassword:
			return fiber.NewError(fiber.ErrBadRequest.Code, "The user password is wrong!")
		default:
			return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
		}
	}

	tokens, err := controller.tokenService.GenerateToken(&tokenService.TokenPayload{
		Id:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	})

	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	_, err = controller.tokenService.SaveToken(tokens.RefreshToken, user.ID)

	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	c.Cookie(&fiber.Cookie{
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
	})

	return c.Status(fiber.StatusOK).JSON(AuthorizationResponse{
		Tokens: tokens,
		Payload: &AuthorizationResponsePayload{
			Name:     user.Name,
			Username: user.Username,
		},
	})
}

func (controller *UserController) UpdateUserData(c fiber.Ctx) error {
	body := new(UpdateUserDataBody)

	if err := c.Bind().Body(body); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	userId, ok := c.Locals("user_id").(int)

	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user_id not found in context")
	}

	user, err := controller.userService.UpdateUserData(userId, body.Name, body.Username)

	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(*user)
}
