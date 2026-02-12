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

type VacancyRepo struct{ d *db.Database }

func NewVacancyRepo(d *db.Database) *VacancyRepo { return &VacancyRepo{d: d} }

// VacancyPublic дополняет вакансию вычисляемым количеством откликов.
type VacancyPublic struct {
	models.Vacancy `bson:",inline"`
	ResponsesCount int64 `bson:"responsesCount" json:"responsesCount"`
}

func (r *VacancyRepo) CountActiveByCompany(ctx context.Context, companyID string) (int64, error) {
	return r.d.Vacancies().CountDocuments(ctx, bson.M{"companyId": companyID, "status": "active"})
}

func (r *VacancyRepo) Create(ctx context.Context, v *models.Vacancy) error {
	now := time.Now().UTC()
	v.VacancyID = ulid.Make().String()
	v.Status = "active"
	v.CreatedAt, v.UpdatedAt = now, now
	_, err := r.d.Vacancies().InsertOne(ctx, v)
	return err
}

func (r *VacancyRepo) GetByID(ctx context.Context, vacancyID string) (*models.Vacancy, error) {
	var v models.Vacancy
	err := r.d.Vacancies().FindOne(ctx, bson.M{"vacancyId": vacancyID}).Decode(&v)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &v, err
}

func (r *VacancyRepo) Update(ctx context.Context, vacancyID, companyID string, set bson.M) error {
	set["updatedAt"] = time.Now().UTC()
	_, err := r.d.Vacancies().UpdateOne(ctx, bson.M{"vacancyId": vacancyID, "companyId": companyID}, bson.M{"$set": set})
	return err
}

func (r *VacancyRepo) Delete(ctx context.Context, vacancyID, companyID string) error {
	_, err := r.d.Vacancies().DeleteOne(ctx, bson.M{"vacancyId": vacancyID, "companyId": companyID})
	return err
}

func (r *VacancyRepo) ListPublic(ctx context.Context, limit, skip int64) ([]VacancyPublic, error) {
	cutoff := time.Now().UTC().Add(-30 * 24 * time.Hour)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"status": "active", "createdAt": bson.M{"$gte": cutoff}}}},
		// Сортировка: сначала премиум (isPremium: true), потом по дате создания
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "isPremium", Value: -1}, // -1 = true первые
			{Key: "createdAt", Value: -1}, // новые первые
		}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "applications",
			"let":  bson.M{"vacancyId": "$vacancyId"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$eq": bson.A{"$vacancyId", "$$vacancyId"}},
				}}},
				bson.D{{Key: "$count", Value: "count"}},
			},
			"as": "appsCount",
		}}},
		bson.D{{Key: "$addFields", Value: bson.M{
			"responsesCount": bson.M{
				"$ifNull": bson.A{bson.M{"$first": "$appsCount.count"}, 0},
			},
		}}},
		bson.D{{Key: "$project", Value: bson.M{"appsCount": 0}}},
	}

	cur, err := r.d.Vacancies().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []VacancyPublic
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CountActive возвращает количество активных вакансий (статус "active")
func (r *VacancyRepo) CountActive(ctx context.Context) (int64, error) {
	return r.d.Vacancies().CountDocuments(ctx, bson.M{"status": "active"})
}

// UpdateAllByCompanyID обновляет все вакансии компании
func (r *VacancyRepo) UpdateAllByCompanyID(ctx context.Context, companyID string, set bson.M) error {
	set["updatedAt"] = time.Now().UTC()
	_, err := r.d.Vacancies().UpdateMany(ctx,
		bson.M{"companyId": companyID, "status": "active"},
		bson.M{"$set": set})
	return err
}

// ListByCompanyID возвращает все вакансии компании
func (r *VacancyRepo) ListByCompanyID(ctx context.Context, companyID string) ([]models.Vacancy, error) {
	cur, err := r.d.Vacancies().Find(
		ctx,
		bson.M{"companyId": companyID},
		options.Find().SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []models.Vacancy
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
