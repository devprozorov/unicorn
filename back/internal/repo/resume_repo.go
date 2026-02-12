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

type ResumeRepo struct{ d *db.Database }

func NewResumeRepo(d *db.Database) *ResumeRepo { return &ResumeRepo{d: d} }

func (r *ResumeRepo) CountActiveByUser(ctx context.Context, userID string) (int64, error) {
	return r.d.Resumes().CountDocuments(ctx, bson.M{"userId": userID, "status": "active"})
}

func (r *ResumeRepo) Create(ctx context.Context, rr *models.Resume) error {
	now := time.Now().UTC()
	rr.ResumeID = ulid.Make().String()
	rr.Status = "active"
	rr.CreatedAt, rr.UpdatedAt = now, now
	_, err := r.d.Resumes().InsertOne(ctx, rr)
	return err
}

func (r *ResumeRepo) GetByID(ctx context.Context, resumeID string) (*models.Resume, error) {
	var rr models.Resume
	err := r.d.Resumes().FindOne(ctx, bson.M{"resumeId": resumeID}).Decode(&rr)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &rr, err
}

func (r *ResumeRepo) ListMine(ctx context.Context, userID string) ([]models.Resume, error) {
	cur, err := r.d.Resumes().Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []models.Resume
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *ResumeRepo) Update(ctx context.Context, resumeID, userID string, set bson.M) error {
	set["updatedAt"] = time.Now().UTC()
	_, err := r.d.Resumes().UpdateOne(ctx, bson.M{"resumeId": resumeID, "userId": userID}, bson.M{"$set": set})
	return err
}

func (r *ResumeRepo) Delete(ctx context.Context, resumeID, userID string) error {
	_, err := r.d.Resumes().DeleteOne(ctx, bson.M{"resumeId": resumeID, "userId": userID})
	return err
}

// CountActive возвращает количество активных резюме (статус "active")
func (r *ResumeRepo) CountActive(ctx context.Context) (int64, error) {
	return r.d.Resumes().CountDocuments(ctx, bson.M{"status": "active"})
}

// UpdateAllByUserID обновляет все резюме пользователя
func (r *ResumeRepo) UpdateAllByUserID(ctx context.Context, userID string, set bson.M) error {
	set["updatedAt"] = time.Now().UTC()
	_, err := r.d.Resumes().UpdateMany(ctx,
		bson.M{"userId": userID, "status": "active"},
		bson.M{"$set": set})
	return err
}
