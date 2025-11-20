package user

import (
	"context"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"backend/internal/model"
)

func (s *ServiceSuite) TestCreateUser() {
	ctx := context.Background()
	passwordHash, err := argon2id.CreateHash("Very_strong_password1235", argon2id.DefaultParams)
	assert.Nil(s.T(), err)

	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func()
		expectedError error
	}{
		{
			name:     "Create user success",
			email:    "goodemail@gmail.com",
			password: passwordHash,
			mockSetup: func() {
				s.userRepository.On("Create", ctx, mock.MatchedBy(func(u model.User) bool {
					return u.Email == "goodemail@gmail.com"
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "Not unique email",
			email:    "checkemail@gmail.com",
			password: passwordHash,
			mockSetup: func() {
				s.userRepository.On("Create", ctx, mock.MatchedBy(func(u model.User) bool {
					return u.Email == "checkemail@gmail.com"
				})).Return(errors.New("такой email уже зарегистирован"))
			},
			expectedError: errors.New("такой email уже зарегистирован"),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			test.mockSetup()

			_, err = s.service.Create(ctx, test.email, test.password)

			if test.expectedError != nil {
				s.Error(err)
				s.Equal(test.expectedError.Error(), err.Error())
			} else {
				s.NoError(err)
			}
		})

		s.userRepository.AssertExpectations(s.T())

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
		expectedUser  model.User
		expectedError error
	}{
		{
			name: "успешное получение",
			id:   existingID,
			mockSetup: func() {
				s.userRepository.On("GetUserByID", ctx, existingID).Return(model.User{
					ID:    existingID,
					Email: "test@gmail.com",
				}, nil).Once()
			},
			expectedUser: model.User{
				ID:    existingID,
				Email: "test@gmail.com",
			},
			expectedError: nil,
		},
		{
			name: "не существующий айди",
			id:   nonExistingID,
			mockSetup: func() {
				s.userRepository.On("GetUserByID", ctx, nonExistingID).Return(model.User{}, errors.New("user not found")).Once()
			},
			expectedUser:  model.User{},
			expectedError: errors.New("user not found"),
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
		})

		s.userRepository.AssertExpectations(s.T())
	}
}

func (s *ServiceSuite) TestGetByEmail() {
	ctx := context.Background()

	tests := []struct {
		name          string
		email         string
		mockSetup     func()
		expectedUser  model.User
		expectedError error
	}{
		{
			name:  "успешное получение",
			email: "existing@gmail.com",
			mockSetup: func() {
				s.userRepository.On("GetUserByEmail", ctx, "existing@gmail.com").Return(model.User{
					ID:    "some-id",
					Email: "existing@gmail.com",
				}, nil).Once()
			},
			expectedUser: model.User{
				ID:    "some-id",
				Email: "existing@gmail.com",
			},
			expectedError: nil,
		},
		{
			name:  "не существующий email",
			email: "nonexisting@gmail.com",
			mockSetup: func() {
				s.userRepository.On("GetUserByEmail", ctx, "nonexisting@gmail.com").Return(model.User{}, errors.New("user not found")).Once()
			},
			expectedUser:  model.User{},
			expectedError: errors.New("user not found"),
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
		})

		s.userRepository.AssertExpectations(s.T())
	}
}
