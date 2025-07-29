package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    int                `bson:"user_id" json:"user_id"`
	Message   string             `bson:"message" json:"message"`
	Type      string             `bson:"type" json:"type"` // e.g., "order", "promo"
	IsRead    bool               `bson:"is_read" json:"is_read"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
}
