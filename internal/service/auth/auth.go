package auth

import (
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"

	"backend/internal/model"
	"backend/pkg/utils"
)

// VerifyToken - проверить токен на подлинность
func (s *Service) VerifyToken(authHeader string) (string, error) {
	tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenStr == "" {
		return "", utils.ErrInvalidToken
	}

	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid || token.Method != jwt.SigningMethodHS256 {
		return "", utils.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", utils.ErrInvalidToken
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", utils.ErrInvalidToken
	}

	return userID, nil
}

func (s *Service) VerifyPassword(user model.User, password string) error {
	valid, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return err
	}

	if !valid {
		return utils.ErrInvalidPassword
	}

	return nil
}

// GenerateToken - создать новый JWT токен
func (s *Service) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(s.expires).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secret))
}
