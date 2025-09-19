package models

import "time"

// Location represents a saved location record
type Location struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string    `bson:"userId,omitempty" json:"userId,omitempty"` // optional user identifier from Flutter
	Lat       float64   `bson:"lat" json:"lat"`
	Lng       float64   `bson:"lng" json:"lng"`
	Address   string    `bson:"address,omitempty" json:"address,omitempty"`
	Accuracy  float64   `bson:"accuracy,omitempty" json:"accuracy,omitempty"`
	Source    string    `bson:"source" json:"source"` // "ip" | "client" | "manual"
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
