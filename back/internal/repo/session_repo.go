package repo

import (
	"context"
	"time"

	"unicorn-auth/internal/db"
	"unicorn-auth/internal/models"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepo struct{ d *db.Database }

func NewSessionRepo(d *db.Database) *SessionRepo { return &SessionRepo{d: d} }

func (r *SessionRepo) Create(ctx context.Context, userID string, refreshHash string, ttl time.Duration) (*models.Session, error) {
	now := time.Now().UTC()
	s := &models.Session{
		SessionID:    ulid.Make().String(),
		UserID:       userID,
		RefreshHash:  refreshHash,
		Revoked:      false,
		CreatedAt:    now,
		ExpiresAt:    now.Add(ttl),
	}
	_, err := r.d.Sessions().InsertOne(ctx, s)
	return s, err
}

func (r *SessionRepo) FindActiveByID(ctx context.Context, sessionID string) (*models.Session, error) {
	var s models.Session
	err := r.d.Sessions().FindOne(ctx, bson.M{"sessionId": sessionID, "revoked": false}).Decode(&s)
	if err == mongo.ErrNoDocuments { return nil, nil }
	return &s, err
}

func (r *SessionRepo) Rotate(ctx context.Context, sessionID string, newHash string, newExp time.Time) error {
	_, err := r.d.Sessions().UpdateOne(ctx, bson.M{"sessionId": sessionID, "revoked": false},
		bson.M{"$set": bson.M{"refreshHash": newHash, "expiresAt": newExp}},
	)
	return err
}

func (r *SessionRepo) Revoke(ctx context.Context, sessionID string) error {
	_, err := r.d.Sessions().UpdateOne(ctx, bson.M{"sessionId": sessionID}, bson.M{"$set": bson.M{"revoked": true}})
	return err
}
