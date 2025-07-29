package models

import "time"

type User struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	UserType      string    `json:"user_type"`
	IsActive      bool      `json:"is_active"`
	EmailVerified bool      `json:"email_verified"`
	PhoneVerified bool      `json:"phone_verified"`
	Password      string    `json:"password,omitempty"` // only used on create, omitted in JSON after
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
