package cuti

import (
	"time"

	"gorm.io/gorm"
)

type PengajuanAbsen struct {
	IDPengajuanAbsen   int       `json:"id_pengajuan_absen"`
	Nik                string    `json:"nik" gorm:"default:null"`
	CompCode           string    `json:"comp_code" gorm:"default:null"`
	TipeAbsenId        *string   `json:"tipe_absen_id" gorm:"default:null"`
	Deskripsi          *string   `json:"deskripsi" gorm:"default:null"`
	MulaiAbsen         time.Time `json:"mulai_absen" gorm:"default:null"`
	AkhirAbsen         time.Time `json:"akhir_absen" gorm:"default:null"`
	TglPengajuan       time.Time `json:"tgl_pengajuan" gorm:"default:null"`
	Status             *string   `json:"status" gorm:"default:null"`
	CreatedBy          *string   `json:"created_by" gorm:"default:null"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Keterangan         *string   `json:"keterangan" gorm:"default:null"`
	Periode            *string   `json:"periode" gorm:"default:null"`
	JumlahHariKalender *int      `json:"jumlah_hari_kalender" gorm:"default:null"`
	JumlahHariKerja    *int      `json:"jumlah_hari_kerja" gorm:"default:null"`
}

type HistoryPengajuanAbsen struct {
	IDHistoryPengajuanAbsen int       `json:"id_history_pengajuan_absen"`
	Nik                     string    `json:"nik" gorm:"default:null"`
	CompCode                string    `json:"comp_code" gorm:"default:null"`
	TipeAbsenId             *string   `json:"tipe_absen_id" gorm:"default:null"`
	Deskripsi               *string   `json:"deskripsi" gorm:"default:null"`
	MulaiAbsen              time.Time `json:"mulai_absen" gorm:"default:null"`
	AkhirAbsen              time.Time `json:"akhir_absen" gorm:"default:null"`
	TglPengajuan            time.Time `json:"tgl_pengajuan" gorm:"default:null"`
	Status                  *string   `json:"status" gorm:"default:null"`
	CreatedBy               *string   `json:"created_by" gorm:"default:null"`
	CreatedAt               time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt               time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Keterangan              *string   `json:"keterangan" gorm:"default:null"`
	Periode                 *string   `json:"periode" gorm:"default:null"`
	ApprovedBy              *string   `json:"approved_by" gorm:"default:null"`
	JumlahHariKalender      *int      `json:"jumlah_hari_kalender" gorm:"default:null"`
	JumlahHariKerja         *int      `json:"jumlah_hari_kerja" gorm:"default:null"`
}

type FileAbsen struct {
	IdFileAbsen      int       `json:"id_file_absen" gorm:"primary_key"`
	PengajuanAbsenId string    `json:"pengajuan_absen_id" gorm:"default:null"`
	FileName         *string   `json:"filename" gorm:"default:null"`
	Url              *string   `json:"url" gorm:"default:null"`
	Extension        *string   `json:"extension" gorm:"default:null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type TipeAbsen struct {
	IdTipeAbsen   string    `json:"id_tipe_absen"`
	NamaTipeAbsen string    `json:"nama_tipe_absen" gorm:"default:null"`
	CompCode      *string   `json:"comp_code" gorm:"default:null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	MaxAbsen      *int      `json:"max_absen" gorm:"default:null"`
	TipeMaxAbsen  *string   `json:"tipe_max_absen" gorm:"default:null"`
}

type SaldoCuti struct {
	IdSaldoCuti     int       `json:"id_saldo_cuti" gorm:"primary_key"`
	TipeAbsenId     string    `json:"tipe_absen_id" gorm:"default:null"`
	Nik             string    `json:"nik" gorm:"default:null"`
	Saldo           int       `json:"saldo" gorm:"default:null"`
	ValidFrom       time.Time `json:"valid_from" gorm:"default:null"`
	ValidTo         time.Time `json:"valid_to" gorm:"default:null"`
	CreatedBy       string    `json:"created_by" gorm:"default:null"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Periode         string    `json:"periode" gorm:"default:null"`
	MaxHutang       int       `json:"max_hutang" gorm:"default:null"`
	ValidFromHutang time.Time `json:"valid_from_hutang" gorm:"default:null"`
}

