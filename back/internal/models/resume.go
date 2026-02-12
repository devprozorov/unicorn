package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Resume struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	ResumeID string `bson:"resumeId" json:"resumeId"`
	UserID   string `bson:"userId" json:"userId"`

	Title  string   `bson:"title" json:"title"`
	About  string   `bson:"about" json:"about"`
	Skills []string `bson:"skills,omitempty" json:"skills,omitempty"`
	Links  []string `bson:"links,omitempty" json:"links,omitempty"`

	IsPremium bool   `bson:"isPremium" json:"isPremium"`
	ColorCode string `bson:"colorCode,omitempty" json:"colorCode,omitempty"` // hex color for premium highlighting

	Status    string    `bson:"status" json:"status"` // active/hidden
	CreatedAt time.Time `bson:"createdAt" json:"-"`
	UpdatedAt time.Time `bson:"updatedAt" json:"-"`
}
