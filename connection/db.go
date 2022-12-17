package connection

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Database() gorm.DB {
	userDB := "root"
	passDB := ""
	nameDB := "golang"
	dsn := userDB + passDB + "@tcp(127.0.0.1:3306)/" + nameDB + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return *db

}
