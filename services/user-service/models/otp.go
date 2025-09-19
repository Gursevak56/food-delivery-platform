package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OTP stored in Mongo
type OTP struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Phone     string             `bson:"phone" json:"phone"`
	Code      string             `bson:"code" json:"code"`
	Used      bool               `bson:"used" json:"used"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	Attempts  int                `bson:"attempts" json:"attempts"`
	Source    string             `bson:"source,omitempty" json:"source,omitempty"`
}
