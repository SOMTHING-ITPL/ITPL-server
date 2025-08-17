package handler

import (
	"net/http"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/auth"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func NewUserHandler(userRepository *user.Repository) *UserHandler {
	return &UserHandler{userRepository: userRepository}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	user, err := h.userRepository.GetById(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not found user! "})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Profile iamge?
// func (h *UserHandler) EditProfile() gin.HandlerFunc {
// 	type req struct {
// 		NickName string `json:"nick_name"`
// 	}
// }

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

func (h *UserHandler) LoginSocialUser() gin.HandlerFunc {
	type req struct {
		SocialProvider string `json:"social_provider"`
		AccessToken    string `json:"access_token"`
	}
	type res struct {
		token string
	}

	return func(c *gin.Context) {
		var request req

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "check request params"})
			return
		}

		client, err := user.GetProviderClient(request.SocialProvider)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to getClient"})
			return
		}

		result, err := client.Login(request.AccessToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to login"})
			return
		}

		targetUser, err := h.userRepository.GetByEmailAndProvider(*result.Email, user.SocialProvider(request.SocialProvider))
		if err != nil {
			nickName := user.GenerateNanoIDNickname()
			targetUser := user.User{
				NickName:       &nickName,
				Email:          result.Email,
				SocialProvider: user.SocialProvider(request.SocialProvider),
			}

			err = h.userRepository.CreateUser(&targetUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}

			jwt, err := auth.GenerateJWT(&targetUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate jwt"})
				return
			}

			c.JSON(http.StatusCreated, res{token: jwt})
			return
		}

		jwt, err := auth.GenerateJWT(&targetUser)
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

func (h *UserHandler) GetArtists() gin.HandlerFunc {
	return func(c *gin.Context) {
		artist, err := h.userRepository.GetArtist()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get artist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": artist})
	}
}

func (h *UserHandler) GetGeners() gin.HandlerFunc {
	return func(c *gin.Context) {
		geners, err := h.userRepository.GetGeners()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get geners"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": geners})
	}
}

func (h *UserHandler) AddUserArtist() gin.HandlerFunc {
	type req struct {
		artistIDs []uint `json:"artist_ids"`
	}
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		userID, ok := userIDInterface.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			return
		}

		var request req
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "check request params"})
			return
		}

		if err := h.userRepository.SetUserArtist(request.artistIDs, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fail to set user artist on db"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User artists updated successfully"})
	}
}

func (h *UserHandler) AddUserGener() gin.HandlerFunc {
	type req struct {
		GenerIDs []uint `json:"gener_ids"`
	}
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		userID, ok := userIDInterface.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			return
		}

		var request req
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "check request params"})
			return
		}

		if err := h.userRepository.SetUserArtist(request.GenerIDs, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fail to set user gener on db"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User gener updated successfully"})
	}
}
