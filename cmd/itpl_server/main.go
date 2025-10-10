package main

import (
	"log"

	"context"

	"github.com/SOMTHING-ITPL/ITPL-server/aws/dynamo"
	aws_client "github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
	"github.com/SOMTHING-ITPL/ITPL-server/config"
	server "github.com/SOMTHING-ITPL/ITPL-server/internal/app"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/storage"
	aconfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/joho/godotenv"
)

func main() {
	//set env file
	err := godotenv.Load()

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

	var s3Cfg config.S3Config
	s3Cfg.Load()

	awsCfg, err := aconfig.LoadDefaultConfig(context.TODO(),
		aconfig.WithRegion("ap-northeast-2"),
	)
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	bucketService := aws_client.NewBucketBasics(awsCfg, &s3Cfg)

	log.Printf("S3 Bucket: %s", bucketService.BucketName)

	// DB Configuration
	dynamoClient := dynamo.NewDynamoDBClient(awsCfg)
	tableBasics := dynamo.NewTableBasics(dynamoClient, "itpl-message-db")

	storage.AutoMigrate(db)
	r := server.SetupRouter(db, rdb, bucketService, &tableBasics)

	r.Run(":8080")
}
