package tjsl

import (
	"time"

	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"gorm.io/gorm"
)

type KegiatanKaryawan struct {
	Id                int       `json:"id" gorm:"primary_key"`
	NIK               string    `json:"nik"`
	KegiatanParentId  *int      `json:"kegiatan_parent_id" gorm:"default:null"`
	KoordinatorId     *int      `json:"koordinator_id" gorm:"default:null"`
	NamaKegiatan      string    `json:"nama_kegiatan"`
	TanggalKegiatan   time.Time `json:"tanggal_kegiatan"`
	LokasiKegiatan    string    `json:"lokasi_kegiatan"`
	DeskripsiKegiatan *string   `json:"deskripsi_kegiatan"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Manager           *string   `json:"manager"`
	Slug              string    `json:"slug"`
	DescDecline       *string   `json:"desc_decline" gorm:"default:null"`
	CompCode          string    `json:"comp_code"`
	Periode           string    `json:"periode"`
}

type MyKegiatanTJSL struct {
	Id                int       `json:"id" gorm:"primary_key"`
	NIK               string    `json:"nik"`
	KegiatanParentId  *int      `json:"kegiatan_parent_id" gorm:"default:null"`
	KoordinatorId     *int      `json:"koordinator_id" gorm:"default:null"`
	NamaKegiatan      string    `json:"nama_kegiatan"`
	TanggalKegiatan   string    `json:"tanggal_kegiatan"`
	LokasiKegiatan    string    `json:"lokasi_kegiatan"`
	DeskripsiKegiatan *string   `json:"deskripsi_kegiatan"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Manager           *string   `json:"manager"`
	Slug              *string   `json:"slug"`
	DescDecline       *string   `json:"desc_decline" gorm:"default:null"`
	CompCode          string    `json:"comp_code"`
	Periode           string    `json:"periode"`
}

type DataKegiatanKaryawan struct {
	KegiatanKaryawan
	pihc.PihcMasterKaryRt
}

type RekapLaporanPerbulan struct {
	Bulan          int `json:"bulan"`
	JumlahPerbulan int `json:"jumlah_perbulan"`
	TotalPertahun  int `json:"total_pertahun"`
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

func (t KegiatanKaryawanRepo) FindDataNIKCompCodePeriode(nik_manager string, tahun string, comp_code string, status string) ([]KegiatanKaryawan, error) {
	var kgtn_krywn []KegiatanKaryawan
	err := t.DB.Where("manager=? AND periode=? AND comp_code=? AND status=?", nik_manager, tahun, comp_code, status).Find(&kgtn_krywn).Error
	if err != nil {
		return kgtn_krywn, err
	}
	return kgtn_krywn, nil
}

func (t KegiatanKaryawanRepo) ListKegiatanKaryawanApprvalWait(nik_manager string, tahun string, comp_code string, status string) ([]DataKegiatanKaryawan, error) {
	results := []DataKegiatanKaryawan{}
	err := t.DB.Raw(`select * 
							from tjsl.kegiatan_karyawan kk 
						join dbo.pihc_master_kary_rt pmkr on pmkr.emp_no = kk.nik
					where manager = ? and comp_code=? and periode =? and status = ?`, nik_manager, comp_code, tahun, status).
		Scan(&results).Error

	if err != nil {
		return results, err
	}

	return results, nil
}

func (t KegiatanKaryawanRepo) DelKegiatanKaryawanID(slug string) error {
	var data []KegiatanKaryawan
	err := t.DB.Where("slug = ?", slug).First(&data).Error
	if err == nil {
		t.DB.Where("slug = ?", slug).Delete(&data)
		return nil
	}
	return err
}

// func (t KegiatanKaryawanRepo) RekapPerbulan(nik string, periode string, status string, bulan int) (int, error) {
// 	var jmlahPerbulan int

// 	err := t.DB.Raw(`
// 	SELECT COUNT(*) AS jumlah_perbulan
// 		FROM tjsl.kegiatan_karyawan kk
// 	WHERE nik = ? AND periode = ? AND status = ? AND EXTRACT(MONTH FROM tanggal_kegiatan) = ?`, nik, periode, status, bulan).Scan(&jmlahPerbulan).Error

// 	if err != nil {
// 		return jmlahPerbulan, err
// 	}

// 	return jmlahPerbulan, nil
// }

func (t KegiatanKaryawanRepo) RekapPerbulan(nik string, periode string, status string) ([]RekapLaporanPerbulan, error) {
	var data []RekapLaporanPerbulan

	err := t.DB.Raw(`
	SELECT NULL AS bulan, NULL AS jumlah_perbulan,count(*) AS total_pertahun
		FROM tjsl.kegiatan_karyawan kk
	WHERE tanggal_kegiatan >= CURRENT_DATE - INTERVAL '1' year AND status = ? AND nik = ?
	UNION ALL
	SELECT EXTRACT(MONTH FROM tanggal_kegiatan) AS bulan, COUNT(*) AS jumlah_perbulan, NULL AS total_pertahun
		FROM tjsl.kegiatan_karyawan kk
	WHERE periode = ? AND nik = ? AND status = ?
	GROUP BY EXTRACT(MONTH FROM tanggal_kegiatan)
	ORDER BY bulan ASC`, status, nik, periode, nik, status).Scan(&data).Error

	if err != nil {
		return data, err
	}

	return data, nil
}
