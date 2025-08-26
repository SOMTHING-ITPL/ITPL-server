package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/scheduler"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/storage"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/joho/godotenv"
)

// 주기적으로 실행되어야 하는 스크립트 파일. 보완 필요함.
// 해당 로직을 그냥 서버에 걸어두는건?
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("fail to read .env file")
	}

	err = config.InitConfigs()
	db, err := storage.InitMySQL(*config.DbCfg)
	storage.AutoMigrate(db)

	rdb, err := storage.InitRedis(*config.RedisCfg)
	if err != nil {
		panic("Failed to init redis: " + err.Error())
	}

	repo := performance.NewRepository(db, rdb)

	scheduler := scheduler.PerformanceScheduler{
		PerformanceRepo: repo,
	}
	//start Date + 6month
	//running
	today := time.Now()

	afterSixMonths := today.AddDate(0, 6, 0)

	layout := "20060102"
	todayStr := today.Format(layout)
	afterSixMonthsStr := afterSixMonths.Format(layout)

	//공연예정 특정 일자 이후로 업데이트 된 것만 반영하면 됨.
	//updateAT 비교 후에 저장하는 형태
	if err := scheduler.PutPerformanceList(todayStr, afterSixMonthsStr, nil, false); err != nil {
		fmt.Errorf("error is occur ! %s", err)
	}
	//공연중
	if err := scheduler.PutPerformanceList(todayStr, afterSixMonthsStr, nil, true); err != nil {
		fmt.Errorf("error is occur ! %s", err)
	}
}
