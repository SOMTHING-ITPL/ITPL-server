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
		//no need to use redis
		// panic("Failed to init redis: " + err.Error())
	}

	repo := performance.NewRepository(db, rdb)

	scheduler := scheduler.PerformanceScheduler{
		PerformanceRepo: repo,
	}
	//start Date + 6month
	//running
	today := time.Now()

	afterSixMonths := today.AddDate(0, 0, -10) //6으로 변경해야함.

	layout := "20060102"
	todayStr := today.Format(layout)
	afterSixMonthsStr := afterSixMonths.Format(layout)

	//공연예정 -> 잠시
	if err := scheduler.PutPerformanceList(afterSixMonthsStr, todayStr, false, nil); err != nil {
		fmt.Errorf("error is occur ! %s", err)
	}
	//공연중
	if err := scheduler.PutPerformanceList(todayStr, afterSixMonthsStr, true, nil); err != nil {
		fmt.Errorf("error is occur ! %s", err)
	}

	fmt.Println("Performance list updated successfully!")
}
