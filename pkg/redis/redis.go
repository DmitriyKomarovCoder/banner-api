package redis

import (
	"context"

	"github.com/DmitriyKomarovCoder/banner-api/pkg/logger"
	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	Addr   string
	DB     int
	Client *redis.Client
	Log    logger.Logger
	ctx    context.Context
}

func NewRedisRepository(addr string, db int, log logger.Logger) *RedisRepository {
	return &RedisRepository{
		Addr: addr,
		DB:   db,
		Log:  log,
	}
}

func (r *RedisRepository) Connect() error {
	r.Client = redis.NewClient(&redis.Options{
		Addr: r.Addr,
		DB:   r.DB,
	})
	_, err := r.Client.Ping(context.Background()).Result()
	return err
}

func (r *RedisRepository) Close(ctx context.Context) error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}
