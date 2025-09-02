package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/gin-gonic/gin"
)

func (h *PlaceHandler) WriteReviewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		splaceId := c.PostForm("place_id")
		placeID, err := strconv.ParseUint(splaceId, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
			return
		}
		placeId := uint(placeID)
		text := c.PostForm("text")
		srating := c.PostForm("rating")
		rating, err := strconv.ParseFloat(srating, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating"})
			return
		}
		var imgUrl []place.ReviewImage

		// 이미지 파일 받기
		form, _ := c.MultipartForm()
		files := form.File["images"]
		for _, fileHeader := range files {
			key, err := aws.UploadToS3(h.BucketBasics.S3Client, h.BucketBasics.BucketName, fmt.Sprintf("reviews/%d", placeId), fileHeader)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
				return
			}
			imgUrl = append(imgUrl, place.ReviewImage{Key: key})
		}
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)
		user, err := h.userRepository.GetById(userID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err := place.WriteReview(h.database, placeId, user, text, rating, imgUrl); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write review: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Review added successfully"})
	}
}

func (h *PlaceHandler) DeleteReviewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		revId, err := strconv.ParseUint(c.Param("review_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
			return
		}
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)

		rev, err := place.GetReviewByID(h.database, uint(revId))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
			return
		}

		if rev.UserId != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own reviews"})
			return
		}

		err = place.DeleteReview(h.database, uint(revId), *h.BucketBasics)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
	}
}

func (h *PlaceHandler) GetPlaceReviewsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		placeID, err := strconv.ParseUint(c.Param("place_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
			return
		}
		reviews, err := place.GetPlaceReviews(h.database, uint(placeID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews: " + err.Error()})
			return
		}
		var response []PlaceReviewResponse
		for _, r := range reviews {
			var imgs []ReviewImageResponse
			for _, img := range r.Images {
				url, _ := aws.GetPresignURL(h.BucketBasics.AwsConfig, h.BucketBasics.BucketName, img.Key)
				imgs = append(imgs, ReviewImageResponse{URL: url})
			}
			response = append(response, PlaceReviewResponse{
				ID:           r.ID,
				UserID:       r.UserId,
				UserNickname: r.UserNickName,
				Rating:       r.Rating,
				Comment:      r.Comment,
				Images:       imgs,
			})
		}

		c.JSON(http.StatusOK, gin.H{"reviews": response})
	}
}

func (h *PlaceHandler) GetMyReviewsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)

		user, err := h.userRepository.GetById(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		reviews, err := place.GetMyReviews(h.database, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews: " + err.Error()})
			return
		}
		var response []PlaceReviewResponse
		for _, r := range reviews {
			var imgs []ReviewImageResponse
			for _, img := range r.Images {
				url, _ := aws.GetPresignURL(h.BucketBasics.AwsConfig, h.BucketBasics.BucketName, img.Key)
				imgs = append(imgs, ReviewImageResponse{URL: url})
			}
			response = append(response, PlaceReviewResponse{
				ID:           r.ID,
				UserID:       r.UserId,
				UserNickname: r.UserNickName,
				Rating:       r.Rating,
				Comment:      r.Comment,
				Images:       imgs,
			})
		}
		c.JSON(http.StatusOK, gin.H{"reviews": response})
	}
}
