package repo

import (
	"context"
	"time"

	"unicorn-auth/internal/db"
	"unicorn-auth/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileRepo struct{ d *db.Database }

func NewProfileRepo(d *db.Database) *ProfileRepo { return &ProfileRepo{d: d} }

func (r *ProfileRepo) GetByUserID(ctx context.Context, userID string) (*models.Profile, error) {
	var p models.Profile
	err := r.d.Profiles().FindOne(ctx, bson.M{"userId": userID}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &p, err
}

func (r *ProfileRepo) Create(ctx context.Context, p *models.Profile) error {
	now := time.Now().UTC()
	p.CreatedAt, p.UpdatedAt = now, now
	_, err := r.d.Profiles().InsertOne(ctx, p)
	return err
}

func (r *ProfileRepo) UpdateByUserID(ctx context.Context, userID string, set bson.M) error {
	set["updatedAt"] = time.Now().UTC()
	_, err := r.d.Profiles().UpdateOne(ctx, bson.M{"userId": userID}, bson.M{"$set": set})
	return err
}

func (r *ProfileRepo) SearchCompanies(ctx context.Context, filter bson.M, limit int64) ([]models.Profile, error) {
	cur, err := r.d.Profiles().Find(ctx, filter)
	if err != nil { return nil, err }
	defer cur.Close(ctx)
	var out []models.Profile
	if err := cur.All(ctx, &out); err != nil { return nil, err }
	if int64(len(out)) > limit {
		out = out[:limit]
	}
	return out, nil
}
