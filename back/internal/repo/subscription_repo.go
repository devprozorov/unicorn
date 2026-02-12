package repo

import (
	"context"
	"time"

	"unicorn-auth/internal/db"
	"unicorn-auth/internal/models"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionRepo struct{ d *db.Database }

func NewSubscriptionRepo(d *db.Database) *SubscriptionRepo {
	return &SubscriptionRepo{d: d}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s *models.Subscription) error {
	now := time.Now().UTC()
	s.SubscriptionID = ulid.Make().String()
	s.CreatedAt = now
	s.UpdatedAt = now
	_, err := r.d.Subscriptions().InsertOne(ctx, s)
	return err
}

func (r *SubscriptionRepo) FindByInvID(ctx context.Context, invID int64) (*models.Subscription, error) {
	var s models.Subscription
	err := r.d.Subscriptions().FindOne(ctx, bson.M{"invId": invID}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &s, err
}

func (r *SubscriptionRepo) UpdateStatus(ctx context.Context, invID int64, status string, startDate, endDate time.Time) error {
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"startDate": startDate,
			"endDate":   endDate,
			"updatedAt": time.Now().UTC(),
		},
	}
	_, err := r.d.Subscriptions().UpdateOne(ctx, bson.M{"invId": invID}, update)
	return err
}

func (r *SubscriptionRepo) GetActiveByUserID(ctx context.Context, userID string) (*models.Subscription, error) {
	var s models.Subscription
	filter := bson.M{
		"userId":  userID,
		"status":  "paid",
		"endDate": bson.M{"$gt": time.Now().UTC()},
	}
	err := r.d.Subscriptions().FindOne(ctx, filter).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &s, err
}
