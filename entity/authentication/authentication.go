package authentication

import (
	"time"

	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/tjsl"
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
	IdKegiatan        int    `json:"id" gorm:"primary_key"`
	NIK               string `json:"nik" binding:"required"`
	NamaKegiatan      string `json:"nama_kegiatan" binding:"required"`
	DeskripsiKegiatan string `json:"deskripsi_kegiatan"`
}

type ValidationSKKgt struct {
	Id                int    `json:"id" gorm:"primary_key"`
	NIK               string `json:"nik" binding:"required"`
	KegiatanParentId  int    `json:"kegiatan_parent_id"`
	KoordinatorId     int    `json:"koordinator_id"`
	NamaKegiatan      string `json:"nama_kegiatan" binding:"required"`
	TanggalKegiatan   string `json:"tanggal" binding:"required"`
	LokasiKegiatan    string `json:"lokasi" binding:"required"`
	DeskripsiKegiatan string `json:"deskripsi" binding:"required"`
	Photos            []struct {
		IDPhoto      int    `json:"id_photo"`
		OriginalName string `json:"original_name"`
		URL          string `json:"url"`
	} `json:"photos"`
	Tahun string `json:"tahun" binding:"required"`
}

type ValidationListApproval struct {
	NIK      string `json:"nik" binding:"required"`
	Tahun    string `json:"tahun" binding:"required"`
	CompCode string `json:"comp_code"`
}

type ValidationApprovalAtasan struct {
	SlugKegiatan string `form:"slug_kegiatan" binding:"required"`
	Status       string `form:"status" binding:"required"`
}

type ValidationMyTjsl struct {
	Nik   string `form:"nik" binding:"required"`
	Tahun string `form:"tahun" binding:"required"`
}

type ValidationListKoordinator struct {
	NIK   string `json:"nik" binding:"required"`
	Tahun string `json:"tahun" binding:"required"`
	Slug  string `json:"slug"`
}

type KegiatanKaryawanPhotos struct {
	IDKegiatan               int                   `json:"id_kegiatan"`
	SlugKegiatan             string                `json:"slug_kegiatan"`
	Nik                      string                `json:"nik"`
	Nama                     string                `json:"nama"`
	PhotoProfile             string                `json:"photo_profile"`
	DeptTitle                string                `json:"dept_title"`
	Jenis                    string                `json:"jenis"`
	KoordinatorID            int                   `json:"koordinator_id"`
	SlugKoordinator          int                   `json:"slug_koordinator"`
	SlugKegiatanParent       int                   `json:"slug_kegiatan_parent"`
	KegiatanParentID         int                   `json:"kegiatan_parent_id"`
	NamaKegiatan             string                `json:"nama_kegiatan"`
	TanggalKegiatan          string                `json:"tanggal_kegiatan"`
	TanggalKegiatanNonFormat string                `json:"tanggal_kegiatan_non_format"`
	LokasiKegiatan           string                `json:"lokasi_kegiatan"`
	Deskripsi                string                `json:"deskripsi"`
	Status                   string                `json:"status"`
	AlasanPenolakan          string                `json:"alasan_penolakan"`
	PhotoKegiatan            []tjsl.KegiatanPhotos `json:"photo_kegiatan"`
	Tahun                    string                `json:"tahun"`
}

type ListApprovalTJSL struct {
	SlugKegiatan    string                `json:"slug_kegiatan"`
	Nik             string                `json:"nik"`
	Nama            string                `json:"nama"`
	PhotoProfile    string                `json:"photo_profile"`
	Email           string                `json:"email"`
	PosID           string                `json:"pos_id"`
	PosTitle        string                `json:"pos_title"`
	DeptTitle       string                `json:"dept_title"`
	Jenis           string                `json:"jenis"`
	NamaKegiatan    string                `json:"nama_kegiatan"`
	TanggalKegiatan string                `json:"tanggal_kegiatan"`
	LokasiKegiatan  string                `json:"lokasi_kegiatan"`
	Deskripsi       string                `json:"deskripsi"`
	Status          string                `json:"status"`
	PhotoKegiatan   []tjsl.KegiatanPhotos `json:"photo_kegiatan"`
	Short           string                `json:"short"`
	LogoCompany     string                `json:"logo_company"`
}

type ValidationKKoor struct {
	Id     int    `json:"id"`
	Nama   string `json:"nama" binding:"required"`
	Nik    string `json:"nik" binding:"required"`
	Photos []struct {
		IDPhoto      string `json:"id_photo"`
		Extension    string `json:"extension"`
		Name         string `json:"name"`
		OriginalName string `json:"original_name"`
		Size         string `json:"size"`
		URL          string `json:"url"`
	} `json:"photos" binding:"required"`
	Person []string `json:"person" binding:"required,min=2"`
	Tahun  string   `json:"tahun" binding:"required"`
}

type KegiatanDetailKoordinatorPhotos struct {
	IDKoordinator    int                   `json:"id_koordinator"`
	KegiatanParentID int                   `json:"kegiatan_parent_id"`
	Nama             string                `json:"nama"`
	CreatedBy        string                `json:"created_by"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	CompCode         string                `json:"comp_code"`
	Slug             string                `json:"slug"`
	Periode          string                `json:"periode"`
	Photos           []tjsl.KegiatanPhotos `json:"photos"`
	Person           []Personal            `json:"person"`
}

type Personal struct {
	ID            int                   `json:"id"`
	KoordinatorID int                   `json:"koordinator_id"`
	Nik           string                `json:"nik"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	Employee      pihc.PihcMasterKaryRt `json:"employee"`
	URLPhoto      string                `json:"url_photo"`
}

type KegiatanListKoordinatorPhotos struct {
	IDKoordinator    int                   `json:"id_koordinator"`
	KegiatanParentID int                   `json:"kegiatan_parent_id"`
	Nama             string                `json:"nama"`
	CreatedBy        string                `json:"created_by"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	CompCode         string                `json:"comp_code"`
	Slug             string                `json:"slug"`
	Periode          string                `json:"periode"`
	Employee         pihc.PihcMasterKaryRt `json:"employee"`
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
