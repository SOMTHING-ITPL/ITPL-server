package main

import (
	"log"

	"context"

	"github.com/SOMTHING-ITPL/ITPL-server/aws/dynamo"
	"github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
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
	bucketService := s3.NewBucketBasics(awsCfg, &s3Cfg)

	// 5. 예시: 버킷 이름 출력
	log.Printf("S3 Bucket: %s", bucketService.BucketName)
	///

	dynamoClient := dynamo.NewDynamoDBClient(awsCfg)
	tableBasics := dynamo.NewTableBasics(dynamoClient, "itpl-message-db")

	// 채팅방 정보를 메모리 상에서 관리해야 함
	// DB에서 불러오면 새로운 ChatRoom 객체가 생성 -> 기존에 접속한 사용자들과 다른 객체가 됨
	// -> 기존 사용자들에게 메시지 브로드캐스트 불가
	// var ChatRooms = make(map[uint]*chat.ChatRoom)

	//db 쪽으로 빼야 하나?
	storage.AutoMigrate(db)
	r := server.SetupRouter(db, rdb, bucketService, &tableBasics)

	r.Run(":8080")
}
