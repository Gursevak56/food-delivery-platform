package models

import "time"

type User struct {
	ID            string    `json:"id"`
	Email         *string   `json:"email"`
	Phone         *string   `json:"phone"`
	FullName      *string   `json:"full_name"`
	UserType      *string   `json:"user_type"`
	IsActive      *bool     `json:"is_active"`
	PhoneVerified *bool     `json:"phone_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
