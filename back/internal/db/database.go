package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func Connect(ctx context.Context, uri, dbName string) *Database {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping: %v", err)
	}
	return &Database{Client: client, DB: client.Database(dbName)}
}

func (d *Database) Close(ctx context.Context) error { return d.Client.Disconnect(ctx) }

func (d *Database) Users() *mongo.Collection         { return d.DB.Collection("users") }
func (d *Database) Sessions() *mongo.Collection      { return d.DB.Collection("sessions") }
func (d *Database) Profiles() *mongo.Collection      { return d.DB.Collection("profiles") }
func (d *Database) Vacancies() *mongo.Collection     { return d.DB.Collection("vacancies") }
func (d *Database) Resumes() *mongo.Collection       { return d.DB.Collection("resumes") }
func (d *Database) Applications() *mongo.Collection  { return d.DB.Collection("applications") }
func (d *Database) ChatMessages() *mongo.Collection  { return d.DB.Collection("chat_messages") }
func (d *Database) Admins() *mongo.Collection        { return d.DB.Collection("admins") }
func (d *Database) Subscriptions() *mongo.Collection { return d.DB.Collection("subscriptions") }
