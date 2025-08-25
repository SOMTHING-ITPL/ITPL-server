package performance

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
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
		fmt.Printf("get facility error : %s\n", result.Error)
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

// func (r *Repository) IsExistFacility(id uint) (bool, error) {
//     var facility Facility

//     result := r.db.First(&facility, id)
//     if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//         return false, nil ] //없을 때
//     }
//     if result.Error != nil { //DB 오류
//         return false, result.Error
//     }
//     return true, nil
// }

func (r *Repository) CreatePerformance(performance *Performance) (uint, error) {
	result := r.db.Create(performance)

	if result.Error != nil {
		fmt.Printf("create performance error : %s\n", result.Error)
		return 0, result.Error
	}

	return performance.ID, nil
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
