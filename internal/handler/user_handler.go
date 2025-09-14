package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/email"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/auth"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewUserHandler(userRepository *user.Repository, smtpRepository *email.Repository, awsBucket *aws.BucketBasics) *UserHandler {
	return &UserHandler{userRepository: userRepository, smtpRepository: smtpRepository, BucketBasics: awsBucket}
}

func (h *UserHandler) SendEmailCode() gin.HandlerFunc {
	type req struct {
		Email string `json:"email" binding:"required,email"`
	}

	return func(c *gin.Context) {
		var request req

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, err := h.userRepository.GetByEmail(request.Email); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This email is already exist"})
			return
		}

		code := email.GenerateCode(6)

		//save code during 10min
		if err := h.smtpRepository.SetEmailCode(c, request.Email, code, 10*time.Minute); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save code"})
			return
		}

		//send mail
		if err := email.SendMail(request.Email, code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send code"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Message": "verification code sent"})
	}

}

func (h *UserHandler) VerifyEmailCode() gin.HandlerFunc {
	type req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	return func(c *gin.Context) {
		var request req
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		savedCode, err := h.smtpRepository.GetEmailCode(c, request.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "code expired or not found"})
			return
		}

		if savedCode != request.Code {
			c.JSON(http.StatusOK, CommonRes{
				Message: "email is not Verified check Code ",
				Data:    false,
			})
			return
		}

		if h.smtpRepository.SetVerifiedEmail(c, request.Email) != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to set verified Email "})
			return
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "email is Verified",
			Data:    true,
		})
	}
}

func (h *UserHandler) GetUser() gin.HandlerFunc {
	type res struct {
		CreatedAt      string  `json:"created_at"`
		UpdatedAt      string  `json:"updated_at"`
		Email          string  `json:"email ,omitempty"`
		NickName       string  `json:"nick_name"`
		SocialProvider string  `json:"social_provider"`
		Birthday       *string `json:"birthday ,omitempty"`
		Photo          *string `json:"profile_url ,omitempty"`
	}
	return func(c *gin.Context) {

		userID, _ := c.Get("userID")

		user, err := h.userRepository.GetById(userID.(uint))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not found user! "})
			return
		}
		var birthday *string

		if user.Birthday != nil {
			formatted := user.Birthday.Format("20060102")
			birthday = &formatted
		}

		var photoURL *string
		if user.Photo != nil {
			url, err := aws.GetPresignURL(h.BucketBasics.AwsConfig, h.BucketBasics.BucketName, *user.Photo)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get photo in aws: " + err.Error()})
				return
			}
			photoURL = &url
		}

		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(c.Writer)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(CommonRes{
			Message: "success",
			Data: CommonRes{
				Message: "success",
				Data: res{
					CreatedAt:      user.CreatedAt.Format(time.RFC3339),
					UpdatedAt:      user.UpdatedAt.Format(time.RFC3339),
					NickName:       user.NickName,
					Email:          *user.Email,
					SocialProvider: string(user.SocialProvider),
					Birthday:       birthday,
					Photo:          photoURL,
				},
			},
		})

	}
}

func (h *UserHandler) UpdateProfile() gin.HandlerFunc {
	type res struct {
		CreatedAt      string  `json:"created_at "`
		UpdatedAt      string  `json:"updated_at"`
		Email          string  `json:"email ,omitempty"`
		NickName       string  `json:"nick_name"`
		SocialProvider string  `json:"social_provider"`
		Birthday       *string `json:"birthday ,omitempty"`
		Photo          *string `json:"profile_url ,omitempty"`
	}

	return func(c *gin.Context) {
		userIDVal, _ := c.Get("userID")
		userID := userIDVal.(uint)
		user, err := h.userRepository.GetById(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
			return
		}
		nickName := c.PostForm("nickname")

		var birthdayTime *time.Time //그냥 stirng 으로 저장하는게 정신건강에 좋을 것 같다는 생각이
		birthdayStr := c.PostForm("birthday")
		if birthdayStr != "" {
			t, err := time.Parse("20060102", birthdayStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "birthday must be in yyyymmdd format"})
				return
			}
			birthdayTime = &t
		}

		var imageURL *string
		file, err := c.FormFile("profile")
		if err == nil {
			if user.Photo != nil {
				err = aws.DeleteImage(h.BucketBasics.S3Client, h.BucketBasics.BucketName, *user.Photo)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete old image"})
					return
				}
			}

			url, err := aws.UploadToS3(h.BucketBasics.S3Client, h.BucketBasics.BucketName, fmt.Sprintf("profile/%d", userID), file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("failed to upload profile image: %v", err),
				})
				return
			}
			imageURL = &url
		}

		updatedUser, err := h.userRepository.UpdateUser(userID, &nickName, imageURL, birthdayTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		var birthdayFormat *string
		if updatedUser.Birthday != nil {
			formatted := updatedUser.Birthday.Format("20060102")
			birthdayFormat = &formatted
		}

		var photoURL *string
		if updatedUser.Photo != nil {
			url, err := aws.GetPresignURL(h.BucketBasics.AwsConfig, h.BucketBasics.BucketName, *updatedUser.Photo)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get photo in aws: " + err.Error()})
				return
			}
			photoURL = &url
		}
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(c.Writer)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(CommonRes{
			Message: "success",
			Data: CommonRes{
				Message: "success",
				Data: res{
					CreatedAt:      updatedUser.CreatedAt.Format(time.RFC3339),
					UpdatedAt:      updatedUser.UpdatedAt.Format(time.RFC3339),
					NickName:       updatedUser.NickName,
					Email:          *updatedUser.Email,
					SocialProvider: string(updatedUser.SocialProvider),
					Birthday:       birthdayFormat,
					Photo:          photoURL,
				},
			},
		})

	}
}

