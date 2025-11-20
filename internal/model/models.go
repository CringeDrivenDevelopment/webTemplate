package model

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
}

// по идее надо бы для каждого слоя свой model файл делать и задрачивать конвертеры для каждого
// но мы заебемся это делать

func (*User) TableName() string {
	return "users"
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	var user User
	err := tx.Where("LOWER(users.email) = ?", strings.ToLower(u.Email)).First(&user).Error
	if err == nil {
		return fmt.Errorf("user with that email already exists")
	}
	return nil
}
