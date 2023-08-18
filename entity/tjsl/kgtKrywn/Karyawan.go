package kgtKrywn

import (
	"time"

	"gorm.io/gorm"
)

type KegiatanKaryawan struct {
	Id                int       `json:"id" gorm:"primary_key"`
	NIK               string    `json:"nik"`
	KegiatanParentId  int       `json:"kegiatan_parent_id"`
	KoordinatorId     int       `json:"koordinator_id"`
	NamaKegiatan      string    `json:"nama_kegiatan"`
	TanggalKegiatan   string    `json:"tanggal_kegiatan"`
	LokasiKegiatan    string    `json:"lokasi_kegiatan"`
	DeskripsiKegiatan string    `json:"deskripsi_kegiatan"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Manager           string    `json:"manager"`
	Slug              string    `json:"slug"`
	DescDecline       string    `json:"desc_decline"`
	CompCode          string    `json:"comp_code"`
	Periode           string    `json:"periode"`
}

func (KegiatanKaryawan) TableName() string {
	return "tjsl.kegiatan_karyawan"
}

type KegiatanKaryawanRepo struct {
	DB *gorm.DB
}

func NewKegiatanKaryawanRepo(db *gorm.DB) *KegiatanKaryawanRepo {
	return &KegiatanKaryawanRepo{DB: db}
}

func (t KegiatanKaryawanRepo) Create(kk KegiatanKaryawan) (KegiatanKaryawan, error) {
	err := t.DB.Create(&kk).Error
	return kk, err
}

func (t KegiatanKaryawanRepo) Update(kk KegiatanKaryawan) (KegiatanKaryawan, error) {
	err := t.DB.Save(&kk).Error
	return kk, err
}

func (t KegiatanKaryawanRepo) FindData(id int) KegiatanKaryawan {
	var kgtn_krywn KegiatanKaryawan
	t.DB.Where("id=?", id).Find(&kgtn_krywn)
	return kgtn_krywn
}
