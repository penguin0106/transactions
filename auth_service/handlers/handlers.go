package handlers

import (
	"auth_service/services"
	"encoding/json"
	"net/http"
)

// AuthHandler представляет хендлер для аутентификации
type AuthHandler struct {
	AuthService *services.AuthService
	JWTService  *services.JWTService
}

// NewAuthHandler создает новый экземпляр хендлера аутентификации
func NewAuthHandler(authService *services.AuthService, jwtService *services.JWTService) *AuthHandler {
	return &AuthHandler{AuthService: authService, JWTService: jwtService}
}

// RegisterHandler обрабатывает запросы на регистрацию пользователей
func (handler *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = handler.AuthService.RegisterUser(user.Username, user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
}

// RegisterAdminHandler обрабатывает запросы на регистрацию администраторов (только для суперадминистратора)
func (handler *AuthHandler) RegisterAdminHandler(w http.ResponseWriter, r *http.Request) {
	var admin struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		AccessLevel int    `json:"access_level"`
	}

	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Получаем токен из заголовка и проверяем права
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	claims, err := handler.JWTService.VerifyToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	superAdmin, err := handler.AuthService.UserRepository.FindByEmail(claims.Email)
	if err != nil || superAdmin == nil {
		http.Error(w, "Super administrator not found", http.StatusUnauthorized)
		return
	}

	if superAdmin.AccessLevel != 1 {
		http.Error(w, "Insufficient privileges to register an admin", http.StatusUnauthorized)
		return
	}

	err = handler.AuthService.RegisterAdmin(admin.Username, admin.Email, admin.Password, admin.AccessLevel, superAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "admin registered successfully"})
}

// AuthenticateHandler обрабатывает запросы на аутентификацию пользователей
func (handler *AuthHandler) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := handler.AuthService.AuthenticateUser(credentials.Identifier, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := handler.JWTService.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// UpdateUserAccessAndRatingHandler обрабатывает запросы на обновление уровня доступа и рейтинга пользователей
func (handler *AuthHandler) UpdateUserAccessAndRatingHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID      string `json:"user_id"`
		AccessLevel int    `json:"access_level"`
		RatingLevel int    `json:"rating_level"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Получаем токен из заголовка и проверяем права
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	claims, err := handler.JWTService.VerifyToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	admin, err := handler.AuthService.UserRepository.FindByEmail(claims.Email)
	if err != nil || admin == nil {
		http.Error(w, "Administrator not found", http.StatusUnauthorized)
		return
	}

	if admin.AccessLevel > 2 {
		http.Error(w, "Insufficient privileges to update user access level and rating", http.StatusUnauthorized)
		return
	}

	err = handler.AuthService.UpdateUserAccessAndRating(request.UserID, request.AccessLevel, request.RatingLevel, admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "user access level and rating updated successfully"})
}

// VerifyTokenHandler обрабатывает запросы на верификацию токена
func (handler *AuthHandler) VerifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	userID, username, email, err := handler.JWTService.VerifyToken(requestBody.Token)
	if err != nil {
		http.Error(w, "Failed to verify token", http.StatusUnauthorized)
		return
	}

	response := struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		UserID:   userID,
		Username: username,
		Email:    email,
	}

	json.NewEncoder(w).Encode(response)
}
