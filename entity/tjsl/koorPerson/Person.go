package koorPerson

import (
	"time"

	"gorm.io/gorm"
)

type KoordinatorPerson struct {
	Id            int       `json:"id" gorm:"primary_key"`
	KoordinatorId int       `json:"koordinator_id"`
	NIK           string    `json:"nik"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (KoordinatorPerson) TableName() string {
	return "tjsl.koordinator_person"
}

type KoordinatorPersonRepo struct {
	DB *gorm.DB
}

func NewKoordinatorPerson(db *gorm.DB) *KoordinatorPersonRepo {
	return &KoordinatorPersonRepo{DB: db}
}
