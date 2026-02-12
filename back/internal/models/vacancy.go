package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vacancy struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	VacancyID string `bson:"vacancyId" json:"vacancyId"`
	CompanyID string `bson:"companyId" json:"companyId"`

	Title       string   `bson:"title" json:"title"`
	Description string   `bson:"description" json:"description"`
	Location    string   `bson:"location,omitempty" json:"location,omitempty"`
	Tags        []string `bson:"tags,omitempty" json:"tags,omitempty"`

	IsPremium bool   `bson:"isPremium" json:"isPremium"`
	ColorCode string `bson:"colorCode,omitempty" json:"colorCode,omitempty"` // hex color for premium highlighting

	Status    string    `bson:"status" json:"status"` // active/closed
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"-"`
}
