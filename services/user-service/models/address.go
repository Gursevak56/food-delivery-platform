package models

import "time"

type Address struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id"`
	Label        string     `json:"label,omitempty"`
	AddressLine1 string     `json:"address_line1,omitempty"`
	AddressLine2 string     `json:"address_line2,omitempty"`
	City         string     `json:"city,omitempty"`
	State        string     `json:"state,omitempty"`
	Pincode      string     `json:"pincode,omitempty"`
	Latitude     *float64   `json:"latitude,omitempty"`
	Longitude    *float64   `json:"longitude,omitempty"`
	IsDefault    bool       `json:"is_default,omitempty"`
	Metadata     string     `json:"metadata,omitempty"` // store JSON string or use sql.NullString / json.RawMessage if preferred
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}
