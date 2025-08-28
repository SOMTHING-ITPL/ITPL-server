package performance

import (
	"errors"
	"fmt"
	"log"
	"strings"

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
