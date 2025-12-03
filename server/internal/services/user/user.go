package userService

import (
	repo "messanger/internal/database/repositories/user"
	"messanger/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repo.UserRepository
}

type UserServiceInterface interface {
	RegisterUser(username string, name string, password string) (*repo.CreateUserResponse, error)
	Authorization(name string, password string) (*models.User, error)
}

func InitUserService(userRepo *repo.UserRepository) *UserService {
	userService := UserService{
		userRepository: userRepo,
	}

	return &userService
}

func (service *UserService) RegisterUser(username string, name string, password string) (*repo.CreateUserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	regUser, err := service.userRepository.CreateUser(username, name, string(hashedPassword))

	if err != nil {
		return nil, err
	}

	return regUser, err
}

func (service *UserService) FindUserByID(userID int) (*FindUserResponse, error) {
	user, err := service.userRepository.FindUserByID(userID)

	if err != nil {
		return nil, err
	}

	return &FindUserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}, nil
}

func (service *UserService) Authorization(user_name string, password string) (*models.User, error) {
	user, err := service.userRepository.FindUserByUserName(user_name)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) UpdateUserData(userID int, username string, name string) (*models.User, error) {
	user, err := service.userRepository.UpdateUserData(userID, username, name)

	if err != nil {
		return nil, err
	}
	return user, nil
}
