package authentication

import (
	"time"

	"github.com/yusufwira/lern-golang-gin/entity/tjsl/photosKgt"
)

type ValidationLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ValidationLMK struct {
	NIK   string `json:"nik" binding:"required"`
	Tahun string `json:"tahun" binding:"required"`
}

type ValidationSMK struct {
	IdKegiatan        int       `json:"id" gorm:"primary_key"`
	NIK               string    `json:"nik" binding:"required"`
	NamaKegiatan      string    `json:"nama_kegiatan"`
	DeskripsiKegiatan string    `json:"deskripsi_kegiatan"`
	Slug              string    `json:"slug"`
	Periode           string    `json:"periode"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type ValidationKK struct {
	Id                int    `json:"id" gorm:"primary_key"`
	NIK               string `json:"nik" binding:"required"`
	KegiatanParentId  int    `json:"kegiatan_parent_id"`
	KoordinatorId     int    `json:"koordinator_id"`
	NamaKegiatan      string `json:"nama_kegiatan"`
	TanggalKegiatan   string `json:"tanggal"`
	LokasiKegiatan    string `json:"lokasi"`
	DeskripsiKegiatan string `json:"deskripsi"`
	Photos            []struct {
		IDPhoto      int    `json:"id_photo"`
		OriginalName string `json:"original_name"`
		URL          string `json:"url"`
	} `json:"photos"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type KegiatanKaryawanPhotos struct {
	IDKegiatan               int                        `json:"id_kegiatan"`
	SlugKegiatan             string                     `json:"slug_kegiatan"`
	Nik                      string                     `json:"nik"`
	Nama                     string                     `json:"nama"`
	PhotoProfile             string                     `json:"photo_profile"`
	DeptTitle                string                     `json:"dept_title"`
	Jenis                    string                     `json:"jenis"`
	KoordinatorID            int                        `json:"koordinator_id"`
	SlugKoordinator          int                        `json:"slug_koordinator"`
	SlugKegiatanParent       int                        `json:"slug_kegiatan_parent"`
	KegiatanParentID         int                        `json:"kegiatan_parent_id"`
	NamaKegiatan             string                     `json:"nama_kegiatan"`
	TanggalKegiatan          string                     `json:"tanggal_kegiatan"`
	TanggalKegiatanNonFormat string                     `json:"tanggal_kegiatan_non_format"`
	LokasiKegiatan           string                     `json:"lokasi_kegiatan"`
	Deskripsi                string                     `json:"deskripsi"`
	Status                   string                     `json:"status"`
	AlasanPenolakan          string                     `json:"alasan_penolakan"`
	PhotoKegiatan            []photosKgt.KegiatanPhotos `json:"photo_kegiatan"`
	Tahun                    string                     `json:"tahun"`
}
