package models

import (
	"encoding/json"
	"time"
)

type MenuCategory struct {
	ID           int64           `json:"id"`
	RestaurantID int64           `json:"restaurant_id"`
	Name         string          `json:"name"`
	Slug         string          `json:"slug,omitempty"`
	ParentID     *int64          `json:"parent_id,omitempty"`
	SortOrder    int             `json:"sort_order,omitempty"`
	IsActive     bool            `json:"is_active,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	CreatedAt    *time.Time      `json:"created_at,omitempty"`
}

type MenuItem struct {
	ID              int64           `json:"id"`
	RestaurantID    int64           `json:"restaurant_id"`
	CategoryID      *int64          `json:"category_id,omitempty"`
	Name            string          `json:"name"`
	Description     string          `json:"description,omitempty"`
	Price           float64         `json:"price"`
	Currency        string          `json:"currency,omitempty"`     // e.g. "INR"
	Availability    string          `json:"availability,omitempty"` // e.g. "IN_STOCK"
	IsVeg           bool            `json:"is_veg,omitempty"`
	SpiceLevel      int             `json:"spice_level,omitempty"`
	PrepTimeMinutes int             `json:"prep_time_minutes,omitempty"`
	Tags            []string        `json:"tags,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`  // free-form json (ingredients etc)
	ImageURL        string          `json:"image_url,omitempty"` // optional
	CreatedAt       *time.Time      `json:"created_at,omitempty"`
	UpdatedAt       *time.Time      `json:"updated_at,omitempty"`
}
