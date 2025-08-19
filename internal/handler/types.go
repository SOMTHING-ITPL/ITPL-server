package handler

import (
	"github.com/SOMTHING-ITPL/ITPL-server/internal/email"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
)

type UserHandler struct {
	userRepository *user.Repository
	smtpRepository *email.Repository
}
