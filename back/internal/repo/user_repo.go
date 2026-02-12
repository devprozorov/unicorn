package repo

import (
	"context"
	"strings"
	"time"

	"unicorn-auth/internal/db"
	"unicorn-auth/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct{ d *db.Database }

func NewUserRepo(d *db.Database) *UserRepo { return &UserRepo{d: d} }

func normLogin(login string) string { return strings.ToLower(strings.TrimSpace(login)) }

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	now := time.Now().UTC()
	u.CreatedAt, u.UpdatedAt = now, now
	u.LoginNorm = normLogin(u.Login)
	_, err := r.d.Users().InsertOne(ctx, u)
	return err
}

func (r *UserRepo) FindByLoginNorm(ctx context.Context, loginNorm string) (*models.User, error) {
	var u models.User
	err := r.d.Users().FindOne(ctx, bson.M{"loginNorm": loginNorm}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepo) FindByUserID(ctx context.Context, userID string) (*models.User, error) {
	var u models.User
	err := r.d.Users().FindOne(ctx, bson.M{"userId": userID}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepo) UpdateByUserID(ctx context.Context, userID string, set bson.M) error {
	set["updatedAt"] = time.Now().UTC()
	_, err := r.d.Users().UpdateOne(ctx, bson.M{"userId": userID}, bson.M{"$set": set})
	return err
}

func (r *UserRepo) SoftDelete(ctx context.Context, userID string) error {
	_, err := r.d.Users().UpdateOne(ctx, bson.M{"userId": userID}, bson.M{"$set": bson.M{"status.deleted": true, "updatedAt": time.Now().UTC()}})
	return err
}

func (r *UserRepo) Block(ctx context.Context, userID string, blocked bool) error {
	_, err := r.d.Users().UpdateOne(ctx, bson.M{"userId": userID}, bson.M{"$set": bson.M{"status.blocked": blocked, "updatedAt": time.Now().UTC()}})
	return err
}

func (r *UserRepo) List(ctx context.Context, typ string, limit, skip int64) ([]models.User, error) {
	filter := bson.M{}
	if typ != "" {
		filter["type"] = typ
	}
	cur, err := r.d.Users().Find(ctx, filter, options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []models.User
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListWithFilters возвращает пользователей с различными фильтрами
func (r *UserRepo) ListWithFilters(ctx context.Context, typ, search, blocked, deleted, premium string, limit, skip int64) ([]models.User, error) {
	filter := bson.M{}

	if typ != "" {
		filter["type"] = typ
	}

	// Поиск по логину или displayName
	if search != "" {
		searchNorm := strings.ToLower(search)
		filter["$or"] = []bson.M{
			{"loginNorm": bson.M{"$regex": searchNorm, "$options": "i"}},
			{"displayName": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Фильтр по статусу блокировки
	if blocked == "true" {
		filter["status.blocked"] = true
	} else if blocked == "false" {
		filter["status.blocked"] = false
	}

	// Фильтр по статусу удаления
	if deleted == "true" {
		filter["status.deleted"] = true
	} else if deleted == "false" {
		filter["status.deleted"] = false
	}

	// Фильтр по премиум-подписке
	if premium == "true" {
		filter["subscription.active"] = true
	} else if premium == "false" {
		filter["subscription.active"] = false
	}

	cur, err := r.d.Users().Find(ctx, filter, options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []models.User
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CountAll возвращает общее количество пользователей (не удаленные)
func (r *UserRepo) CountAll(ctx context.Context) (int64, error) {
	return r.d.Users().CountDocuments(ctx, bson.M{"status.deleted": false})
}

// CountByType возвращает количество пользователей по типу (не удаленные)
func (r *UserRepo) CountByType(ctx context.Context, typ string) (int64, error) {
	return r.d.Users().CountDocuments(ctx, bson.M{"type": typ, "status.deleted": false})
}
