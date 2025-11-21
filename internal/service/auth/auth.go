package auth

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"

	"backend/internal/infra/queries"
	"backend/pkg/utils"
)

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", utils.ErrInvalidUser
		}
		return "", err
	}

	err = s.VerifyPassword(user, password)
	if err != nil {
		return "", err
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

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

func (s *Service) VerifyPassword(user queries.User, password string) error {
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
