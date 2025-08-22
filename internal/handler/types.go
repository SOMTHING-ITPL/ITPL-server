package handler

import (
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

type UserHandler struct {
	userRepository *user.Repository
}

type PlaceHandler struct {
	database       *gorm.DB
	userRepository *user.Repository
}

type CourseHandler struct {
	database       *gorm.DB
	userRepository *user.Repository
}
