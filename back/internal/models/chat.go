package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatMessage struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	MessageID     string `bson:"messageId" json:"messageId"`
	ApplicationID string `bson:"applicationId" json:"applicationId"`

	SenderID   string `bson:"senderId" json:"senderId"`
	SenderType string `bson:"senderType" json:"senderType"` // user/company
	Text       string `bson:"text" json:"text"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
