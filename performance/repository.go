package performance

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, rdb *redis.Client) *Repository {
	return &Repository{db: db, rdb: rdb}
}

func (r *Repository) CreateFacility(facility *Facility) (uint, error) {
	result := r.db.Create(facility)

	if result.Error != nil {
		fmt.Printf("create facility error : %s\n", result.Error)
		return 0, result.Error
	}
	return facility.ID, nil

}

func (r *Repository) GetFacilityById(id uint) (*Facility, error) {
	var facility Facility

	result := r.db.First(&facility, id)
	if result.Error != nil {
		fmt.Printf("get facility error %d : %s\n ", id, result.Error)
		return &Facility{}, result.Error
	}
	return &facility, nil
}

func (r *Repository) GetFacilityByKopisID(id string) (*Facility, error) {
	var facility Facility

	result := r.db.Where("kopis_facility_key = ?", id).First(&facility)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &facility, nil
}

func (r *Repository) GetPerformanceByKopisID(id string) (*Performance, error) {
	var performance Performance

	result := r.db.Where("kopis_performance_key = ?", id).First(&performance)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &performance, nil
}

func (r *Repository) CreatePerformance(performance *Performance) (uint, error) {
	result := r.db.Create(performance)

	if result.Error != nil {
		fmt.Printf("create performance error : %s\n", result.Error)
		return 0, result.Error
	}

	return performance.ID, nil
}

func (r *Repository) GetPerformanceById(id uint) (*Performance, error) {
	performance := &Performance{}
	result := r.db.First(performance, id)

	if result.Error != nil {
		fmt.Printf("get performance error : %s\n", result.Error)
		return performance, result.Error
	}

	return performance, nil
}

func (r *Repository) GetPerformancesByIDs(ids []uint) ([]Performance, error) {
	var performances []Performance
	err := r.db.Where("id IN ?", ids).Find(&performances).Error
	if err != nil {
		return nil, err
	}
	return performances, nil
}

// TODO : 이거 제대로 동작하는지 테스트해봐야 함.
func (r *Repository) FindPerformances(page, limit, genre int, region, keyword string) ([]Performance, error) {
	var performances []Performance
	db := r.db.Model(&Performance{})

	//this should be code 01 02 03 04 ...
	if genre != 0 {
		db = db.Where("genre = ?", genre)
	}
	//
	if region != "" {
		db = db.Where("region = ?", region)
	}
	if keyword != "" { //이거 키워드에서 찾게 해야 하나?
		db = db.Where("title LIKE ?", "%"+keyword+"%")
		keyword = strings.TrimSpace(keyword)

	}

	offset := (page - 1) * limit
	if err := db.
		Where("status IN (?)", []string{"공연중", "공연예정"}).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&performances).Error; err != nil {
		return nil, err
	}

	return performances, nil
}

func (r *Repository) FindFacilities(page, limit int, region string) ([]Facility, error) {
	var facilities []Facility
	db := r.db.Model(&Facility{})

	if region != "" {
		db = db.Where("region = ?", region)
	}

	offset := (page - 1) * limit
	if err := db.
		Order("created_at DESC"). // 최신순으로
		Limit(limit).
		Offset(offset).
		Find(&facilities).Error; err != nil {
		return nil, err
	}

	return facilities, nil
}

func (r *Repository) CreatePerformanceTicketSite(site *PerformanceTicketSite) error {
	result := r.db.Create(site)

	if result.Error != nil {
		fmt.Printf("create ticket site error : %s\n", result.Error)
		return result.Error
	}

	return nil
}

func (r *Repository) CreatePerformanceImage(site *PerformanceImage) error {
	result := r.db.Create(site)

	if result.Error != nil {
		fmt.Printf("create ticket site error : %s\n", result.Error)
		return result.Error
	}

	return nil
}

func (r *Repository) GetPerformanceImages(prefId uint) ([]PerformanceImage, error) {
	var images []PerformanceImage

	if err := r.db.
		Where("performance_id = ?", prefId).
		Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (r *Repository) GetPerformanceWithTicketsAndImages(perfID uint) (*PerformanceWithTicketsAndImage, error) {
	var perf Performance

	err := r.db.Preload("TicketSites").Preload("Images").First(&perf, perfID).Error
	if err != nil {
		return nil, err
	}

	log.Printf("Fetched Performance: ID=%d, FacilityID=%d, Title=%s\n", perf.ID, perf.FacilityID, perf.Title)

	result := &PerformanceWithTicketsAndImage{
		Performance:       perf,
		TicketSites:       perf.TicketSites,
		PerformanceImages: perf.Images,
	}

	return result, nil
}
func (r *Repository) GetUserLike(userID uint) ([]Performance, error) {
	var favorites []PerformanceUserLike

	if err := r.db.
		Preload("Performance").
		Where("user_id = ?", userID).
		Find(&favorites).Error; err != nil {
		return nil, err
	}

	performances := make([]Performance, len(favorites))
	for i, f := range favorites {
		performances[i] = f.Performance
	}

	return performances, nil
}

func (r *Repository) CreateUserLike(perfID uint, userID uint) error {
	var existing PerformanceUserLike

	err := r.db.
		Where("performance_id = ? AND user_id = ?", perfID, userID).
		First(&existing).Error

	if err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	fav := PerformanceUserLike{
		PerformanceID: perfID,
		UserID:        userID,
	}

	if err := r.db.Create(&fav).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteUserLike(perfID uint, userID uint) error {
	result := r.db.
		Where("performance_id = ? AND user_id = ?", perfID, userID).
		Delete(&PerformanceUserLike{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// gin context 그대로 사용
func (r *Repository) CreateRecentView(userID uint, performanceID uint, ctx context.Context) error {
	key := fmt.Sprintf("user:%d:recent_views", userID)

	if err := r.rdb.LPush(ctx, key, performanceID).Err(); err != nil {
		return err
	}

	//user 당 최대 10개까지만 저장함.
	if err := r.rdb.LTrim(ctx, key, 0, 9).Err(); err != nil {
		return err
	}

	// 공연 별 조회수 집계
	if err := r.rdb.ZIncrBy(ctx, "performance_views", 1, fmt.Sprintf("%d", performanceID)).Err(); err != nil {
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
	zs, err := r.rdb.ZRevRangeWithScores(ctx, "performance_views", 0, topN-1).Result()
	if err != nil {
		return nil, err
	}

	topPerformances := make([]PerformanceScore, len(zs))
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
	ttl := 3 * 24 * 60 * 60

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
