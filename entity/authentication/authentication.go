package authentication

import (
	"time"
)

type AuthenticationInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthenticationLMK struct {
	NIK   string `json:"nik" binding:"required"`
	Tahun string `json:"tahun" binding:"required"`
}

type AuthenticationSMK struct {
	IdKegiatan        int       `json:"id" gorm:"primary_key"`
	NIK               string    `json:"nik" binding:"required"`
	NamaKegiatan      string    `json:"nama_kegiatan"`
	DeskripsiKegiatan string    `json:"deskripsi_kegiatan"`
	Slug              string    `json:"slug"`
	Periode           string    `json:"periode"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type AuthenticationKK struct {
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
