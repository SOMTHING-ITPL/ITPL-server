package performance

import (
	"fmt"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateFacility(facility *Facility) error {
	result := r.db.Create(facility)

	if result.Error != nil {
		fmt.Printf("create facility error : %s\n", result.Error)
		return result.Error
	}
	return nil

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

func (r *Repository) CreatePerformance(performance *Performance) error {
	result := r.db.Create(performance)

	if result.Error != nil {
		fmt.Printf("create performance error : %s\n", result.Error)
		return result.Error
	}

	return nil
}
