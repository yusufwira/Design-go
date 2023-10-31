package authentication

import "time"

type ValidasiStoreCutiKaryawan struct {
	IDPengajuanAbsen int                          `form:"id_pengajuan_absen" json:"id_pengajuan_absen"`
	Nik              string                       `form:"nik" json:"nik" binding:"required"`
	TipeAbsenId      string                       `form:"tipe_absen_id" json:"tipe_absen_id" binding:"required"`
	Deskripsi        string                       `form:"deskripsi" json:"deskripsi"`
	MulaiAbsen       string                       `form:"mulai_absen" json:"mulai_absen"`
	AkhirAbsen       string                       `form:"akhir_absen" json:"akhir_absen"`
	CreatedBy        string                       `form:"created_by" json:"created_by"`
	FileAbsen        []FileAbsenStoreCutiKaryawan `form:"file_absen" json:"file_absen"`
}

type FileAbsenStoreCutiKaryawan struct {
	FileName  string `json:"file_name"`
	URL       string `json:"url"`
	Extension string `json:"extension"`
}

type ValidasiStoreSaldoCuti struct {
	IDSaldo     int    `form:"id_saldo" json:"id_saldo"`
	TipeAbsenId string `form:"tipe_absen_id" json:"tipe_absen_id" binding:"required"`
	ValidasiKonfirmasiNik
	Saldo           int    `form:"saldo" json:"saldo" binding:"required"`
	ValidFrom       string `form:"valid_from" json:"valid_from" binding:"required"`
	ValidTo         string `form:"valid_to" json:"valid_to" binding:"required"`
	CreatedBy       string `form:"created_by" json:"created_by"`
	MaxHutang       int    `form:"max_hutang" json:"max_hutang"`
	ValidFromHutang string `form:"valid_from_hutang" json:"valid_from_hutang"`
}

type ValidasiListSaldoCuti struct {
	ValidasiKonfirmasiNik
	Tahun string `form:"tahun" json:"tahun" binding:"required"`
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
	IdSaldoCuti     int       `json:"id_saldo_cuti"`
	TipeAbsenId     string    `json:"tipe_absen_id"`
	NamaTipeAbsen   string    `json:"nama_tipe_absen"`
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

type HistorySaldoCutiKaryawan struct {
	IdHistorySaldoCuti int       `json:"id_history_saldo_cuti"`
	TipeAbsenId        string    `json:"tipe_absen_id"`
	Nik                string    `json:"nik"`
	Saldo              int       `json:"saldo"`
	ValidFrom          string    `json:"valid_from"`
	ValidTo            string    `json:"valid_to"`
	CreatedBy          string    `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Periode            string    `json:"periode"`
}

type GetTipeAbsenSaldoIndiv struct {
	IdTipeAbsen   string `json:"id_tipe_absen"`
	NamaTipeAbsen string `json:"nama_tipe_absen"`
	Saldo         int    `json:"saldo"`
	ValidFrom     string `json:"valid_from"`
	ValidTo       string `json:"valid_to"`
}

type GetTipeAbsenKaryawan struct {
	IdTipeAbsen   string `json:"id_tipe_absen"`
	NamaTipeAbsen string `json:"nama_tipe_absen"`
}
