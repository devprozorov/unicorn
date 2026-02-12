package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	AdminID       string `bson:"adminId" json:"adminId"`
	Login         string `bson:"login" json:"login"`
	LoginNorm     string `bson:"loginNorm" json:"-"`
	PasswordHash  string `bson:"passwordHash" json:"-"`
	CreatedAt     time.Time `bson:"createdAt" json:"-"`
}
