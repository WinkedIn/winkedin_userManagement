package services

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBConnection(v *viper.Viper, dbModels ...interface{}) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		v.GetString("database.host"),
		v.GetString("database.user"),
		v.GetString("database.password"),
		v.GetString("database.dbname"),
		v.GetInt("database.port"),
		v.GetBool("database.sslmode"),
	)
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	db.AutoMigrate(dbModels)
	return db, nil
}
