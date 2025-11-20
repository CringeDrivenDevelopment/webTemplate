package user

import (
	"context"
	"strings"

	"backend/internal/model"
)

func (r *Repository) Create(ctx context.Context, user model.User) error {
	err := r.db.WithContext(ctx).Create(&user).Error
	return err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("LOWER(users.email) = ?", strings.ToLower(email)).First(&user).Error
	return user, err
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, &model.User{ID: id}).Error
	return user, err
}
