package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Application struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	ApplicationID string `bson:"applicationId" json:"applicationId"`
	VacancyID     string `bson:"vacancyId" json:"vacancyId"`
	ResumeID      string `bson:"resumeId" json:"resumeId"`

	UserID    string `bson:"userId" json:"userId"`
	CompanyID string `bson:"companyId" json:"companyId"`

	Status  string `bson:"status" json:"status"` // pending/accepted/rejected
	Message string `bson:"message,omitempty" json:"message,omitempty"`

	// Hidden state: скрыто ли для пользователя или компании
	Hidden   Hidden    `bson:"hidden" json:"-"`
	HiddenAt time.Time `bson:"hiddenAt,omitempty" json:"-"`

	CreatedAt time.Time `bson:"createdAt" json:"-"`
	UpdatedAt time.Time `bson:"updatedAt" json:"-"`

	Viewed   bool      `bson:"viewed" json:"viewed"`
	ViewedAt time.Time `bson:"viewedAt,omitempty" json:"viewedAt,omitempty"`
}

// Hidden представляет состояние скрытия для пользователя и компании
type Hidden struct {
	User    bool `bson:"user" json:"user"`       // скрыто пользователем
	Company bool `bson:"company" json:"company"` // скрыто компанией
}
