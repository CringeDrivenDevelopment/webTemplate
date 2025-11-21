package auth

import (
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/infra"
	"backend/internal/infra/queries"
)

func TestGenerateToken(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}

	service := NewService(cfg)

	userID := "test-user-123"

	token, err := service.GenerateToken(userID)

	require.NoError(t, err)

	assert.NotEmpty(t, token)
}

func TestVerifyToken_Success(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}

	service := NewService(cfg)

	userID := "test-user-123"

	token, err := service.GenerateToken(userID)

	require.NoError(t, err)

	extractedUserID, err := service.VerifyToken("Bearer " + token)

	require.NoError(t, err)

	assert.Equal(t, userID, extractedUserID)
}

func TestVerifyToken_EmptyToken(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}

	service := NewService(cfg)

	_, err := service.VerifyToken("")

	assert.Error(t, err)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}

	service := NewService(cfg)

	_, err := service.VerifyToken("Bearer invalid-token")

	assert.Error(t, err)
}

func TestVerifyPassword_Success(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}

	service := NewService(cfg)

	password := "test-password"

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)

	require.NoError(t, err)

	user := queries.User{
		PasswordHash: hash,
	}

	err = service.VerifyPassword(user, password)

	assert.NoError(t, err)
}

func TestVerifyPassword_InvalidPassword(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}

	service := NewService(cfg)

	password := "test-password"

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)

	require.NoError(t, err)

	user := queries.User{
		PasswordHash: hash,
	}

	err = service.VerifyPassword(user, "wrong-password")

	assert.Error(t, err)
}

func TestNewService(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}

	service := NewService(cfg)

	assert.NotNil(t, service)

	assert.Equal(t, "test-secret", service.secret)

	assert.Equal(t, time.Hour, service.expires)
}
