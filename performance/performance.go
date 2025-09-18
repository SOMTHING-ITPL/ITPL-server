package performance

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

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

// use This function when server recommend performance to user.
func (r *Repository) GetPerformancesByIDs(ids []uint) ([]Performance, error) {
	var performances []Performance
	err := r.db.Where("id IN ?", ids).Find(&performances).Error
	if err != nil {
		return nil, err
	}
	return performances, nil
}

// TODO : 이거 제대로 동작하는지 테스트해봐야 함.
func (r *Repository) FindPerformances(page, limit, genre int, region, keyword string) ([]Performance, int64, error) {
	var performances []Performance
	var total int64

	db := r.db.Model(&Performance{})

	//this should be code 01 02 03 04 ...
	if genre != 0 && genre != 9 {
		db = db.Where("genre = ?", genre)
	}
	if genre == 9 { //내한 여부
		db = db.Where("is_foreign = ?", "Y")
	}
	//서울 .. 득
	if region != "" {
		db = db.Where("region LIKE ?", "%"+region+"%")
	}
	if keyword != "" {
		keyword = strings.TrimSpace(keyword)
		likePattern := "%" + keyword + "%"

		db = db.Where(
			"title LIKE ? OR `cast` LIKE ? OR `keyword` LIKE ?",
			likePattern, likePattern, likePattern,
		)
	}

	db = db.Where("status IN (?)", []string{"공연중", "공연예정"})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := db.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&performances).Error; err != nil {
		return nil, 0, err
	}

	return performances, total, nil
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

func (r *Repository) UpdatePerformance(perf *Performance) error {
	return r.db.Save(perf).Error
}

func (r *Repository) GetRecentPerformance(targetDate time.Time, num int) ([]Performance, error) {
	var perfs []Performance

	//start Date 바로 직전인 것만 가져옴
	err := r.db.Where("start_date >= ?", targetDate).Order("start_date ASC").Limit(num).Find(&perfs).Error

	if err != nil {
		log.Printf("Fail to Get start_date : %s", err)
		return nil, err
	}

	return perfs, nil
}
