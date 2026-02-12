package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	UserID      string   `bson:"userId" json:"userId"`
	Type        UserType `bson:"type" json:"type"`
	DisplayName string   `bson:"displayName" json:"displayName"`

	About    string   `bson:"about,omitempty" json:"about,omitempty"`
	Location string   `bson:"location,omitempty" json:"location,omitempty"`
	Links    []string `bson:"links,omitempty" json:"links,omitempty"`

	Industry string `bson:"industry,omitempty" json:"industry,omitempty"`
	Website  string `bson:"website,omitempty" json:"website,omitempty"`

	CreatedAt time.Time `bson:"createdAt" json:"-"`
	UpdatedAt time.Time `bson:"updatedAt" json:"-"`

	AvatarURL string `bson:"avatarUrl,omitempty" json:"avatarUrl,omitempty"`
}
