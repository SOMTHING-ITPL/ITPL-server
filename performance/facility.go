package performance

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

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
