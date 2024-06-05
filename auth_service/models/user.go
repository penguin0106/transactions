package models

// User представляет модель пользователя
type User struct {
	ID          string `json:"id"` // 25-символьный уникальный ID
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessLevel int    `json:"access_level"` // Уровень доступа (1, 2, 3)
	RatingLevel int    `json:"rating_level"` // Оценка пользователя (0-9)
}