func (h *UserHandler) RegisterLocalUser() gin.HandlerFunc {
	type req struct {
		NickName string `json:"nick_name" binding:"required"`
		Pwd      string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Birthday string `json:"birthday"` //required 아님
		// Agree1 bool `json:"agree_personal_info"`    // 필수 동의
		// Agree2 bool `json:"agree_marketing"`        // 마케팅 동의
		// Agree3 bool `json:"agree_terms_of_service"` // 서비스 약관 동의
		// Agree4 bool `json:"agree_privacy_policy"`   // 개인정보 처리방침 동의
		Agree1 bool `json:"agree_1"` // 필수 동의
		Agree2 bool `json:"agree_2"` // 마케팅 동의
		Agree3 bool `json:"agree_3"` // 서비스 약관 동의
		Agree4 bool `json:"agree_4"` // 개인정보 처리방침 동의

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
		//format check
		birthdayTime, err := time.Parse("20060102", request.Birthday)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "birthday must be in yyyymmdd format"})
			return
		}

		verified, err := h.smtpRepository.CheckVerifiedEmail(c, request.Email)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to get data through redis db"})
			return
		}
		if !verified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "verifing your email is first!"})
			return
		}

		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(request.Pwd), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}

		hashedPwdStr := string(hashedPwd)
		user := user.User{
			NickName:       request.NickName,
			Email:          &request.Email,
			SocialProvider: user.ProviderLocal,
			EncryptPwd:     &hashedPwdStr,
			Birthday:       &birthdayTime,
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

		c.JSON(http.StatusCreated, CommonRes{
			Message: "success",
			Data:    res{Token: jwt},
		})
	}
}

func (h *UserHandler) LoginSocialUser() gin.HandlerFunc { //Access 이런 거 다 구조화 해야 하는 건가?
	type req struct {
		SocialProvider string `json:"social_provider" binding:"required" `
		AccessToken    string `json:"access_token" binding:"required"`
	}
	type res struct {
		Token string `json:"token"`
		IsNew bool   `json:"is_new"`
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

		targetUser, err := h.userRepository.GetBySocialIDAndProvider(result.ID, user.SocialProvider(request.SocialProvider))
		if errors.Is(err, gorm.ErrRecordNotFound) { //Not found 일 경우,
			targetUser, err = h.RegisterSocialUser(user.SocialProvider(request.SocialProvider), result.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to create new user"})
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		jwt, err := auth.GenerateJWT(&targetUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate jwt"})
			return
		}

		c.JSON(http.StatusCreated, CommonRes{
			Message: "success",
			Data:    res{Token: jwt},
		})
	}
}

func (h *UserHandler) RegisterSocialUser(provider user.SocialProvider, socialID string) (user.User, error) {
	nickName := user.GenerateNanoIDNickname()
	targetUser := user.User{
		NickName:       nickName,
		SocialID:       &socialID,
		SocialProvider: provider,
	}

	err := h.userRepository.CreateUser(&targetUser)
	if err != nil {
		return user.User{}, err
	}
	return targetUser, nil
}

func (h *UserHandler) LoginLocalUser() gin.HandlerFunc {
	type req struct {
		Email string `json:"email" binding:"required"`
		Pwd   string `json:"password" binding:"required"`
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

		user, err := h.userRepository.GetByEmail(request.Email)
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

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    res{Token: jwt},
		})
	}
}

func (h *UserHandler) GetGenres() gin.HandlerFunc {
	return func(c *gin.Context) {
		genres, err := h.userRepository.GetGenres()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get genres"})
			return
		}

		res := make([]PreferSearchResponse, 0, len(genres))
		for _, g := range genres {
			url, err := aws.GetPresignURL(h.BucketBasics.AwsConfig, h.BucketBasics.BucketName, g.ImageKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Image URL From AWS"})
				return
			}
			res = append(res, PreferSearchResponse{
				ID:       g.ID,
				Name:     g.Name,
				ImageUrl: url,
			})
		}

		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(c.Writer)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(CommonRes{
			Message: "success",
			Data:    res,
		})
	}
}

func (h *UserHandler) AddUserGenre() gin.HandlerFunc {
	type req struct {
		GenreIDs []uint `json:"genre_ids" binding:"required"`
	}
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		var request req
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "check request params"})
			return
		}

		if err := h.userRepository.UpdateUserGenres(request.GenreIDs, userID.(uint)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fail to set user genre on db"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User genre updated successfully"})
	}
}

func (h *UserHandler) GetUserGenres() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		genres, err := h.userRepository.GetUserGenres(userID.(uint))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Fail to get genres"})
			return
		}
		genreIDs := make([]uint, 0, len(genres))
		for _, g := range genres {
			genreIDs = append(genreIDs, g.ID)
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    genreIDs,
		})
	}
}

// func (h *UserHandler) uploadProfileImage(fileHeader *multipart.FileHeader, userID uint) (string, error) {
// 	key := fmt.Sprintf("profile")

// 	uploadedKey, err := aws.UploadToS3(h.BucketBasics.S3Client, h.BucketBasics.BucketName, key, fileHeader)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to upload profile image: %w", err)
// 	}

// 	return uploadedKey, nil
// }

func (h *UserHandler) PutGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")

		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file "})
			return
		}
		url, err := aws.UploadToS3(h.BucketBasics.S3Client, h.BucketBasics.BucketName, fmt.Sprintf("genre/%s", name), file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to upload image: %v", err),
			})
			return
		}

		genre := &user.Genre{
			Name:     name,
			ImageKey: url,
		}

		if err = h.userRepository.PutGenre(genre); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save genre "})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	}
}
