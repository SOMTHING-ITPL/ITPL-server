package email

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{rdb: rdb}
}
func (r *Repository) SetEmailCode(ctx context.Context, email string, code string, ttl time.Duration) error {
	return r.rdb.Set(ctx, email, code, ttl).Err()
}

func (r *Repository) GetEmailCode(ctx context.Context, email string) (string, error) {
	return r.rdb.Get(ctx, email).Result()
}

func (r *Repository) SetVerifiedEmail(ctx context.Context, email string) error {
	if err := r.rdb.Del(ctx, "email:code:"+email).Err(); err != nil {
		return fmt.Errorf("failed to delete email code: %w", err)
	}

	if err := r.rdb.Set(ctx, "email:verified:"+email, "true", 10*time.Minute).Err(); err != nil {
		return fmt.Errorf("failed to set verified email: %w", err)
	}

	return nil
}

func (r *Repository) CheckVerifiedEmail(ctx context.Context, email string) (bool, error) {
	result, err := r.rdb.Get(ctx, "email:verified:"+email).Result()

	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return result == "true", nil
}
