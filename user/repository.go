package user

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

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

func (r *Repository) GetByUserName(id string) (User, error) {
	var user User

	result := r.db.Where("user_name = ?", id).First(&user)
	if result.Error != nil {
		fmt.Printf("get user error : %s\n", result.Error)
		return User{}, result.Error
	}
	return user, nil
}

func (r *Repository) GetByEmailAndProvider(email string, provider SocialProvider) (User, error) {
	var user User

	result := r.db.Where("email = ? AND social_provider = ?", email, provider).First(&user)
	if result.Error != nil {
		fmt.Printf("get user error : %s\n", result.Error)
		return User{}, result.Error
	}
	return user, nil
}

func (r *Repository) GetArtist() ([]Artist, error) {
	var artists []Artist

	result := r.db.Find(&artists)

	if result.Error != nil {
		return nil, result.Error
	}
	return artists, nil
}

func (r *Repository) GetGenres() ([]Genre, error) {
	var genres []Genre

	result := r.db.Find(&genres)

	if result.Error != nil {
		return nil, result.Error
	}
	return genres, nil
}

func (r *Repository) SetUserArtist(artistIDs []uint, userID uint) error {
	var userArtists []UserArtist
	for _, artistID := range artistIDs {
		userArtists = append(userArtists, UserArtist{
			UserID:   userID,
			ArtistID: artistID,
		})
	}
	return r.db.Create(&userArtists).Error
}

func (r *Repository) SetUserGenres(genresIDs []uint, userID uint) error {
	var userGenres []UserGenre
	for _, artistID := range genresIDs {
		userGenres = append(userGenres, UserGenre{
			UserID:  userID,
			GenreID: artistID,
		})
	}
	return r.db.Create(&userGenres).Error
}

func (r *Repository) GetUserArtists(userID uint) ([]Artist, error) {
	var artists []Artist

	// SELECT a.* FROM artists a
	// INNER JOIN user_artists ua ON a.id = ua.artist_id
	// WHERE ua.user_id = ?
	err := r.db.
		Joins("JOIN user_artists ua ON ua.artist_id = artists.id").
		Where("ua.user_id = ?", userID).
		Find(&artists).Error

	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (r *Repository) GetUserGenres(userID uint) ([]Genre, error) {
	var genres []Genre

	err := r.db.
		Joins("JOIN user_genres ug ON ug.genre_id = genre.id").
		Where("ug.user_id = ?", userID).
		Find(&genres).Error

	if err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *Repository) DeleteFavArtists(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&UserArtist{})
	return result.Error
}

func (r *Repository) DeleteFavGenres(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&UserGenre{})
	return result.Error
}
