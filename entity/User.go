package entity

import (
	"time"

	"github.com/yusufwira/lern-golang-gin/connection"
)

type User struct {
	Username   string
	Password   string
	Name       string
	Email      string
	created_at time.Time `json:"created_at" gorm:"autoCreateTime;not null"`
	updated_at time.Time `json:"updated_at" gorm:"autoCreateTime;not null"`
}

func (user User) Save() {
	db := connection.Database()
	db.Table("User").Create(user)

}

func (user User) GetAll() []User {
	db := connection.Database()
	var data []User
	db.Table("User").Find(&data)
	return data
}
