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
