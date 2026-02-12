package repo

import (
	"context"
	"time"

	"unicorn-auth/internal/db"
	"unicorn-auth/internal/models"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ApplicationRepo struct{ d *db.Database }

func NewApplicationRepo(d *db.Database) *ApplicationRepo { return &ApplicationRepo{d: d} }

func (r *ApplicationRepo) Create(ctx context.Context, a *models.Application) error {
	now := time.Now().UTC()
	a.ApplicationID = ulid.Make().String()
	a.Status = "pending"
	a.CreatedAt, a.UpdatedAt = now, now
	_, err := r.d.Applications().InsertOne(ctx, a)
	return err
}

func (r *ApplicationRepo) GetByID(ctx context.Context, appID string) (*models.Application, error) {
	var a models.Application
	err := r.d.Applications().FindOne(ctx, bson.M{"applicationId": appID}).Decode(&a)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &a, err
}

func (r *ApplicationRepo) ListInbox(ctx context.Context, companyID, status string, limit, skip int64) ([]models.Application, error) {
	f := bson.M{"companyId": companyID, "hidden.company": bson.M{"$ne": true}}
	if status != "" {
		f["status"] = status
	}
	cur, err := r.d.Applications().Find(ctx, f, options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []models.Application
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *ApplicationRepo) UpdateStatus(ctx context.Context, appID, companyID, status string) error {
	_, err := r.d.Applications().UpdateOne(ctx,
		bson.M{"applicationId": appID, "companyId": companyID},
		bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now().UTC()}},
	)
	return err
}

func (r *ApplicationRepo) ExistsCompanyResume(ctx context.Context, companyID, resumeID string) (bool, error) {
	n, err := r.d.Applications().CountDocuments(ctx, bson.M{
		"companyId": companyID,
		"resumeId":  resumeID,
	})
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *ApplicationRepo) ListByUser(ctx context.Context, userID, status string, limit, skip int64) ([]models.Application, error) {
	f := bson.M{"userId": userID, "hidden.user": bson.M{"$ne": true}}
	if status != "" {
		f["status"] = status
	}

	cur, err := r.d.Applications().Find(
		ctx,
		f,
		options.Find().
			SetLimit(limit).
			SetSkip(skip).
			SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Application
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
func (r *ApplicationRepo) MarkViewed(ctx context.Context, appID, companyID string) error {
	now := time.Now().UTC()
	_, err := r.d.Applications().UpdateOne(
		ctx,
		bson.M{"applicationId": appID, "companyId": companyID},
		bson.M{"$set": bson.M{
			"viewed":    true,
			"viewedAt":  now,
			"updatedAt": now,
		}},
	)
	return err
}

// ListByUserID возвращает все applications пользователя (кроме скрытых)
func (r *ApplicationRepo) ListByUserID(ctx context.Context, userID string) ([]models.Application, error) {
	cur, err := r.d.Applications().Find(
		ctx,
		bson.M{"userId": userID, "hidden.user": bson.M{"$ne": true}},
		options.Find().SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Application
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListByCompanyID возвращает все applications компании (кроме скрытых)
func (r *ApplicationRepo) ListByCompanyID(ctx context.Context, companyID string) ([]models.Application, error) {
	cur, err := r.d.Applications().Find(
		ctx,
		bson.M{"companyId": companyID, "hidden.company": bson.M{"$ne": true}},
		options.Find().SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Application
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Hide скрывает отклик для пользователя или компании
func (r *ApplicationRepo) Hide(ctx context.Context, appID, userID string, isCompany bool) error {
	field := "hidden.user"
	if isCompany {
		field = "hidden.company"
	}

	update := bson.M{
		"$set": bson.M{
			field:       true,
			"hiddenAt":  time.Now().UTC(),
			"updatedAt": time.Now().UTC(),
		},
	}

	filter := bson.M{"applicationId": appID}
	if isCompany {
		filter["companyId"] = userID
	} else {
		filter["userId"] = userID
	}

	_, err := r.d.Applications().UpdateOne(ctx, filter, update)
	return err
}

// UnhideOnNewMessage снимает скрытие при получении нового сообщения
func (r *ApplicationRepo) UnhideOnNewMessage(ctx context.Context, appID string) error {
	update := bson.M{
		"$set": bson.M{
			"hidden.user":    false,
			"hidden.company": false,
			"updatedAt":      time.Now().UTC(),
		},
	}

	_, err := r.d.Applications().UpdateOne(ctx, bson.M{"applicationId": appID}, update)
	return err
}

// DeleteOldHidden удаляет отклики, скрытые более месяца назад
func (r *ApplicationRepo) DeleteOldHidden(ctx context.Context) (int64, error) {
	oneMonthAgo := time.Now().UTC().AddDate(0, -1, 0)

	// Удаляем те, где оба скрыли и прошёл месяц
	filter := bson.M{
		"$and": []bson.M{
			{"hidden.user": true},
			{"hidden.company": true},
			{"hiddenAt": bson.M{"$lt": oneMonthAgo}},
		},
	}

	result, err := r.d.Applications().DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
