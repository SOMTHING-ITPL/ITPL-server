package performance

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, rdb *redis.Client) *Repository {
	return &Repository{db: db, rdb: rdb}
}
