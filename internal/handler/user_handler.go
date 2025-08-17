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

func (h *UserHandler) CheckValidId() gin.HandlerFunc {
	type req struct {
		UserID string `json:"user_id"`
	}

	return func(c *gin.Context) {
		var request req

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "check your body"})
			return
		}

		_, err := h.userRepository.GetByUserID(request.UserID)
		c.JSON(http.StatusOK, gin.H{"valid": err == nil}) //true or false
	}

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
		UserID   string `json:"user_id"`
		Pwd      string `json:"password"`
		Email    string `json:"email"`
	}
	type res struct {
		Token string `json:"token"`
	}
	return func(c *gin.Context) {
		var request req

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, err := h.userRepository.GetByUserID(request.UserID); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID is already exist check if is USERID is validate"})
			return
		}

		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(request.Pwd), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}

		hashedPwdStr := string(hashedPwd)
		user := user.User{
			UserID:         request.UserID,
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

		c.JSON(http.StatusCreated, res{Token: jwt})
	}
}

func (h *UserHandler) LoginSocialUser() gin.HandlerFunc {
	type req struct {
		SocialProvider string `json:"social_provider"`
		AccessToken    string `json:"access_token"`
	}
	type res struct {
		Token string `json:"token"`
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

			c.JSON(http.StatusCreated, res{Token: jwt})
			return
		}

		jwt, err := auth.GenerateJWT(&targetUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate jwt"})
			return
		}

		c.JSON(http.StatusCreated, res{Token: jwt})
	}
}

func (h *UserHandler) LoginLocalUser() gin.HandlerFunc {
	type req struct {
		UserID string `json:"user_id"`
		Pwd    string `json:"password"`
	}
	type res struct {
		Token string `json:"token"`
	}

	return func(c *gin.Context) {
		var request req

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "check your body"})
			return
		}

		user, err := h.userRepository.GetByUserID(request.UserID)
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

		c.JSON(http.StatusOK, res{Token: jwt})
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

func (h *UserHandler) GetGenres() gin.HandlerFunc {
	return func(c *gin.Context) {
		genres, err := h.userRepository.GetGenres()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get genres"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": genres})
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

func (h *UserHandler) AddUserGenre() gin.HandlerFunc {
	type req struct {
		GenreIDs []uint `json:"genre_ids"`
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

		if err := h.userRepository.SetUserArtist(request.GenreIDs, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fail to set user genre on db"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User genre updated successfully"})
	}
}
