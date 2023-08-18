package koorKgt

import (
	"time"

	"gorm.io/gorm"
)

type KegiatanKoordinator struct {
	IdKoordinator    int       `json:"id_koordinator" gorm:"primary_key"`
	KegiatanParentId int       `json:"kegiatan_parent_id"`
	Nama             string    `json:"nama"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CompCode         string    `json:"comp_code"`
	Slug             string    `json:"slug"`
	Periode          string    `json:"periode"`
}

func (KegiatanKoordinator) TableName() string {
	return "tjsl.kegiatan_koordinator"
}

type KegiatanKoordinatorRepo struct {
	DB *gorm.DB
}

func NewKegiatanKoordinatorRepo(db *gorm.DB) *KegiatanKoordinatorRepo {
	return &KegiatanKoordinatorRepo{DB: db}
}
