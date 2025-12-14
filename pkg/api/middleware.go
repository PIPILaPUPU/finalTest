package api

import (
	"net/http"

	"github.com/PIPILaPUPU/finalTest/pkg/auth"
)

// AuthMiddleware middleware для проверки аутентификации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Если аутентификация не требуется, пропускаем
		if !auth.IsAuthRequired() {
			next(w, r)
			return
		}

		// Получаем токен из куки
		var token string
		cookie, err := r.Cookie("token")
		if err == nil {
			token = cookie.Value
		}

		// Проверяем токен
		valid, err := auth.ValidateToken(token)
		if !valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Токен валиден, продолжаем выполнение
		next(w, r)
	}
}
