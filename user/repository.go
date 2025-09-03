package user

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *User) error {
	result := r.db.Create(user)

	if result.Error != nil {
		fmt.Printf("create user error : %s\n", result.Error)
		return result.Error
	}
	return nil
}

func (r *Repository) GetById(id uint) (User, error) {
	var user User

	result := r.db.First(&user, id)
	if result.Error != nil {
		fmt.Printf("get user error : %s\n", result.Error)
		return User{}, result.Error
	}
	return user, nil
}

func (r *Repository) GetByEmail(email string) (User, error) {
	var user User

	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		fmt.Printf("get user error : %s\n", result.Error)
		return User{}, result.Error
	}
	return user, nil
}

func (r *Repository) GetBySocialIDAndProvider(socialID string, provider SocialProvider) (User, error) {
	var user User

	result := r.db.Where("social_id = ? AND social_provider = ?", socialID, provider).First(&user)
	if result.Error != nil {
		return User{}, result.Error //Error is check
	}
	return user, nil
}

func (r *Repository) GetGenres() ([]Genre, error) {
	var genres []Genre

	result := r.db.Find(&genres)

	if result.Error != nil {
		return nil, result.Error
	}
	return genres, nil
}

func (r *Repository) UpdateUserGenres(genresIDs []uint, userID uint) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&UserGenre{}).Error; err != nil {
		return err
	}

	var userGenres []UserGenre
	for _, genreID := range genresIDs {
		userGenres = append(userGenres, UserGenre{
			UserID:  userID,
			GenreID: genreID,
		})
	}

	if len(userGenres) > 0 {
		if err := r.db.Create(&userGenres).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) GetUserGenres(userID uint) ([]Genre, error) {
	var genres []Genre

	err := r.db.
		Joins("JOIN user_genres ug ON ug.genre_id = genres.id").
		Where("ug.user_id = ?", userID).
		Find(&genres).Error

	if err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *Repository) DeleteFavGenres(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&UserGenre{})
	return result.Error
}

func (r *Repository) UpdateUser(userID uint, nickname *string, photo *string, birthday *time.Time) (*User, error) {
	updates := map[string]interface{}{}

	if nickname != nil {
		updates["nickname"] = *nickname
	}
	if photo != nil {
		updates["photo"] = *photo
	}
	if birthday != nil {
		updates["birthday"] = *birthday
	}

	var updated User
	if err := r.db.Model(&User{}).
		Where("id = ?", userID).
		Updates(updates).
		First(&updated).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}
