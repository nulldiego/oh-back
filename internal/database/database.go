package database

import (
	"fmt"

	"github.com/nulldiego/oh-back/config"
	"github.com/nulldiego/oh-back/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectAutoMigrateDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Conf.DbUser, config.Conf.DbPassword, config.Conf.DbHost, config.Conf.DbPort, config.Conf.DbName)

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&model.User{}, &model.Chat{}, &model.Message{})
}
