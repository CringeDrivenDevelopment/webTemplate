package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/mock"

	"backend/internal/infra/queries"
	"backend/pkg/utils"
)

func (s *ServiceSuite) TestRegister() {
	ctx := context.Background()

	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func()
		expectedError error
		checkID       bool
	}{
		{
			name:     "successful registration",
			email:    "newuser@gmail.com",
			password: "SecurePassword123",
			mockSetup: func() {
				// Email doesn't exist (returns ErrNoRows)
				s.userRepository.On("GetUserByEmail", ctx, "newuser@gmail.com").
					Return(queries.User{}, pgx.ErrNoRows).Once()

				// Create succeeds
				s.userRepository.On("Create", ctx, mock.MatchedBy(func(u queries.User) bool {
					return u.Email == "newuser@gmail.com" && u.ID != "" && u.PasswordHash != ""
				})).Return(nil).Once()
			},
			expectedError: nil,
			checkID:       true,
		},
		{
			name:     "email already exists",
			email:    "existing@gmail.com",
			password: "SecurePassword123",
			mockSetup: func() {
				// Email already exists
				s.userRepository.On("GetUserByEmail", ctx, "existing@gmail.com").
					Return(queries.User{
						ID:    "existing-id",
						Email: "existing@gmail.com",
					}, nil).Once()
			},
			expectedError: utils.ErrEmailAlreadySignup,
			checkID:       false,
		},
		{
			name:     "database error on email check",
			email:    "test@gmail.com",
			password: "SecurePassword123",
			mockSetup: func() {
				// Database error
				s.userRepository.On("GetUserByEmail", ctx, "test@gmail.com").
					Return(queries.User{}, errors.New("database connection error")).Once()
			},
			expectedError: errors.New("database connection error"),
			checkID:       false,
		},
		{
			name:     "database error on create",
			email:    "newuser2@gmail.com",
			password: "SecurePassword123",
			mockSetup: func() {
				// Email doesn't exist
				s.userRepository.On("GetUserByEmail", ctx, "newuser2@gmail.com").
					Return(queries.User{}, pgx.ErrNoRows).Once()

				// Create fails
				s.userRepository.On("Create", ctx, mock.MatchedBy(func(u queries.User) bool {
					return u.Email == "newuser2@gmail.com"
				})).Return(errors.New("insert failed")).Once()
			},
			expectedError: errors.New("insert failed"),
			checkID:       false,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			test.mockSetup()

			id, err := s.service.Register(ctx, test.email, test.password)

			if test.expectedError != nil {
				s.Error(err)
				s.Equal(test.expectedError.Error(), err.Error())
				s.Empty(id)
			} else {
				s.NoError(err)
				if test.checkID {
					s.NotEmpty(id)
					// Verify it's a valid ULID
					_, parseErr := ulid.Parse(id)
					s.NoError(parseErr)
				}
			}

			s.userRepository.AssertExpectations(s.T())
		})
	}
}

func (s *ServiceSuite) TestGetByID() {
	ctx := context.Background()

	existingID := ulid.Make().String()
	nonExistingID := ulid.Make().String()

	tests := []struct {
		name          string
		id            string
		mockSetup     func()
		expectedUser  queries.User
		expectedError error
	}{
		{
			name: "successful retrieval",
			id:   existingID,
			mockSetup: func() {
				s.userRepository.On("GetUserByID", ctx, existingID).Return(queries.User{
					ID:    existingID,
					Email: "test@gmail.com",
				}, nil).Once()
			},
			expectedUser: queries.User{
				ID:    existingID,
				Email: "test@gmail.com",
			},
			expectedError: nil,
		},
		{
			name: "user not found",
			id:   nonExistingID,
			mockSetup: func() {
				s.userRepository.On("GetUserByID", ctx, nonExistingID).
					Return(queries.User{}, pgx.ErrNoRows).Once()
			},
			expectedUser:  queries.User{},
			expectedError: pgx.ErrNoRows,
		},
		{
			name: "database error",
			id:   "some-id",
			mockSetup: func() {
				s.userRepository.On("GetUserByID", ctx, "some-id").
					Return(queries.User{}, errors.New("connection timeout")).Once()
			},
			expectedUser:  queries.User{},
			expectedError: errors.New("connection timeout"),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			test.mockSetup()

			user, err := s.service.GetByID(ctx, test.id)

			if test.expectedError != nil {
				s.Error(err)
				s.Equal(test.expectedError.Error(), err.Error())
			} else {
				s.NoError(err)
				s.Equal(test.expectedUser.ID, user.ID)
				s.Equal(test.expectedUser.Email, user.Email)
			}

			s.userRepository.AssertExpectations(s.T())
		})
	}
}

func (s *ServiceSuite) TestGetByEmail() {
	ctx := context.Background()

	tests := []struct {
		name          string
		email         string
		mockSetup     func()
		expectedUser  queries.User
		expectedError error
	}{
		{
			name:  "successful retrieval",
			email: "existing@gmail.com",
			mockSetup: func() {
				s.userRepository.On("GetUserByEmail", ctx, "existing@gmail.com").Return(queries.User{
					ID:    "some-id",
					Email: "existing@gmail.com",
				}, nil).Once()
			},
			expectedUser: queries.User{
				ID:    "some-id",
				Email: "existing@gmail.com",
			},
			expectedError: nil,
		},
		{
			name:  "user not found",
			email: "nonexisting@gmail.com",
			mockSetup: func() {
				s.userRepository.On("GetUserByEmail", ctx, "nonexisting@gmail.com").
					Return(queries.User{}, pgx.ErrNoRows).Once()
			},
			expectedUser:  queries.User{},
			expectedError: pgx.ErrNoRows,
		},
		{
			name:  "database error",
			email: "error@gmail.com",
			mockSetup: func() {
				s.userRepository.On("GetUserByEmail", ctx, "error@gmail.com").
					Return(queries.User{}, errors.New("database error")).Once()
			},
			expectedUser:  queries.User{},
			expectedError: errors.New("database error"),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			test.mockSetup()

			user, err := s.service.GetByEmail(ctx, test.email)

			if test.expectedError != nil {
				s.Error(err)
				s.Equal(test.expectedError.Error(), err.Error())
			} else {
				s.NoError(err)
				s.Equal(test.expectedUser.Email, user.Email)
			}

			s.userRepository.AssertExpectations(s.T())
		})
	}
}
