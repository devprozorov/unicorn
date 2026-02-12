package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	UserID      string   `bson:"userId" json:"userId"`
	Login       string   `bson:"login" json:"login"`
	LoginNorm   string   `bson:"loginNorm" json:"-"`
	DisplayName string   `bson:"displayName" json:"displayName"`
	Type        UserType `bson:"type" json:"type"`

	PasswordHash string `bson:"passwordHash" json:"-"`

	Status struct {
		Deleted bool `bson:"deleted" json:"-"`
		Blocked bool `bson:"blocked" json:"-"`
	} `bson:"status" json:"-"`

	MFA struct {
		TOTP struct {
			Enabled        bool   `bson:"enabled" json:"-"`
			SecretEncB64   string `bson:"secretEncB64,omitempty" json:"-"`
			PendingEncB64  string `bson:"pendingEncB64,omitempty" json:"-"`
			PendingExpires int64  `bson:"pendingExpires,omitempty" json:"-"`
			LastStep       int64  `bson:"lastStep,omitempty" json:"-"`
		} `bson:"totp" json:"-"`
	} `bson:"mfa" json:"-"`

	Subscription struct {
		Active bool      `bson:"active" json:"-"`
		Until  time.Time `bson:"until,omitempty" json:"-"`
	} `bson:"subscription" json:"-"`

	CreatedAt time.Time `bson:"createdAt" json:"-"`
	UpdatedAt time.Time `bson:"updatedAt" json:"-"`
}
