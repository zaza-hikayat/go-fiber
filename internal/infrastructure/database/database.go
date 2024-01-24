package database

import (
	"fmt"

	"github.com/zaza-hikayat/go-fiber/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDBConnection(conf configs.Config) (*gorm.DB, error) {
	logConfig := logger.Default.LogMode(logger.Info)
	configStr := "Info"

	if !conf.Server.IsLogging {
		logConfig = logger.Default.LogMode(logger.Silent)
		configStr = "Silent"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		conf.Database.Host,
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Name,
		conf.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logConfig,
	})
	fmt.Printf("Check DB Log Config %s : %t", configStr, conf.Server.IsLogging)

	if err != nil {
		panic(err)
	}
	return db, err
}
