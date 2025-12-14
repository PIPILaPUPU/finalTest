package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/PIPILaPUPU/finalTest/pkg/auth"
)

// SignInRequest структура запроса аутентификации
type SignInRequest struct {
	Password string `json:"password"`
}

// SignInResponse структура ответа аутентификации
type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// signInHandler обрабатывает запрос аутентификации
func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, SignInResponse{Error: "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	// Декодируем запрос
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, SignInResponse{Error: "Ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	// Проверяем пароль
	expectedPassword := os.Getenv("TODO_PASSWORD")

	// Если пароль не настроен, пропускаем аутентификацию
	if expectedPassword == "" {
		writeJSON(w, SignInResponse{Error: "Аутентификация не настроена"}, http.StatusBadRequest)
		return
	}

	// Сравниваем пароли
	if req.Password != expectedPassword {
		writeJSON(w, SignInResponse{Error: "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	// Генерируем JWT-токен
	token, err := auth.GenerateToken()
	if err != nil {
		writeJSON(w, SignInResponse{Error: "Ошибка при генерации токена"}, http.StatusInternalServerError)
		return
	}

	// Возвращаем токен
	writeJSON(w, SignInResponse{Token: token}, http.StatusOK)
}