type HistorySaldoCuti struct {
	IdHistorySaldoCuti int       `json:"id_history_saldo_cuti"`
	TipeAbsenId        string    `json:"tipe_absen_id" gorm:"default:null"`
	Nik                string    `json:"nik" gorm:"default:null"`
	Saldo              int       `json:"saldo" gorm:"default:null"`
	ValidFrom          time.Time `json:"valid_from" gorm:"default:null"`
	ValidTo            time.Time `json:"valid_to" gorm:"default:null"`
	CreatedBy          string    `json:"created_by" gorm:"default:null"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Periode            string    `json:"periode" gorm:"default:null"`
	MaxHutang          int       `json:"max_hutang" gorm:"default:null"`
	ValidFromHutang    time.Time `json:"valid_from_hutang" gorm:"default:null"`
}

type TransaksiCuti struct {
	IdTransaksiCuti  int       `json:"id_transaksi_cuti" gorm:"primary_key"`
	PengajuanAbsenId string    `json:"pengajuan_absen_id" gorm:"default:null"`
	Nik              string    `json:"nik" gorm:"default:null"`
	Periode          string    `json:"periode" gorm:"default:null"`
	JumlahCuti       int       `json:"jumlah_cuti" gorm:"default:null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PengajuanAbsen) TableName() string {
	return "cuti_karyawan.pengajuan_absen"
}
func (HistoryPengajuanAbsen) TableName() string {
	return "cuti_karyawan.history_pengajuan_absen"
}
func (SaldoCuti) TableName() string {
	return "cuti_karyawan.saldo_cuti"
}
func (HistorySaldoCuti) TableName() string {
	return "cuti_karyawan.history_saldo_cuti"
}
func (TipeAbsen) TableName() string {
	return "cuti_karyawan.tipe_absen"
}
func (FileAbsen) TableName() string {
	return "cuti_karyawan.file_absen"
}
func (TransaksiCuti) TableName() string {
	return "cuti_karyawan.transaksi_cuti"
}

type PengajuanAbsenRepo struct {
	DB *gorm.DB
}
type HistoryPengajuanAbsenRepo struct {
	DB *gorm.DB
}
type SaldoCutiRepo struct {
	DB *gorm.DB
}
type HistorySaldoCutiRepo struct {
	DB *gorm.DB
}
type TipeAbsenRepo struct {
	DB *gorm.DB
}
type FileAbsenRepo struct {
	DB *gorm.DB
}
type TransaksiCutiRepo struct {
	DB *gorm.DB
}

func NewPengajuanAbsenRepo(db *gorm.DB) *PengajuanAbsenRepo {
	return &PengajuanAbsenRepo{DB: db}
}
func NewHistoryPengajuanAbsenRepo(db *gorm.DB) *HistoryPengajuanAbsenRepo {
	return &HistoryPengajuanAbsenRepo{DB: db}
}
func NewSaldoCutiRepo(db *gorm.DB) *SaldoCutiRepo {
	return &SaldoCutiRepo{DB: db}
}
func NewHistorySaldoCutiRepo(db *gorm.DB) *HistorySaldoCutiRepo {
	return &HistorySaldoCutiRepo{DB: db}
}
func NewTipeAbsenRepo(db *gorm.DB) *TipeAbsenRepo {
	return &TipeAbsenRepo{DB: db}
}
func NewFileAbsenRepo(db *gorm.DB) *FileAbsenRepo {
	return &FileAbsenRepo{DB: db}
}
func NewTransaksiCutiRepo(db *gorm.DB) *TransaksiCutiRepo {
	return &TransaksiCutiRepo{DB: db}
}

// SALDO CUTI
func (t SaldoCutiRepo) Create(sc SaldoCuti) (SaldoCuti, error) {
	err := t.DB.Create(&sc).Error
	if err != nil {
		return sc, err
	}
	return sc, nil
}

func (t SaldoCutiRepo) Update(sc SaldoCuti) (SaldoCuti, error) {
	err := t.DB.Where("nik=? AND tipe_absen_id=? AND periode=?", sc.Nik, sc.TipeAbsenId, sc.Periode).Save(&sc).Error
	if err != nil {
		return sc, err
	}
	return sc, nil
}

func (t SaldoCutiRepo) DelAdminSaldoCuti(idSaldo int) (SaldoCuti, error) {
	var sc SaldoCuti
	err := t.DB.Where("id_saldo_cuti = ?", idSaldo).First(&sc).Error
	if err == nil {
		t.DB.Where("id_saldo_cuti = ?", idSaldo).Delete(&sc)
		return sc, nil
	}
	return sc, err
}

