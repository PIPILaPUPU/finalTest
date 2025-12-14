package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret    []byte
	passwordHash string
)

// Init инициализирует модуль аутентификации
func Init() error {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return nil // Аутентификация не требуется
	}

	jwtSecret = []byte(password)

	// Создаем хэш пароля для хранения в токене
	hash := sha256.Sum256([]byte(password))
	passwordHash = hex.EncodeToString(hash[:])

	return nil
}

// Claims структура claims для JWT
type Claims struct {
	PasswordHash string `json:"pwd_hash"`
	jwt.RegisteredClaims
}

// GenerateToken создает JWT-токен
func GenerateToken() (string, error) {
	if jwtSecret == nil {
		return "", errors.New("аутентификация не настроена")
	}

	expirationTime := time.Now().Add(8 * time.Hour)

	claims := &Claims{
		PasswordHash: passwordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken проверяет JWT-токен
func ValidateToken(tokenString string) (bool, error) {
	if jwtSecret == nil {
		return true, nil // Аутентификация не требуется
	}

	if tokenString == "" {
		return false, errors.New("токен не предоставлен")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, errors.New("невалидный токен")
	}

	// Проверяем хэш пароля в claims
	if claims, ok := token.Claims.(*Claims); ok {
		currentHash := sha256.Sum256(jwtSecret)
		currentHashStr := hex.EncodeToString(currentHash[:])

		if claims.PasswordHash != currentHashStr {
			return false, errors.New("пароль был изменен")
		}

		return true, nil
	}

	return false, errors.New("ошибка при разборе claims")
}

// IsAuthRequired проверяет, требуется ли аутентификация
func IsAuthRequired() bool {
	return os.Getenv("TODO_PASSWORD") != ""
}

// GetPasswordHash возвращает хэш пароля
func GetPasswordHash() string {
	return passwordHash
}
