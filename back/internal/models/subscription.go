package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscription struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	SubscriptionID string `bson:"subscriptionId" json:"subscriptionId"`
	UserID         string `bson:"userId" json:"userId"`

	Amount   float64 `bson:"amount" json:"amount"`
	Currency string  `bson:"currency" json:"currency"`

	Status string `bson:"status" json:"status"` // pending/paid/cancelled

	// Robokassa fields
	InvID  int64  `bson:"invId" json:"invId"`
	OutSum string `bson:"outSum" json:"outSum"`

	StartDate time.Time `bson:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate   time.Time `bson:"endDate,omitempty" json:"endDate,omitempty"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
