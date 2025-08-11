package user

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *User) error
	GetById(id uint) (User, error)
	GetByEmailAndProvider(email string, provider SocialProvider) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user *User) error {
	result := r.db.Create(user)

	if result.Error != nil {
		fmt.Printf("create user error : %s\n", result.Error)
		return result.Error
	}
	return nil
}

func (r *repository) GetById(id uint) (User, error) {
	var user User

	result := r.db.First(&user, id)
	if result.Error != nil {
		fmt.Printf("get user error : %s\n", result.Error)
		return User{}, result.Error
	}
	return user, nil
}

func (r *repository) GetByEmailAndProvider(email string, provider SocialProvider) (User, error) {
	var user User

	result := r.db.Where("email = ? AND social_provider = ?", email, provider).First(&user)
	if result.Error != nil {
		fmt.Printf("get user error : %s\n", result.Error)
		return User{}, result.Error
	}
	return user, nil
}
