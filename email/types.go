package email

import "github.com/go-redis/redis/v8"

type Repository struct {
	rdb *redis.Client
}
