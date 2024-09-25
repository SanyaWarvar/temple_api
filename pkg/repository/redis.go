package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	db      *redis.Client
	codeExp time.Duration
}

func NewRedisDb(opt *redis.Options) (*redis.Client, error) {
	ctx := context.Background()
	client := redis.NewClient(opt)
	err := client.Ping(ctx).Err()
	return client, err
}

func NewCache(db *redis.Client, codeExp time.Duration) *Cache {
	return &Cache{
		db:      db,
		codeExp: codeExp,
	}
}

func (c *Cache) SaveConfirmCode(email, code string) error {
	ctx := context.Background()
	err := c.db.Set(ctx, email, code, c.codeExp).Err()
	return err
}

func (c *Cache) GetConfirmCode(email string) (string, time.Duration, error) {
	var ttl time.Duration
	ctx := context.Background()
	code, err := c.db.Get(ctx, email).Result()
	if err != nil {
		return "", ttl, err
	}
	ttl, err = c.db.TTL(ctx, email).Result()

	return code, ttl, err
}
