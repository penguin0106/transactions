package repositories

import (
	"auth_service/models"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

// UserRepository представляет репозиторий пользователей
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// generateUniqueID генерирует уникальный 25-символьный ID пользователя
func generateUniqueID(userType, accessLevel, ratingLevel int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 22) // Оставшиеся 22 символа (первые 3 - для типа, прав и рейтинга)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return fmt.Sprintf("%d%d%d%s", userType, accessLevel, ratingLevel, string(b))
}

// Save сохраняет пользователя в базе данных
func (repo *UserRepository) Save(user *models.User) error {
	existingUser, err := repo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}
	existingUserName, err := repo.FindByUserName(user.Username)
	if err != nil {
		return err
	}
	if existingUserName != nil {
		return errors.New("username already exists")
	}

	userType := 2 // 1 - администратор, 2 - пользователь

	user.ID = generateUniqueID(userType, user.AccessLevel, user.RatingLevel)

	_, err = repo.DB.Exec("INSERT INTO users (id, email, username, password, access_level, rating_level) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID, user.Email, user.Username, user.Password, user.AccessLevel, user.RatingLevel)
	if err != nil {
		return err
	}
	return nil
}

// SaveAdmin сохраняет администратора в базе данных
func (repo *UserRepository) SaveAdmin(admin *models.User) error {
	existingUser, err := repo.FindByEmail(admin.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}
	existingUserName, err := repo.FindByUserName(admin.Username)
	if err != nil {
		return err
	}
	if existingUserName != nil {
		return errors.New("username already exists")
	}

	userType := 1 // 1 - администратор, 2 - пользователь

	admin.ID = generateUniqueID(userType, admin.AccessLevel, 0) // Администраторы не имеют рейтинга

	_, err = repo.DB.Exec("INSERT INTO users (id, email, username, password, access_level, rating_level) VALUES ($1, $2, $3, $4, $5, $6)",
		admin.ID, admin.Email, admin.Username, admin.Password, admin.AccessLevel, admin.RatingLevel)
	if err != nil {
		return err
	}
	return nil
}

// FindByEmail ищет пользователя по его email
func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := repo.DB.QueryRow("SELECT id, username, email, password, access_level, rating_level FROM users WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.AccessLevel, &user.RatingLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь с таким email не найден
		}
		return nil, err
	}

	return &user, nil
}

// FindByUserName ищет пользователя по его username
func (repo *UserRepository) FindByUserName(username string) (*models.User, error) {
	var user models.User

	err := repo.DB.QueryRow("SELECT id, username, email, password, access_level, rating_level FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.AccessLevel, &user.RatingLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) Authenticate(identifier, password string) (*models.User, error) {
	user, err := repo.FindByEmail(identifier)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = repo.FindByUserName(identifier)
		if err != nil {
			return nil, err
		}
	}

	if user != nil && user.Password == password {
		return user, nil
	}

	return nil, errors.New("invalid username or password")
}

// UpdateAccessAndRatingLevel обновляет уровень доступа и рейтинг пользователя
func (repo *UserRepository) UpdateAccessAndRatingLevel(userID string, accessLevel, ratingLevel int) error {
	_, err := repo.DB.Exec("UPDATE users SET access_level = $1, rating_level = $2 WHERE id = $3", accessLevel, ratingLevel, userID)
	if err != nil {
		return err
	}
	return nil
}
