package performance

import (
	"errors"

	"gorm.io/gorm"
)

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
