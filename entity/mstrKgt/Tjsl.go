package mstrKgt

import (
	"time"

	"gorm.io/gorm"
)

type KegiatanMaster struct {
	IdKegiatan        int       `json:"id_kegiatan" gorm:"primary_key"`
	NamaKegiatan      string    `json:"nama_kegiatan"`
	DeskripsiKegiatan string    `json:"deskripsi_kegiatan"`
	CompCode          string    `json:"comp_code"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy         string    `json:"created_by"`
	Slug              string    `json:"slug"`
	Periode           string    `json:"periode"`
}

func (KegiatanMaster) TableName() string {
	return "tjsl.kegiatan_mstr"
}

type TJSLRepo struct {
	DB *gorm.DB
}

func NewTJSLRepo(db *gorm.DB) *TJSLRepo {
	return &TJSLRepo{DB: db}
}

func (t TJSLRepo) FindUserByCompCodeYear(comp_code string, tahun string) ([]KegiatanMaster, error) {
	var kgtn_mstr []KegiatanMaster
	err := t.DB.Where("comp_code=? AND periode=?", comp_code, tahun).Find(&kgtn_mstr).Error
	if err != nil {
		return kgtn_mstr, err
	}
	return kgtn_mstr, nil
}

func (t TJSLRepo) FindUserByID(id string) (KegiatanMaster, error) {
	var kgtn_mstr KegiatanMaster
	err := t.DB.Where("id_kegiatan=?", id).Take(&kgtn_mstr).Error
	if err != nil {
		return kgtn_mstr, err
	}
	return kgtn_mstr, nil
}

func (t TJSLRepo) FindNIKbySlug(slug string) (KegiatanMaster, error) {
	var kgtn_mstr KegiatanMaster
	err := t.DB.Where("slug=?", slug).Take(&kgtn_mstr).Error
	if err != nil {
		return kgtn_mstr, err
	}
	return kgtn_mstr, nil
}

func (t TJSLRepo) Create(km KegiatanMaster) (KegiatanMaster, error) {
	err := t.DB.Create(&km).Error
	return km, err
}

func (t TJSLRepo) Update(km KegiatanMaster) (KegiatanMaster, error) {
	err := t.DB.Save(&km).Error
	return km, err
}

func (t TJSLRepo) DelMasterKegiatanID(slug string) ([]KegiatanMaster, error) {
	var data []KegiatanMaster
	err := t.DB.Where("slug = ?", slug).First(&data).Error
	if err == nil {
		t.DB.Where("slug = ?", slug).Delete(&data)
		return data, nil
	}
	return data, err
}
