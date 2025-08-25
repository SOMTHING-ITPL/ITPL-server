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

//performance table 채우는 cmd

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("fail to read .env file")
	}

	err = config.InitConfigs()
	db, err := storage.InitMySQL(*config.DbCfg)
	storage.AutoMigrate(db)

	repo := performance.NewRepository(db)

	scheduler := scheduler.PerformanceScheduler{
		PerformanceRepo: repo,
	}
	//start Date + 6month
	//running
	today := time.Now()

	// afterSixMonths := today.AddDate(0, 6, 0)
	afterSixMonths := today.AddDate(0, 0, 5)

	layout := "20060102"
	todayStr := today.Format(layout)
	afterSixMonthsStr := afterSixMonths.Format(layout)

	if err := scheduler.PutPerformanceList(todayStr, afterSixMonthsStr, nil, false); err != nil {
		fmt.Errorf("error is occur ! %s", err)
	}
	// if err := scheduler.PutPerformanceList(todayStr, afterSixMonthsStr, nil, false); err != nil {
	// 	fmt.Errorf("error is occur ! %s", err)
	// }

}
