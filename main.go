package main

import (
	"log"

	"context"

	"github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/config"
	server "github.com/SOMTHING-ITPL/ITPL-server/internal/app"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/storage"
	aconfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/joho/godotenv"
)

func main() {
	//set env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//init function 쓰기. 패키지 로드될 때 init function 실행됨.
	err = config.InitConfigs()

	if err != nil {
		panic("Failed to load configs: " + err.Error())
	}

	db, err := storage.InitMySQL(*config.DbCfg)
	if err != nil {
		panic("Failed to init mysql: " + err.Error())
	}

	rdb, err := storage.InitRedis(*config.RedisCfg)
	if err != nil {
		panic("Failed to init redis: " + err.Error())
	}

	////////////////////////////////
	var s3Cfg config.S3Config
	s3Cfg.Load()

	// 3. AWS SDK 기본 설정 로드
	awsCfg, err := aconfig.LoadDefaultConfig(context.TODO(),
		aconfig.WithRegion("ap-northeast-2"), // 서울 리전
	)
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	// 4. BucketBasics 객체 생성
	bucketService := aws.NewBucketBasics(awsCfg, &s3Cfg)

	// 5. 예시: 버킷 이름 출력
	log.Printf("S3 Bucket: %s", bucketService.BucketName)
	///

	//db 쪽으로 빼야 하나?
	storage.AutoMigrate(db)
	r := server.SetupRouter(db, rdb, bucketService)

	r.Run(":8080")
}
