package authentication

import (
	"time"

	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl"
)

type ValidationLogin struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type ValidationRegister struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Name     string `json:"name" form:"name" binding:"required"`
	Nik      string `json:"nik" form:"nik" binding:"required"`
	Type     string `json:"type" form:"type"`
}

type ValidationLMK struct {
	NIK   string `json:"nik" form:"nik" binding:"required"`
	Tahun string `json:"tahun" form:"tahun" binding:"required"`
}

type ValidationSMK struct {
	IdKegiatan        int     `json:"id" form:"id" gorm:"primary_key"`
	NIK               string  `json:"nik" form:"nik" binding:"required"`
	NamaKegiatan      string  `json:"nama_kegiatan" form:"nama_kegiatan" binding:"required"`
	DeskripsiKegiatan *string `json:"deskripsi_kegiatan" form:"deskripsi_kegiatan"`
}

type ValidationSKKgt struct {
	Id                int     `json:"id" gorm:"primary_key"`
	NIK               string  `json:"nik" binding:"required"`
	KegiatanParentId  *int    `json:"kegiatan_parent_id" gorm:"default:null"`
	KoordinatorId     *int    `json:"koordinator_id" gorm:"default:null"`
	NamaKegiatan      string  `json:"nama_kegiatan" binding:"required"`
	TanggalKegiatan   string  `json:"tanggal" binding:"required"`
	Status            string  `json:"status" binding:"required"`
	Manager           *string `json:"manager"`
	LokasiKegiatan    string  `json:"lokasi" binding:"required"`
	DeskripsiKegiatan *string `json:"deskripsi" binding:"required"`
	Photos            []struct {
		IDPhoto      int    `json:"id_photo"`
		OriginalName string `json:"original_name"`
		URL          string `json:"url"`
	} `json:"photos"`
	Tahun string `json:"tahun" binding:"required"`
}

type ValidationListApproval struct {
	ValidationLMK
	CompCode string `json:"comp_code"`
}
type ValidationGetLeaderBoard struct {
	ValidationLMK
	IsMobile int    `json:"isMobile" form:"isMobile"`
	Company  string `json:"company" form:"company" binding:"required"`
}

type ValidationApprovalAtasan struct {
	SlugKegiatan string `form:"slug_kegiatan" binding:"required"`
	Status       string `form:"status" binding:"required"`
}

type ValidationMyTjsl struct {
	ValidationLMK
}

type ValidationListKoordinator struct {
	ValidationLMK
	Slug string `json:"slug" form:"slug"`
}

type KegiatanKaryawanPhotos struct {
	IDKegiatan               int                   `json:"id_kegiatan"`
	SlugKegiatan             string                `json:"slug_kegiatan"`
	Nik                      string                `json:"nik"`
	Nama                     string                `json:"nama"`
	PhotoProfile             string                `json:"photo_profile"`
	DeptTitle                string                `json:"dept_title"`
	Jenis                    string                `json:"jenis"`
	KoordinatorID            *int                  `json:"koordinator_id"`
	SlugKoordinator          *string               `json:"slug_koordinator"`
	SlugKegiatanParent       *string               `json:"slug_kegiatan_parent"`
	KegiatanParentID         *int                  `json:"kegiatan_parent_id"`
	NamaKegiatan             string                `json:"nama_kegiatan"`
	TanggalKegiatan          string                `json:"tanggal_kegiatan"`
	TanggalKegiatanNonFormat string                `json:"tanggal_kegiatan_non_format"`
	LokasiKegiatan           string                `json:"lokasi_kegiatan"`
	Deskripsi                *string               `json:"deskripsi"`
	Status                   string                `json:"status"`
	AlasanPenolakan          *string               `json:"alasan_penolakan"`
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
	Deskripsi       *string               `json:"deskripsi"`
	Status          string                `json:"status"`
	PhotoKegiatan   []tjsl.KegiatanPhotos `json:"photo_kegiatan"`
	Short           *string               `json:"short"`
	LogoCompany     string                `json:"logo_company"`
}

type ListChartSummary struct {
	RekapPerbulan `json:"month"`
	Employee      `json:"employee"`
}

type ListChartNotFoundDataSummary struct {
	RekapPerbulan `json:"month"`
}

type RekapPerbulan struct {
	Month
	TotalIndividu int `json:"total_individu"`
}

type Month struct {
	Num1  int `json:"1"`
	Num2  int `json:"2"`
	Num3  int `json:"3"`
	Num4  int `json:"4"`
	Num5  int `json:"5"`
	Num6  int `json:"6"`
	Num7  int `json:"7"`
	Num8  int `json:"8"`
	Num9  int `json:"9"`
	Num10 int `json:"10"`
	Num11 int `json:"11"`
	Num12 int `json:"12"`
}

type Employee struct {
	EmpNama   string `json:"emp_nama"`
	Nik       string `json:"nik"`
	PosID     string `json:"pos_id"`
	PosTitle  string `json:"pos_title"`
	DeptID    string `json:"dept_id"`
	DeptTitle string `json:"dept_title"`
	KompID    string `json:"komp_id"`
	KompTitle string `json:"komp_title"`
	DirID     string `json:"dir_id"`
	DirTitle  string `json:"dir_title"`
	Photo     string `json:"photo"`
}

type ValidationKKoor struct {
	Id               int    `json:"id"`
	KegiatanParentId *int   `json:"kegiatan_parent_id"`
	Nama             string `json:"nama" form:"nama" binding:"required"`
	Nik              string `json:"nik" form:"nik" binding:"required"`
	Photos           []struct {
		IDPhoto      string `json:"id_photo" form:"id_photo"`
		Extension    string `json:"extension" form:"extension"`
		Name         string `json:"name" form:"name"`
		OriginalName string `json:"original_name" form:"original_name"`
		Size         string `json:"size" form:"size"`
		URL          string `json:"url" form:"url"`
	} `json:"photos" form:"photos[]" binding:"required"`
	Person []string `json:"person" form:"person[]" binding:"required,min=2"`
	Tahun  string   `json:"tahun" form:"tahun" binding:"required"`
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
	ID            int                 `json:"id"`
	KoordinatorID int                 `json:"koordinator_id"`
	Nik           string              `json:"nik"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Employee      pihc.PihcMasterKary `json:"employee"`
	URLPhoto      string              `json:"url_photo"`
}

type KegiatanListKoordinatorPhotos struct {
	tjsl.KegiatanKoordinator
	Employee pihc.PihcMasterKaryRt `json:"employee" gorm:"foreignkey:EmpNo;association_foreignkey:CreatedBy"`
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
