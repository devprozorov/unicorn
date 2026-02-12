package repo

import (
	"context"
	"time"

	"unicorn-auth/internal/db"
	"unicorn-auth/internal/models"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepo struct{ d *db.Database }

func NewChatRepo(d *db.Database) *ChatRepo { return &ChatRepo{d: d} }

func (r *ChatRepo) Create(ctx context.Context, m *models.ChatMessage) error {
	m.MessageID = ulid.Make().String()
	m.CreatedAt = time.Now().UTC()
	_, err := r.d.ChatMessages().InsertOne(ctx, m)
	return err
}

func (r *ChatRepo) List(ctx context.Context, appID string, limit int64) ([]models.ChatMessage, error) {
	cur, err := r.d.ChatMessages().Find(ctx, bson.M{"applicationId": appID},
		options.Find().SetLimit(limit).SetSort(bson.M{"createdAt": 1}),
	)
	if err != nil { return nil, err }
	defer cur.Close(ctx)
	var out []models.ChatMessage
	if err := cur.All(ctx, &out); err != nil { return nil, err }
	return out, nil
}
