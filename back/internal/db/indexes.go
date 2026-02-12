package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *Database) EnsureIndexes(ctx context.Context) {
	must := func(err error) {
		if err != nil {
			log.Fatalf("ensure indexes: %v", err)
		}
	}

	_, err := d.Users().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "loginNorm", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_loginNorm")},
		{Keys: bson.D{{Key: "userId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_userId")},
	})
	must(err)

	_, err = d.Sessions().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "sessionId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_sessionId")},
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "revoked", Value: 1}}, Options: options.Index().SetName("sess_user_revoked")},
		{Keys: bson.D{{Key: "expiresAt", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0).SetName("ttl_sessions")},
	})
	must(err)

	_, err = d.Profiles().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_profile_userId")},
		{Keys: bson.D{{Key: "type", Value: 1}, {Key: "displayName", Value: 1}}, Options: options.Index().SetName("profile_type_name")},
		{Keys: bson.D{{Key: "industry", Value: 1}}, Options: options.Index().SetName("profile_industry")},
		{Keys: bson.D{{Key: "location", Value: 1}}, Options: options.Index().SetName("profile_location")},
	})
	must(err)

	_, err = d.Vacancies().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "vacancyId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_vacancyId")},
		{Keys: bson.D{{Key: "companyId", Value: 1}, {Key: "createdAt", Value: -1}}, Options: options.Index().SetName("vac_company_created")},
		{Keys: bson.D{{Key: "createdAt", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(int32(60 * 60 * 24 * 30)).SetName("ttl_vacancies_30d")},
	})
	must(err)

	_, err = d.Resumes().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "resumeId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_resumeId")},
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}, Options: options.Index().SetName("res_user_created")},
	})
	must(err)

	_, err = d.Applications().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "applicationId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_applicationId")},
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "vacancyId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_user_vacancy")},
		{Keys: bson.D{{Key: "companyId", Value: 1}, {Key: "createdAt", Value: -1}}, Options: options.Index().SetName("app_company_created")},
	})
	must(err)

	_, err = d.ChatMessages().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "applicationId", Value: 1}, {Key: "createdAt", Value: 1}}, Options: options.Index().SetName("chat_app_created")},
	})
	must(err)

	_, err = d.Admins().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "loginNorm", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_admin_login")},
		{Keys: bson.D{{Key: "adminId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uniq_adminId")},
	})
	must(err)
}
