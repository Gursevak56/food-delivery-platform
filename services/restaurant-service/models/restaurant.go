package models

import (
	"encoding/json"
	"time"
)

type Restaurant struct {
	ID              int64           `json:"id"`
	OwnerAuthUserID *int64          `json:"owner_auth_user_id,omitempty"`
	Name            string          `json:"name"`
	Slug            string          `json:"slug,omitempty"`
	Description     string          `json:"description,omitempty"`
	Status          string          `json:"status,omitempty"`
	AddressLine1    string          `json:"address_line1,omitempty"`
	AddressLine2    string          `json:"address_line2,omitempty"`
	City            string          `json:"city,omitempty"`
	State           string          `json:"state,omitempty"`
	Pincode         string          `json:"pincode,omitempty"`
	Latitude        *float64        `json:"latitude,omitempty"`
	Longitude       *float64        `json:"longitude,omitempty"`
	AvgRating       *float64        `json:"avg_rating,omitempty"`
	RatingCount     *int64          `json:"rating_count,omitempty"`
	Tags            []string        `json:"tags,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	CreatedAt       *time.Time      `json:"created_at,omitempty"`
	UpdatedAt       *time.Time      `json:"updated_at,omitempty"`
}

type RestaurantHour struct {
	ID           int64  `json:"id"`
	RestaurantID int64  `json:"restaurant_id"`
	Weekday      int    `json:"weekday"`              // 0..6
	OpenTime     string `json:"open_time,omitempty"`  // "15:04:05"
	CloseTime    string `json:"close_time,omitempty"` // "15:04:05"
	IsClosed     bool   `json:"is_closed,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type RestaurantTable struct {
	ID              int64  `json:"id"`
	RestaurantID    int64  `json:"restaurant_id"`
	TableIdentifier string `json:"table_identifier"`
	Seats           int    `json:"seats"`
	QRToken         string `json:"qr_token,omitempty"`
	QRUrl           string `json:"qr_url,omitempty"`
	IsActive        bool   `json:"is_active"`
	CreatedAt       string `json:"created_at,omitempty"`
}