func (t SaldoCutiRepo) GetSaldoCutiByID(idSaldo int) (SaldoCuti, error) {
	var sc SaldoCuti
	err := t.DB.Where("id_saldo_cuti=?", idSaldo).Take(&sc).Error
	if err != nil {
		return sc, err
	}
	return sc, nil
}

func (t SaldoCutiRepo) FindSaldoCutiKaryawanAdmin(nik string, tahun string) ([]SaldoCuti, error) {
	var sc []SaldoCuti
	err := t.DB.Where("created_by=? AND periode=?", nik, tahun).Find(&sc).Error
	if err != nil {
		return sc, err
	}
	return sc, nil
}

func (t SaldoCutiRepo) GetSaldoCutiPerTipe(id string, nik string, tahun string) ([]SaldoCuti, error) {
	var sc []SaldoCuti
	err := t.DB.Where("tipe_absen_id=? AND nik=? AND periode=?", id, nik, tahun).Take(&sc).Error
	if err != nil {
		return sc, err
	}

	return sc, nil
}

func (t SaldoCutiRepo) FindExistSaldo(tipe_absen_id string, nik string, dateStart string, dateEnd string) (bool, []SaldoCuti, error) {
	var sc_count int64
	var sc []SaldoCuti
	err := t.DB.Table("cuti_karyawan.saldo_cuti").Where("tipe_absen_id=? AND nik=? AND (tsrange(?::date, ?::date, '[]') && tsrange(valid_from, valid_to, '[]'))", tipe_absen_id, nik, dateStart, dateEnd).
		Count(&sc_count).Order("periode ASC").Find(&sc).Error
	if err != nil {
		return false, sc, err
	}
	return sc_count != 0, sc, nil
}

// HISTORY SALDO CUTI
func (t HistorySaldoCutiRepo) Create(hsc HistorySaldoCuti) (HistorySaldoCuti, error) {
	err := t.DB.Create(&hsc).Error
	if err != nil {
		return hsc, err
	}
	return hsc, nil
}

// TIPE CUTI
func (t TipeAbsenRepo) Create(tc TipeAbsen) (TipeAbsen, error) {
	err := t.DB.Create(&tc).Error
	if err != nil {
		return tc, err
	}
	return tc, nil
}

func (t TipeAbsenRepo) Update(tc TipeAbsen) (TipeAbsen, error) {
	err := t.DB.Where("comp_code=?", tc.CompCode).Save(&tc).Error
	if err != nil {
		return tc, err
	}
	return tc, nil

}

func (t TipeAbsenRepo) FindTipeAbsen(compCode string) ([]TipeAbsen, error) {
	var tc []TipeAbsen
	err := t.DB.Where("comp_code=?", compCode).Order("nama_tipe_absen ASC").Find(&tc).Error
	if err != nil {
		return nil, err
	}

	return tc, nil
}

func (t TipeAbsenRepo) FindTipeAbsenByID(id string) (TipeAbsen, error) {
	var tc TipeAbsen
	err := t.DB.Where("id_tipe_absen=?", id).Take(&tc).Error
	if err != nil {
		return tc, err
	}

	return tc, nil
}

// FILE CUTI
func (t FileAbsenRepo) Create(fc FileAbsen) (FileAbsen, error) {
	err := t.DB.Create(&fc).Error
	if err != nil {
		return fc, err
	}
	return fc, nil
}

func (t FileAbsenRepo) Update(fc FileAbsen) (FileAbsen, error) {
	err := t.DB.Where("pengajuan_absen_id=?", fc.PengajuanAbsenId).Save(&fc).Error
	if err != nil {
		return fc, err
	}
	return fc, nil
}

// PENGAJUAN CUTI
func (t PengajuanAbsenRepo) Create(pc PengajuanAbsen) (PengajuanAbsen, error) {
	err := t.DB.Create(&pc).Error
	if err != nil {
		return pc, err
	}
	return pc, nil
}

func (t PengajuanAbsenRepo) Update(pc PengajuanAbsen) (PengajuanAbsen, error) {
	err := t.DB.Where("id_pengajuan_absen=?", pc.IDPengajuanAbsen).Save(&pc).Error
	if err != nil {
		return pc, err
	}
	return pc, nil
}
