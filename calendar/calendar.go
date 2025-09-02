package calendar

import (
	"time"
)

func (r *Repository) CreateCalendar(perfID uint, userID uint, date time.Time) error {
	data := Calendar{
		PerformanceID: perfID,
		UserID:        userID,
		Year:          date.Year(),
		Month:         int(date.Month()),
		Day:           date.Day(),
	}

	if err := r.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetCalendar(userID uint, month int, year int) ([]Calendar, error) {
	var calendar []Calendar

	if err := r.db.Where("user_id = ? AND month = ? AND year = ?", userID, month, year).
		Preload("Performance").
		Find(&calendar).Error; err != nil {
		return nil, err
	}

	return calendar, nil
}

func (r *Repository) DeleteCalendar(id uint) error {
	result := r.db.Where("id = ?", id).Delete(&Calendar{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
