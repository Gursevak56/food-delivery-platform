package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Gursevak56/food-delivery-platform/services/user-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OTPRepo interface {
	Save(ctx context.Context, otp *models.OTP) error
	FindValid(ctx context.Context, phone, code string) (*models.OTP, error)
	MarkUsed(ctx context.Context, id interface{}) error
	IncrementAttempts(ctx context.Context, id interface{}) error
}

type otpRepo struct {
	col *mongo.Collection
}

func NewOTPRepo(col *mongo.Collection) OTPRepo {
	return &otpRepo{col: col}
}

// EnsureOTPIndexes creates phone index + TTL on expires_at
func EnsureOTPIndexes(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "phone", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	})
	return err
}

func (r *otpRepo) Save(ctx context.Context, otp *models.OTP) error {
	if otp == nil {
		return errors.New("otp required")
	}
	_, err := r.col.InsertOne(ctx, otp)
	return err
}

func (r *otpRepo) FindValid(ctx context.Context, phone, code string) (*models.OTP, error) {
	filter := bson.M{
		"phone": phone,
		"code":  code,
		"used":  false,
		"expires_at": bson.M{
			"$gt": time.Now().UTC(),
		},
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	var otp models.OTP
	err := r.col.FindOne(ctx, filter, opts).Decode(&otp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepo) MarkUsed(ctx context.Context, id interface{}) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"used": true}})
	return err
}

func (r *otpRepo) IncrementAttempts(ctx context.Context, id interface{}) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"attempts": 1}})
	return err
}
