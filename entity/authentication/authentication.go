package authentication

import "time"

type AuthenticationInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthenticationLMK struct {
	NIK   string `json:"nik" binding:"required"`
	Tahun string `json:"tahun" binding:"required"`
}

type AuthenticationSMK struct {
	NIK               string    `json:"nik" binding:"required"`
	NamaKegiatan      string    `json:"nama_kegiatan"`
	DeskripsiKegiatan string    `json:"deskripsi_kegiatan"`
	Slug              string    `json:"slug"`
	Periode           string    `json:"periode"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
