package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/scheduler"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/storage"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/update_client"
	"github.com/joho/godotenv"
)

// 최초 1회만 실행되면 되는 코드임. +6month 까지만 담아놓는거.
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
		// not use redis in real code
		// panic("Failed to init redis: " + err.Error())
	}

	repo := performance.NewRepository(db, rdb)

	scheduler := scheduler.PerformanceScheduler{
		PerformanceRepo: repo,
	}
	//start Date + 6month
	//running 3일 전 업데이트 된 거
	today := time.Now()
	startDate := today.AddDate(0, 0, 5)

	afterSixMonths := today.AddDate(0, 6, 0)

	layout := "20060102"
	todayStr := today.Format(layout)
	startDayStr := startDate.Format(layout)
	afterSixMonthsStr := afterSixMonths.Format(layout)

	cli, err := update_client.NewClient(config.GrpcCfg.Host + ":" + config.GrpcCfg.UpdatePort)
	if err != nil {
		log.Fatal("Failed to create update client:", err)
	}

	//걍 자정 이후로 추가되는 데이터 있으면 가져오면 될 것 같은데? 추가로 공연중 / 공연예정도 다 담아야함.
	if err := scheduler.PutPerformanceList(startDayStr, afterSixMonthsStr, false, &todayStr, cli); err != nil {
		fmt.Errorf("error is occur ! %s", err)
	}

	if err := scheduler.UpdateStatusList(); err != nil {
		fmt.Errorf("error is occur: update status! %s", err)
	}
}
