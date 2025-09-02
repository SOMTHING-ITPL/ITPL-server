package artist

import "gorm.io/gorm"

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetArtist() ([]Artist, error) {
	var artists []Artist

	result := r.db.Find(&artists)

	if result.Error != nil {
		return nil, result.Error
	}
	return artists, nil
}

func (r *Repository) UpdateUserArtist(artistIDs []uint, userID uint) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&UserArtist{}).Error; err != nil {
		return err
	}
	var userArtists []UserArtist
	for _, artistID := range artistIDs {
		userArtists = append(userArtists, UserArtist{
			UserID:   userID,
			ArtistID: artistID,
		})
	}

	if len(userArtists) > 0 {
		if err := r.db.Create(&userArtists).Error; err != nil {
			return err
		}
	}

	return nil
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

func (r *Repository) DeleteFavArtists(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&UserArtist{})
	return result.Error
}
