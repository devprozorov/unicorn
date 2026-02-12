package repo

import (
	"context"
	"strings"
	"time"

	"unicorn-auth/internal/db"
	"unicorn-auth/internal/models"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepo struct{ d *db.Database }

func NewAdminRepo(d *db.Database) *AdminRepo { return &AdminRepo{d: d} }

func normAdminLogin(login string) string { return strings.ToLower(strings.TrimSpace(login)) }

func (r *AdminRepo) Create(ctx context.Context, login, passwordHash string) (*models.Admin, error) {
	a := &models.Admin{
		AdminID:      ulid.Make().String(),
		Login:        strings.TrimSpace(login),
		LoginNorm:    normAdminLogin(login),
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}
	_, err := r.d.Admins().InsertOne(ctx, a)
	return a, err
}

func (r *AdminRepo) FindByLoginNorm(ctx context.Context, loginNorm string) (*models.Admin, error) {
	var a models.Admin
	err := r.d.Admins().FindOne(ctx, bson.M{"loginNorm": loginNorm}).Decode(&a)
	if err == mongo.ErrNoDocuments { return nil, nil }
	return &a, err
}
