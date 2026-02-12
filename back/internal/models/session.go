package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	SessionID string `bson:"sessionId" json:"sessionId"`
	UserID    string `bson:"userId" json:"userId"`

	RefreshHash string `bson:"refreshHash" json:"-"`
	Revoked     bool   `bson:"revoked" json:"-"`
	CreatedAt   time.Time `bson:"createdAt" json:"-"`
	ExpiresAt   time.Time `bson:"expiresAt" json:"-"`
}
