package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey" json:"id"`
	UserName     string `json:"user_name"`
	Email        string `gorm:"unique" json:"email"`
	PasswordHash string `json:"-"` // პაროლის ჰეში არ შედის JSON პასუხებში
	ProfileImage []byte `gorm:"type:longblob" json:"profile_image,omitempty"`
}

type UserResponse struct {
	ID           uint   `json:"id"`
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	ProfileImage string `json:"profile_image,omitempty"` // Base64 დაშიფრულია
}

type RegisterRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Image    []byte `json:"image,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
