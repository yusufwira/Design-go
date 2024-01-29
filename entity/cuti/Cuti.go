package cuti

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PengajuanAbsen struct {
	IdPengajuanAbsen int             `json:"id_pengajuan_absen" gorm:"primary_key"`
	Nik              string          `json:"nik" gorm:"default:null"`
	CompCode         string          `json:"comp_code" gorm:"default:null"`
	TipeAbsenId      *string         `json:"tipe_absen_id" gorm:"default:null"`
	Deskripsi        *string         `json:"deskripsi" gorm:"default:null"`
	MulaiAbsen       time.Time       `json:"mulai_absen" gorm:"default:null"`
	AkhirAbsen       time.Time       `json:"akhir_absen" gorm:"default:null"`
	TglPengajuan     time.Time       `json:"tgl_pengajuan" gorm:"default:null"`
	Status           *string         `json:"status" gorm:"default:null"`
	CreatedBy        *string         `json:"created_by" gorm:"default:null"`
	CreatedAt        time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	Keterangan       *string         `json:"keterangan" gorm:"default:null"`
	Periode          *string         `json:"periode" gorm:"default:null"`
	ApprovedBy       json.RawMessage `json:"approved_by" gorm:"default:null;type:json"`
	JmlHariKalendar  *int            `json:"jml_hari_kalendar" gorm:"default:null"`
	JmlHariKerja     *int            `json:"jml_hari_kerja" gorm:"default:null"`
}

type AtasanApproved struct {
	Nik          string  `json:"nik"`
	Name         string  `json:"name"`
	Position     string  `json:"position"`
	TypeApprover *string `json:"type_approver"`
	Status       *string `json:"status"`
	Keterangan   *string `json:"keterangan"`
	Photo        string  `json:"photo"`
}

type MyPengajuanAbsen struct {
	IdPengajuanAbsen int    `json:"id_pengajuan_absen"`
	Nik              string `json:"nik"`
	CompCode         string `json:"comp_code"`
	TipeAbsen        `json:"tipe_absen"`
	Deskripsi        *string         `json:"deskripsi"`
	MulaiAbsen       string          `json:"mulai_absen"`
	AkhirAbsen       string          `json:"akhir_absen"`
	TglPengajuan     string          `json:"tgl_pengajuan"`
	Status           *string         `json:"status"`
	CreatedBy        *string         `json:"created_by"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	Keterangan       *string         `json:"keterangan"`
	Periode          *string         `json:"periode"`
	ApprovedBy       json.RawMessage `json:"approved_by" gorm:"default:null;type:json"`
	JmlHariKalendar  *int            `json:"jml_hari_kalendar"`
	JmlHariKerja     *int            `json:"jml_hari_kerja"`
}

type HistoryPengajuanAbsen struct {
	IDHistoryPengajuanAbsen int             `json:"id_history_pengajuan_absen"`
	Nik                     string          `json:"nik" gorm:"default:null"`
	CompCode                string          `json:"comp_code" gorm:"default:null"`
	TipeAbsenId             *string         `json:"tipe_absen_id" gorm:"default:null"`
	Deskripsi               *string         `json:"deskripsi" gorm:"default:null"`
	MulaiAbsen              time.Time       `json:"mulai_absen" gorm:"default:null"`
	AkhirAbsen              time.Time       `json:"akhir_absen" gorm:"default:null"`
	TglPengajuan            time.Time       `json:"tgl_pengajuan" gorm:"default:null"`
	Status                  *string         `json:"status" gorm:"default:null"`
	CreatedBy               *string         `json:"created_by" gorm:"default:null"`
	CreatedAt               time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt               time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	Keterangan              *string         `json:"keterangan" gorm:"default:null"`
	Periode                 *string         `json:"periode" gorm:"default:null"`
	ApprovedBy              json.RawMessage `json:"approved_by" gorm:"default:null;type:json"`
	JmlHariKalendar         *int            `json:"jml_hari_kalendar" gorm:"default:null"`
	JmlHariKerja            *int            `json:"jml_hari_kerja" gorm:"default:null"`
}

type FileAbsen struct {
	IdFileAbsen      int       `json:"id_file_absen" gorm:"primary_key"`
	PengajuanAbsenId int       `json:"pengajuan_absen_id" gorm:"default:null"`
	Filename         *string   `json:"filename" gorm:"default:null"`
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
	Saldo           int       `json:"saldo" gorm:"default:0"`
	ValidFrom       time.Time `json:"valid_from" gorm:"default:null"`
	ValidTo         time.Time `json:"valid_to" gorm:"default:null"`
	CreatedBy       string    `json:"created_by" gorm:"default:null"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Periode         string    `json:"periode" gorm:"default:null"`
	MaxHutang       int       `json:"max_hutang" gorm:"default:0"`
	ValidFromHutang time.Time `json:"valid_from_hutang" gorm:"default:null"`
}

