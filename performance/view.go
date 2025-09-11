package performance

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (r *Repository) CreateRecentView(userID uint, performanceID uint, ctx context.Context) error {
	key := fmt.Sprintf("user:%d:recent_views", userID)

	if err := r.rdb.LPush(ctx, key, performanceID).Err(); err != nil {
		return err
	}

	//user 당 최대 10개까지만 저장함.
	if err := r.rdb.LTrim(ctx, key, 0, 9).Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetRecentViews(userID uint, ctx context.Context) ([]uint, error) {
	key := fmt.Sprintf("user:%d:recent_views", userID)

	ids, err := r.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	uintIDs := make([]uint, 0, len(ids))
	for _, s := range ids {
		idUint, convErr := strconv.ParseUint(s, 10, 64)
		if convErr != nil {
			fmt.Printf("invalid performance id in redis: %s\n", s)
			continue
		}
		uintIDs = append(uintIDs, uint(idUint))
	}

	return uintIDs, nil
}

// 현재 가장 조회가 많은 공연 조회하는 부분
func (r *Repository) GetTopPerformances(topN int64, ctx context.Context) ([]PerformanceScore, error) {
	zs, err := r.rdb.ZRevRangeWithScores(ctx, "performance_views", 0, topN-1).Result() // 0~topN-1 개만 가져옴.
	if err != nil {
		return nil, err
	}

	topPerformances := make([]PerformanceScore, 0, len(zs))
	for _, z := range zs {
		idStr, ok := z.Member.(string)
		if !ok {
			continue
		}
		idUint, convErr := strconv.ParseUint(idStr, 10, 64)
		if convErr != nil {
			continue
		}

		topPerformances = append(topPerformances, PerformanceScore{
			ID:    uint(idUint),
			Score: z.Score,
		})
	}

	return topPerformances, nil
}

func (r *Repository) IncrementPerformanceScore(perfID uint, score float64, ctx context.Context) error {
	key := "performance_views"
	ttl := 3 * 24 * 60 * 60 //3일

	//점수 증가 시키는 부분
	if err := r.rdb.ZIncrBy(ctx, key, score, fmt.Sprintf("%d", perfID)).Err(); err != nil {
		log.Printf("Failed to increment score for performance %d: %v", perfID, err)
		return err
	}

	//TTL 최대 3일?
	err := r.rdb.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		log.Printf("Failed to set TTL: %v", err)
	}

	return nil
}
