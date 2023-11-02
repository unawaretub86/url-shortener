package mongodb

import (
	"context"
	"time"

	errs "github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/unawaretub86/url-shortener/shortener"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoCLient(mongoURL string, mongoTimeOut int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeOut)*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeOut int) (shortener.RedirectRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeOut) * time.Second,
		database: mongoDB,
	}
	client, err := newMongoCLient(mongoURL, mongoTimeOut)
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewMongoRepo.mongodb")
	}

	repo.client = client

	return repo, nil
}

// Find implements shortener.RedirectRepository.
func (r *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)

	defer cancel()

	redirect := &shortener.Redirect{}

	collection := r.client.Database(r.database).Collection("redirects")

	filter := bson.M{"code": code}

	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
		}

		return nil, errs.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}

	return redirect, nil
}

// Store implements shortener.RedirectRepository.
func (r *mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)

	defer cancel()

	collection := r.client.Database(r.database).Collection("redirects")
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":       redirect.Code,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		},
	)
	if err != nil {
		return errs.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Store")
	}

	return nil
}