type HistorySaldoCuti struct {
	IdHistorySaldoCuti int       `json:"id_history_saldo_cuti"`
	TipeAbsenId        string    `json:"tipe_absen_id" gorm:"default:null"`
	Nik                string    `json:"nik" gorm:"default:null"`
	Saldo              int       `json:"saldo" gorm:"default:0"`
	ValidFrom          time.Time `json:"valid_from" gorm:"default:null"`
	ValidTo            time.Time `json:"valid_to" gorm:"default:null"`
	CreatedBy          string    `json:"created_by" gorm:"default:null"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Periode            string    `json:"periode" gorm:"default:null"`
	MaxHutang          int       `json:"max_hutang" gorm:"default:0"`
	ValidFromHutang    time.Time `json:"valid_from_hutang" gorm:"default:null"`
}

type TransaksiCuti struct {
	IdTransaksiCuti  int       `json:"id_transaksi_cuti" gorm:"primary_key"`
	PengajuanAbsenId int       `json:"pengajuan_absen_id" gorm:"default:null"`
	TipeAbsenId      string    `json:"tipe_absen_id" gorm:"default:notnull"`
	Nik              string    `json:"nik" gorm:"default:null"`
	Periode          string    `json:"periode" gorm:"default:null"`
	TipeHari         string    `json:"tipe_hari" gorm:"default:null"`
	JumlahCuti       int       `json:"jumlah_cuti" gorm:"default:null"`
	Keterangan       *string   `json:"keterangan" gorm:"default:null"`
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
	err := t.DB.Where("id_saldo_cuti = ?", sc.IdSaldoCuti).Save(&sc).Error
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

func (t SaldoCutiRepo) GetSaldoCutiByID(idSaldo interface{}) (SaldoCuti, error) {
	var sc SaldoCuti
	err := t.DB.Where("id_saldo_cuti=?", idSaldo).Take(&sc).Error
	if err != nil {
		return sc, err
	}
	return sc, nil
}

func (t SaldoCutiRepo) FindSaldoCutiKaryawanAdmin(key string, company string, direktorat string, departemen string, kompartemen string, nik string, tahun string) ([]SaldoCuti, error) {
	var sc []SaldoCuti

	query := t.DB.
		Select("cuti_karyawan.saldo_cuti.*").
		Joins("inner join dbo.pihc_master_kary_rt pmkr on cuti_karyawan.saldo_cuti.nik = pmkr.emp_no")

	if key != "" {
		query = query.Where("(pmkr.emp_no like ? OR lower(pmkr.nama) like lower(?) OR pmkr.pos_id like ? OR lower(pmkr.pos_title) like lower(?))",
			key, "%"+key+"%", key, "%"+key+"%")
	}

	if company != "" {
		query = query.Where("pmkr.company = ?", company)
	}

	if direktorat != "" {
		query = query.Where("pmkr.dir_id = ?", direktorat)
	}

	if kompartemen != "" {
		query = query.Where("pmkr.komp_id = ?", kompartemen)
	}

	if departemen != "" {
		query = query.Where("pmkr.dept_id = ?", departemen)
	}

	query = query.Where("cuti_karyawan.saldo_cuti.created_by = ? AND cuti_karyawan.saldo_cuti.periode = ?", nik, tahun)

	err := query.Order("cuti_karyawan.saldo_cuti.valid_from ASC").Find(&sc).Error
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func (t SaldoCutiRepo) GetSaldoCutiPerTipeArr(nik string, tipe string, tahun string) ([]SaldoCuti, error) {
	var sc []SaldoCuti

	time_periode_start, _ := time.Parse(time.DateOnly, tahun+"-01-01")
	time_periode_end, _ := time.Parse(time.DateOnly, tahun+"-12-31")
	fmt.Println(time_periode_start, time_periode_end)
	err := t.DB.Where("tipe_absen_id=? AND nik=? AND (tsrange(valid_from, valid_to) && tsrange(?::date, ?::date))",
		tipe, nik, time_periode_start, time_periode_end).Order("periode asc").
		Find(&sc).Error

	fmt.Println(sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func (t SaldoCutiRepo) GetSaldoCutiPerTipe(nik string, tipe string, tahun string) (SaldoCuti, error) {
	var sc SaldoCuti
	err := t.DB.Where("tipe_absen_id=? AND nik=? AND periode=?", tipe, nik, tahun).Take(&sc).Error
	if err != nil {
		return sc, err
	}
	return sc, nil
}

// HISTORY SALDO CUTI
func (t HistorySaldoCutiRepo) Create(hsc HistorySaldoCuti) (HistorySaldoCuti, error) {
	err := t.DB.Create(&hsc).Error
	if err != nil {
		return hsc, err
	}
	return hsc, nil
}

// HISTORY PENGAJUAN CUTI
func (t HistoryPengajuanAbsenRepo) Create(hsp HistoryPengajuanAbsen) (HistoryPengajuanAbsen, error) {
	err := t.DB.Create(&hsp).Error
	if err != nil {
		return hsp, err
	}
	return hsp, nil
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

func (t TipeAbsenRepo) FindTipeAbsenSaldo(compCode string) ([]TipeAbsen, error) {
	var tc []TipeAbsen
	err := t.DB.Where("comp_code=? AND (max_absen is null or max_absen = 0)", compCode).Order("nama_tipe_absen ASC").Find(&tc).Error
	if err != nil {
		return nil, err
	}

	return tc, nil
}
func (t TipeAbsenRepo) FindTipeAbsenPengajuan(compCode string) ([]TipeAbsen, error) {
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
func (t TipeAbsenRepo) FindTipeAbsenByIDArray(id []string) ([]TipeAbsen, error) {
	var tc []TipeAbsen
	err := t.DB.Where("id_tipe_absen in(?)", id).Find(&tc).Error
	if err != nil {
		return tc, err
	}

	return tc, nil
}

// FILE CUTI
func (t FileAbsenRepo) CreateArr(fc []FileAbsen) ([]FileAbsen, error) {
	err := t.DB.Create(&fc).Error
	if err != nil {
		return fc, err
	}
	return fc, nil
}
func (t FileAbsenRepo) Delete(fc FileAbsen) (FileAbsen, error) {
	err := t.DB.Delete(&fc).Error
	if err != nil {
		return fc, err
	}
	return fc, nil
}

func (t FileAbsenRepo) FindFileAbsenByIDPengajuan(id_pengajuan int) ([]FileAbsen, error) {
	var fc []FileAbsen
	err := t.DB.Where("pengajuan_absen_id=?", id_pengajuan).Find(&fc).Error
	if err != nil {
		return fc, err
	}
	return fc, nil
}
func (t FileAbsenRepo) FindFileAbsenByIDPengajuanArray(id_pengajuan []int) ([]FileAbsen, error) {
	var fc []FileAbsen
	err := t.DB.Where("pengajuan_absen_id in(?)", id_pengajuan).Find(&fc).Error
	if err != nil {
		return fc, err
	}
	return fc, nil
}

// PENGAJUAN ABSEN
func (t PengajuanAbsenRepo) Create(pc PengajuanAbsen) (PengajuanAbsen, error) {
	err := t.DB.Create(&pc).Error
	if err != nil {
		return pc, err
	}
	return pc, nil
}

func (t PengajuanAbsenRepo) Update(pc PengajuanAbsen) (PengajuanAbsen, error) {
	err := t.DB.Where("id_pengajuan_absen=?", pc.IdPengajuanAbsen).Save(&pc).Error
	if err != nil {
		return pc, err
	}
	return pc, nil
}

func (t PengajuanAbsenRepo) FindDataNIKPeriode(nik string, tahun string) ([]PengajuanAbsen, error) {
	var pengajuan_absen []PengajuanAbsen
	err := t.DB.Where("nik=? AND periode=?", nik, tahun).Find(&pengajuan_absen).Error
	if err != nil {
		return pengajuan_absen, err
	}
	return pengajuan_absen, nil
}
func (t PengajuanAbsenRepo) FindDataIdPengajuan(id interface{}) (PengajuanAbsen, error) {
	var pengajuan_absen PengajuanAbsen
	err := t.DB.Where("id_pengajuan_absen=?", id).Take(&pengajuan_absen).Error
	if err != nil {
		fmt.Println("ERR")
		return pengajuan_absen, err
	}
	return pengajuan_absen, nil
}

func (t PengajuanAbsenRepo) FindDataNIKPeriodeApproval(nik string, tahun string, manager bool) ([]PengajuanAbsen, error) {
	var pengajuan_absen []PengajuanAbsen
	var err error

	query := t.DB.
		Table("cuti_karyawan.pengajuan_absen").
		Select("*").
		Where("periode = ?", tahun)

	if manager {
		query = query.Where("EXISTS (SELECT 1 FROM json_array_elements(approved_by) AS approver WHERE approver->>'status' IS NOT NULL AND approver->>'nik' = ?)", nik)
	} else {
		query = query.Where("EXISTS (SELECT 1 FROM json_array_elements(approved_by) AS approver WHERE approver->>'status' IS NOT NULL)")
	}

	err = query.Find(&pengajuan_absen).Error
	if err != nil {
		return pengajuan_absen, err
	}
	return pengajuan_absen, nil
}
func (t PengajuanAbsenRepo) FindDataNIKPeriodeApprovalWaiting(nik string, tahun string, status string, manager bool) ([]PengajuanAbsen, error) {
	var pengajuan_absen []PengajuanAbsen
	var err error

	query := t.DB.
		Table("cuti_karyawan.pengajuan_absen").
		Select("*").
		Where("periode = ? AND status = ?", tahun, status)
	if manager {
		query = query.Where("EXISTS (SELECT 1 FROM json_array_elements(approved_by) AS approver WHERE approver->>'status' = 'WaitApv' and approver->>'nik' = ?)", nik)
	} else {
		query = query.Where("EXISTS (SELECT 1 FROM json_array_elements(approved_by) AS approver WHERE approver->>'status' = 'WaitApv')")
	}

	err = query.Find(&pengajuan_absen).Error
	if err != nil {
		return pengajuan_absen, err
	}

	return pengajuan_absen, nil
}

func (t PengajuanAbsenRepo) DelPengajuanCuti(id int) (PengajuanAbsen, error) {
	var sc PengajuanAbsen
	err := t.DB.Where("id_pengajuan_absen = ?", id).First(&sc).Error
	if err == nil {
		t.DB.Where("id_pengajuan_absen = ?", id).Delete(&sc)
		return sc, nil
	}
	return sc, err
}

// Transaksi Cuti
func (t TransaksiCutiRepo) FindDataTransaksiIDPengajuan(id int) ([]TransaksiCuti, error) {
	var traksaksi_cuti []TransaksiCuti
	err := t.DB.Where("pengajuan_absen_id=?", id).Find(&traksaksi_cuti).Error
	if err != nil {
		return traksaksi_cuti, err
	}
	return traksaksi_cuti, nil
}

func (t TransaksiCutiRepo) Create(tc TransaksiCuti) (TransaksiCuti, error) {
	err := t.DB.Create(&tc).Error
	if err != nil {
		return tc, err
	}
	return tc, nil
}

func (t TransaksiCutiRepo) Delete(tc TransaksiCuti) (TransaksiCuti, error) {
	err := t.DB.Delete(&tc).Error
	if err != nil {
		return tc, err
	}
	return tc, nil
}
