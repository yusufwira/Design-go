package connection

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Database() *gorm.DB {
	dsn := "postgres://postgres:V3ry5tr0n94dm1nP@$$w0rd@192.168.188.232:5432/pi-smart"
	//dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//db, err := gorm.Open(postgres.Open("postgres://postgres:postgres@localhost:5432/postgres"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
