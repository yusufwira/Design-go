package authentication

import (
	"time"

	"github.com/yusufwira/lern-golang-gin/entity/cuti"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
)

type ValidasiStoreCutiKaryawan struct {
	IdPengajuanAbsen interface{}                  `form:"id_pengajuan_absen" json:"id_pengajuan_absen"`
	Nik              string                       `form:"nik" json:"nik" binding:"required"`
	TipeAbsenId      string                       `form:"tipe_absen_id" json:"tipe_absen_id" binding:"required"`
	Deskripsi        string                       `form:"deskripsi" json:"deskripsi"`
	MulaiAbsen       string                       `form:"mulai_absen" json:"mulai_absen"`
	AkhirAbsen       string                       `form:"akhir_absen" json:"akhir_absen"`
	CreatedBy        string                       `form:"created_by" json:"created_by"`
	FileAbsen        []FileAbsenStoreCutiKaryawan `form:"file_absen" json:"file_absen"`
	Status           string                       `form:"status" json:"status"`
}

type SaldoCutiTransaksiPengajuan struct {
	cuti.SaldoCuti
	JmlhCuti int `json:"jumlah_cuti"`
}

type FileAbsenStoreCutiKaryawan struct {
	Filename  *string `form:"filename" json:"filename"`
	URL       *string `form:"url" json:"url"`
	Extension *string `form:"extension" json:"extension"`
}

type ValidasiStoreSaldoCuti struct {
	IDSaldo     interface{} `form:"id_saldo" json:"id_saldo"`
	TipeAbsenId string      `form:"tipe_absen_id" json:"tipe_absen_id" binding:"required"`
	ValidasiKonfirmasiNik
	Saldo           int    `form:"saldo" json:"saldo" binding:"required"`
	ValidFrom       string `form:"valid_from" json:"valid_from" binding:"required"`
	ValidTo         string `form:"valid_to" json:"valid_to" binding:"required"`
	CreatedBy       string `form:"created_by" json:"created_by"`
	Periode         string `form:"periode" json:"periode"`
	MaxHutang       int    `form:"max_hutang" json:"max_hutang"`
	ValidFromHutang string `form:"valid_from_hutang" json:"valid_from_hutang"`
}

type ValidasiListSaldoCuti struct {
	ValidasiKonfirmasiNik
	Key         string `form:"key" json:"key"`
	Perusahaan  string `form:"perusahaan" json:"perusahaan"`
	Direktorat  string `form:"direktorat" json:"direktorat"`
	Kompartemen string `form:"kompartemen" json:"kompartemen"`
	Departemen  string `form:"departemen" json:"departemen"`
	Tahun       string `form:"tahun" json:"tahun" binding:"required"`
}

type ValidationNIKTahun struct {
	ValidationLMK
}
type ValidationNIKTahunStatus struct {
	ValidationLMK
	Status    string `form:"status" json:"status"`
	IsManager bool   `form:"is_manager" json:"is_manager"`
}

type ValidationApprovalAtasanPengajuanAbsen struct {
	IdPengajuanAbsen int    `form:"id_pengajuan_absen" json:"id_pengajuan_absen" binding:"required"`
	Nik              string `form:"nik" json:"nik"`
	IsManager        bool   `form:"is_manager" json:"is_manager"`
	Status           string `form:"status" json:"status" binding:"required"`
	Keterangan       string `form:"keterangan" json:"keterangan"`
}

type SaldoCutiKaryawan struct {
	IdSaldoCuti     int       `json:"id_saldo_cuti"`
	TipeAbsenId     string    `json:"tipe_absen_id"`
	Nik             string    `json:"nik"`
	Saldo           int       `json:"saldo"`
	ValidFrom       string    `json:"valid_from"`
	ValidTo         string    `json:"valid_to"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Periode         string    `json:"periode"`
	MaxHutang       int       `json:"max_hutang"`
	ValidFromHutang string    `json:"valid_from_hutang"`
}

type ListSaldoCutiKaryawan struct {
	IdSaldoCuti               int `json:"id_saldo_cuti"`
	GetTipeAbsenKaryawanSaldo `json:"tipe_absen"`
	Nik                       string `json:"nik"`
	Nama                      string `json:"nama"`
	CompanyKaryawan           `json:"company"`
	Saldo                     int       `json:"saldo"`
	ValidFrom                 string    `json:"valid_from"`
	ValidTo                   string    `json:"valid_to"`
	CreatedBy                 string    `json:"created_by"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
	Periode                   string    `json:"periode"`
	MaxHutang                 int       `json:"max_hutang"`
	ValidFromHutang           string    `json:"valid_from_hutang"`
}

type CompanyKaryawan struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type GetTipeAbsenSaldoIndiv struct {
	IdTipeAbsen   string       `json:"id_tipe_absen"`
	NamaTipeAbsen string       `json:"nama_tipe_absen"`
	MaxAbsen      *int         `json:"max_absen"`
	TipeMaxAbsen  string       `json:"tipe_max_absen"`
	MySaldo       []SaldoIndiv `json:"my_saldo"`
}

type SaldoIndiv struct {
	Saldo           int    `json:"saldo"`
	ValidFrom       string `json:"valid_from"`
	ValidTo         string `json:"valid_to"`
	Periode         string `json:"periode"`
	MaxHutang       int    `json:"max_hutang"`
	ValidFromHutang string `json:"valid_from_hutang"`
}

type PengajuanAbsens struct {
	cuti.MyPengajuanAbsen
	File []cuti.FileAbsen `json:"files"`
}

type GetTipeAbsenKaryawanSaldo struct {
	IdTipeAbsen   string    `json:"id_tipe_absen"`
	NamaTipeAbsen string    `json:"nama_tipe_absen" gorm:"default:null"`
	CompCode      *string   `json:"comp_code" gorm:"default:null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	MaxAbsen      *int      `json:"max_absen" gorm:"default:null"`
	TipeMaxAbsen  *string   `json:"tipe_max_absen" gorm:"default:null"`
}

type ListApprovalCuti struct {
	IdPengajuanAbsen       int `json:"id_pengajuan_absen"`
	pihc.PihcMasterKaryRt  `json:"karyawan"`
	pihc.PihcMasterCompany `json:"companys"`
	cuti.TipeAbsen         `json:"tipe_absen"`
	Deskripsi              string                `json:"deskripsi"`
	TglPengajuan           string                `json:"tgl_pengajuan"`
	MulaiAbsen             string                `json:"mulai_absen"`
	AkhirAbsen             string                `json:"akhir_absen"`
	Status                 string                `json:"status"`
	FileAbsen              []cuti.FileAbsen      `json:"file_absen"`
	ApprovedBy             []cuti.AtasanApproved `json:"approved_by"`
	Periode                *string               `json:"periode"`
	Foto                   string                `json:"foto"`
	FotoDefault            string                `json:"foto_default"`
}
