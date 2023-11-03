package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	errs "github.com/pkg/errors"

	"github.com/unawaretub86/url-shortener/shortener"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisUrl string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	_, err = client.Ping().Result()

	return client, nil
}

func NewRedisRepository(redisURL string) (shortener.RedirectRepository, error) {
	repo := &redisRepository{}

	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewRedisRepo.redis")
	}

	repo.client = client

	return repo, nil
}

// Find implements shortener.RedirectRepository.
func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}

	key := r.generateKey(code)

	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewRedisRepo.Find")
	}

	if len(data) == 0 {
		return nil, errs.Wrap(shortener.ErrRedirectNotFound, "repository.NewRedisRepo.Find")
	}

	//converting timestamp from string to int , based 10 as in 64 bit
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewRedisRepo.Find")
	}

	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt

	return redirect, nil
}

// Store implements shortener.RedirectRepository.
func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)

	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errs.Wrap(err, "repository.NewRedisRepo.Store")
	}

	return nil
}

// Redis is just a key value store inside of memory
//
// This is a utility func which allows us to crate the key that we're to use to get the data from DB
func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}
