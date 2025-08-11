package handler

import (
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userRepo user.Repository, userService user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := h.userService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not found user! "})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) RegisterLocalUser() gin.HandlerFunc {
	//when you register specific
	type req struct {
		NickName string `json:"nickName"`
		UserName string `json:"userName"`
		Pwd      string `json:"pwd"`
		Email    string `json:"email"`
	}

	return func(c *gin.Context) {
		var request req

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := user.User{
			Username:       request.UserName,
			NickName:       &request.NickName,
			Email:          &request.Email,
			SocialProvider: user.ProviderLocal,
			EncryptPwd:     &request.Pwd,
		}

		err := h.userService.CreateUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func (h *UserHandler) RegisterSocialUser() func(c *gin.Context) {

	// return func(c *gin.Context) {
	// 	NickName string `json:"nickName"`
	// 	UserName string `json:"userName"`
	// 	Pwd      string `json:"pwd"`
	// 	Email    string `json:"email"`

	// }
}
