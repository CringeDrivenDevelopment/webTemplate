package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/infra"
	"backend/internal/infra/queries"
	repositoryMocks "backend/internal/repository/mocks"
	"backend/pkg/utils"
)

func TestLogin(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}
	ctx := context.Background()

	// Create a test user with hashed password
	password := "SecurePassword123"
	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	require.NoError(t, err)

	userID := ulid.Make().String()
	testUser := queries.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
	}

	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func(*repositoryMocks.MockUserRepository)
		expectedError error
		checkToken    bool
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: password,
			mockSetup: func(mockRepo *repositoryMocks.MockUserRepository) {
				mockRepo.On("GetUserByEmail", ctx, "test@example.com").
					Return(testUser, nil).Once()
			},
			expectedError: nil,
			checkToken:    true,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: password,
			mockSetup: func(mockRepo *repositoryMocks.MockUserRepository) {
				mockRepo.On("GetUserByEmail", ctx, "nonexistent@example.com").
					Return(queries.User{}, pgx.ErrNoRows).Once()
			},
			expectedError: utils.ErrInvalidUser,
			checkToken:    false,
		},
		{
			name:     "database error",
			email:    "test@example.com",
			password: password,
			mockSetup: func(mockRepo *repositoryMocks.MockUserRepository) {
				mockRepo.On("GetUserByEmail", ctx, "test@example.com").
					Return(queries.User{}, errors.New("database connection error")).Once()
			},
			expectedError: errors.New("database connection error"),
			checkToken:    false,
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "WrongPassword123",
			mockSetup: func(mockRepo *repositoryMocks.MockUserRepository) {
				mockRepo.On("GetUserByEmail", ctx, "test@example.com").
					Return(testUser, nil).Once()
			},
			expectedError: utils.ErrInvalidPassword,
			checkToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repositoryMocks.NewMockUserRepository(t)
			tt.mockSetup(mockRepo)
			service := NewService(cfg, mockRepo)

			token, err := service.Login(ctx, tt.email, tt.password)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, token)
			} else {
				require.NoError(t, err)
				if tt.checkToken {
					assert.NotEmpty(t, token)
					// Verify the token is valid
					extractedUserID, verifyErr := service.VerifyToken("Bearer " + token)
					require.NoError(t, verifyErr)
					assert.Equal(t, userID, extractedUserID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}
	mockRepo := repositoryMocks.NewMockUserRepository(t)
	service := NewService(cfg, mockRepo)
	userID := "test-user-123"

	token, err := service.GenerateToken(userID)

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token contains correct user ID
	extractedUserID, verifyErr := service.VerifyToken("Bearer " + token)
	require.NoError(t, verifyErr)
	assert.Equal(t, userID, extractedUserID)
}

func TestVerifyToken(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}
	mockRepo := repositoryMocks.NewMockUserRepository(t)
	service := NewService(cfg, mockRepo)
	userID := "test-user-123"

	tests := []struct {
		name          string
		setupToken    func() string
		expectedError error
		expectedID    string
	}{
		{
			name: "valid token",
			setupToken: func() string {
				token, _ := service.GenerateToken(userID)
				return "Bearer " + token
			},
			expectedError: nil,
			expectedID:    userID,
		},
		{
			name: "empty token",
			setupToken: func() string {
				return ""
			},
			expectedError: utils.ErrInvalidToken,
			expectedID:    "",
		},
		{
			name: "token without Bearer prefix",
			setupToken: func() string {
				token, _ := service.GenerateToken(userID)
				return token
			},
			expectedError: nil,
			expectedID:    userID,
		},
		{
			name: "invalid token format",
			setupToken: func() string {
				return "Bearer invalid-token-string"
			},
			expectedError: errors.New("token is malformed"),
			expectedID:    "",
		},
		{
			name: "expired token (different secret)",
			setupToken: func() string {
				// Create a token with a different secret
				differentCfg := &infra.Config{JwtSecret: "different-secret"}
				differentRepo := repositoryMocks.NewMockUserRepository(t)
				differentService := NewService(differentCfg, differentRepo)
				token, _ := differentService.GenerateToken(userID)
				return "Bearer " + token
			},
			expectedError: errors.New("token signature is invalid"),
			expectedID:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenStr := tt.setupToken()
			extractedID, err := service.VerifyToken(tokenStr)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Empty(t, extractedID)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, extractedID)
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}
	mockRepo := repositoryMocks.NewMockUserRepository(t)
	service := NewService(cfg, mockRepo)

	password := "SecurePassword123"
	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	require.NoError(t, err)

	tests := []struct {
		name          string
		user          queries.User
		password      string
		expectedError error
	}{
		{
			name: "valid password",
			user: queries.User{
				ID:           "user-123",
				Email:        "test@example.com",
				PasswordHash: passwordHash,
			},
			password:      password,
			expectedError: nil,
		},
		{
			name: "invalid password",
			user: queries.User{
				ID:           "user-123",
				Email:        "test@example.com",
				PasswordHash: passwordHash,
			},
			password:      "WrongPassword123",
			expectedError: utils.ErrInvalidPassword,
		},
		{
			name: "empty password",
			user: queries.User{
				ID:           "user-123",
				Email:        "test@example.com",
				PasswordHash: passwordHash,
			},
			password:      "",
			expectedError: utils.ErrInvalidPassword,
		},
		{
			name: "malformed hash",
			user: queries.User{
				ID:           "user-123",
				Email:        "test@example.com",
				PasswordHash: "invalid-hash",
			},
			password:      password,
			expectedError: errors.New("argon2id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.VerifyPassword(tt.user, tt.password)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}
	mockRepo := repositoryMocks.NewMockUserRepository(t)
	service := NewService(cfg, mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, "test-secret", service.secret)
	assert.Equal(t, time.Hour, service.expires)
	assert.NotNil(t, service.repository)
}

func TestTokenExpiration(t *testing.T) {
	cfg := &infra.Config{JwtSecret: "test-secret"}
	mockRepo := repositoryMocks.NewMockUserRepository(t)
	service := NewService(cfg, mockRepo)
	userID := "test-user-123"

	// Generate a token
	token, err := service.GenerateToken(userID)
	require.NoError(t, err)

	// Verify the token immediately (should be valid)
	extractedID, err := service.VerifyToken("Bearer " + token)
	require.NoError(t, err)
	assert.Equal(t, userID, extractedID)

	// Note: Testing actual expiration would require either:
	// 1. Mocking time (complex)
	// 2. Creating a service with very short expiration and waiting
	// 3. Manually crafting an expired token
	// For now, we just verify the token is valid when freshly created
}
