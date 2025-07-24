package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLDB(name string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewRedisDB(addrs string, pass string) (*redis.Client, error) {
	ctx := context.Background()

	rclient := redis.NewClient(
		&redis.Options{
			Addr:     addrs,
			Password: pass,
			DB:       0,
			Protocol: 3,
		})

	_, err := rclient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rclient, nil
}
