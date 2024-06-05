package services

import (
	"auth_service/models"
	"auth_service/repositories"
	"errors"
	"regexp"
)

// AuthService представляет сервис авторизации
type AuthService struct {
	UserRepository *repositories.UserRepository
}

// NewAuthService создает новый экземпляр сервиса авторизации
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{UserRepository: userRepo}
}

// RegisterUser регистрирует нового пользователя
func (service *AuthService) RegisterUser(username, email, password string) error {
	if err := service.validatePassword(password); err != nil {
		return err
	}
	user := &models.User{
		Username:    username,
		Email:       email,
		Password:    password,
		AccessLevel: 1, // Начальный уровень доступа
		RatingLevel: 1, // Начальный рейтинг
	}

	err := service.UserRepository.Save(user)
	if err != nil {
		return err
	}
	return nil
}

// RegisterAdmin регистрирует нового администратора (только для суперадминистратора)
func (service *AuthService) RegisterAdmin(username, email, password string, accessLevel int, superAdmin *models.User) error {
	if superAdmin.AccessLevel != 1 {
		return errors.New("only super administrators can create administrator accounts")
	}
	if err := service.validatePassword(password); err != nil {
		return err
	}
	admin := &models.User{
		Username:    username,
		Email:       email,
		Password:    password,
		AccessLevel: accessLevel,
		RatingLevel: 0, // Администраторы не имеют рейтинга
	}

	err := service.UserRepository.SaveAdmin(admin)
	if err != nil {
		return err
	}
	return nil
}

// AuthenticateUser аутентифицирует пользователя
func (service *AuthService) AuthenticateUser(identifier, password string) (*models.User, error) {
	user, err := service.UserRepository.Authenticate(identifier, password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid user")
	}
	return user, nil
}

// UpdateUserAccessAndRating обновляет уровень доступа и рейтинг пользователя (только для администраторов)
func (service *AuthService) UpdateUserAccessAndRating(userID string, accessLevel, ratingLevel int, admin *models.User) error {
	if admin.AccessLevel > 2 {
		return errors.New("insufficient privileges to update user access level and rating")
	}
	if accessLevel < 1 || accessLevel > 3 {
		return errors.New("invalid access level")
	}
	if ratingLevel < 0 || ratingLevel > 9 {
		return errors.New("invalid rating level")
	}
	return service.UserRepository.UpdateAccessAndRatingLevel(userID, accessLevel, ratingLevel)
}

// validatePassword проверяет валидность пароля
func (service *AuthService) validatePassword(password string) error {
	if len(password) == 0 {
		return errors.New("password cannot be empty")
	}

	// Пример валидации пароля
	var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9]{6,}$`)
	if !passwordRegex.MatchString(password) {
		return errors.New("password must be at least 6 characters long and contain only letters and numbers")
	}
	return nil
}
