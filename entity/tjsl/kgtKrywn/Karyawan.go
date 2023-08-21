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
	if err != nil {
		return kk, err
	}
	return kk, nil
}

func (t KegiatanKaryawanRepo) Update(kk KegiatanKaryawan) (KegiatanKaryawan, error) {
	err := t.DB.Save(&kk).Error
	if err != nil {
		return kk, err
	}
	return kk, nil
}

func (t KegiatanKaryawanRepo) FindDataID(id int) (KegiatanKaryawan, error) {
	var kgtn_krywn KegiatanKaryawan
	err := t.DB.Where("id=?", id).First(&kgtn_krywn).Error
	if err != nil {
		return kgtn_krywn, err
	}
	return kgtn_krywn, nil
}

func (t KegiatanKaryawanRepo) FindDataSlug(slug string) (KegiatanKaryawan, error) {
	var kgtn_krywn KegiatanKaryawan
	err := t.DB.Where("slug=?", slug).First(&kgtn_krywn).Error
	if err != nil {
		return kgtn_krywn, err
	}
	return kgtn_krywn, nil
}

func (t KegiatanKaryawanRepo) FindDataNIKPeriode(nik string, tahun string) ([]KegiatanKaryawan, error) {
	var kgtn_krywn []KegiatanKaryawan
	err := t.DB.Where("nik=? AND periode=?", nik, tahun).Find(&kgtn_krywn).Error
	if err != nil {
		return kgtn_krywn, err
	}
	return kgtn_krywn, nil
}

func (t KegiatanKaryawanRepo) DelKegiatanKaryawanID(slug string, status string) error {
	var data []KegiatanKaryawan
	err := t.DB.Where("slug = ? AND status=?", slug, status).First(&data).Error
	if err == nil {
		t.DB.Where("slug = ? AND status=?", slug, status).Delete(&data)
		return nil
	}
	return err
}
