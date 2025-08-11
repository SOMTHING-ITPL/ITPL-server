package handler

import (
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/auth"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepository *user.Repository
}

func NewUserHandler(userRepository *user.Repository) *UserHandler {
	return &UserHandler{userRepository: userRepository}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userRepository.GetById(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not found user! "})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) RegisterLocalUser() gin.HandlerFunc {
	type req struct {
		NickName string `json:"nick_name"`
		UserName string `json:"user_name"`
		Pwd      string `json:"password"`
		Email    string `json:"email"`
	}
	type res struct {
		token string
	}

	return func(c *gin.Context) {
		var request req

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(request.Pwd), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}

		hashedPwdStr := string(hashedPwd)
		user := user.User{
			Username:       request.UserName,
			NickName:       &request.NickName,
			Email:          &request.Email,
			SocialProvider: user.ProviderLocal,
			EncryptPwd:     &hashedPwdStr,
		}

		err = h.userRepository.CreateUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		jwt, err := auth.GenerateJWT(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate jwt"})
			return
		}

		c.JSON(http.StatusCreated, res{token: jwt})
	}
}

func (h *UserHandler) RegisterSocialUser() gin.HandlerFunc {
	type req struct {
		SocialProvider string `json:"social_provider"`
		AccessToken    string `json:"user_name"`
	}
	type res struct {
		token string
	}

	return func(c *gin.Context) {
		var request req

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client, err := user.GetProviderClient(request.SocialProvider)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		nickName := user.GenerateNanoIDNickname()
		result, err := client.Login(request.AccessToken)

		user := user.User{
			NickName:       &nickName,
			Email:          result.Email,
			SocialProvider: user.SocialProvider(request.SocialProvider),
		}

		jwt, err := auth.GenerateJWT(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate jwt"})
			return
		}

		c.JSON(http.StatusCreated, res{token: jwt})
	}
}

func (h *UserHandler) LoginLocalUser() gin.HandlerFunc {
	type req struct {
		UserName string `json:"user_name"`
		Pwd      string `json:"password"`
	}
	type res struct {
		token string
	}

	return func(c *gin.Context) {
		var request req

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := h.userRepository.GetByUserName(request.UserName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get user"})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(*user.EncryptPwd), []byte(request.Pwd))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		jwt, err := auth.GenerateJWT(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate jwt"})
			return
		}

		c.JSON(http.StatusOK, res{token: jwt})
	}
}
