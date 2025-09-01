package aws

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewBucketBasics(cfg aws.Config, s3Cfg *config.S3Config) *BucketBasics {
	client := s3.NewFromConfig(cfg)
	return &BucketBasics{
		S3Client:   client,
		BucketName: s3Cfg.BucketName,
		AwsConfig:  cfg,
	}
}

func UploadToS3(client *s3.Client, bucket, prefix string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	var contentType string
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	default:
		contentType = "image/jpeg"
	}

	// file seek
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// key (ex: reviews/{prefix}/{filename})
	key := fmt.Sprintf("%s/%s", prefix, fileHeader.Filename)

	// S3 upload
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to s3: %w", err)
	}

	return key, nil
}

func DeleteImage(client *s3.Client, bucketName, KeyName string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(KeyName),
	}
	_, err := client.DeleteObject(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

func GetPresignURL(cfg aws.Config, bucketName, keyName string) (string, error) {
	s3client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3client)
	presignedUrl, err := presignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(keyName),
		},
		s3.WithPresignExpires(time.Minute*15))
	if err != nil {
		return "", err
	}
	return presignedUrl.URL, nil
}
