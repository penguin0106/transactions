package repositories

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// JWTRepository представляет репозиторий для работы с JWT токенами
type JWTRepository struct {
	DB        *sql.DB
	secretKey []byte
}

// NewJWTRepository создает новый экземпляр репозитория JWT
func NewJWTRepository(db *sql.DB, secretKey []byte) *JWTRepository {
	return &JWTRepository{DB: db, secretKey: secretKey}
}

// GenerateToken генерирует новый JWT токен с данными пользователя
func (repo *JWTRepository) GenerateToken(userID string, username, email string) (string, error) {
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	claims := map[string]interface{}{
		"userID":   userID,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Токен действителен в течение 24 часов
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerBase64 := base64.RawURLEncoding.EncodeToString(headerBytes)
	claimsBase64 := base64.RawURLEncoding.EncodeToString(claimsBytes)

	unsignedToken := fmt.Sprintf("%s.%s", headerBase64, claimsBase64)

	signature := repo.signToken(unsignedToken)

	signedToken := fmt.Sprintf("%s.%s", unsignedToken, signature)

	return signedToken, nil
}

// VerifyToken проверяет действительность JWT токена и возвращает данные пользователя
func (repo *JWTRepository) VerifyToken(tokenString string) (string, string, string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("неверный формат токена")
	}

	unsignedToken := fmt.Sprintf("%s.%s", parts[0], parts[1])
	signature := parts[2]

	expectedSignature := repo.signToken(unsignedToken)
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", "", "", fmt.Errorf("недействительная подпись токена")
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", "", err
	}

	var claims map[string]interface{}
	err = json.Unmarshal(claimsBytes, &claims)
	if err != nil {
		return "", "", "", err
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", "", "", fmt.Errorf("недопустимое значение поля exp")
	}

	if int64(exp) < time.Now().Unix() {
		return "", "", "", fmt.Errorf("токен истек")
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("недопустимое значение поля userID")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("недопустимое значение поля username")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("недопустимое значение поля email")
	}

	return userID, username, email, nil
}

// signToken подписывает токен с использованием HMAC-SHA256
func (repo *JWTRepository) signToken(token string) []byte {
	hash := hmac.New(sha256.New, repo.secretKey)
	hash.Write([]byte(token))
	return hash.Sum(nil)
}
