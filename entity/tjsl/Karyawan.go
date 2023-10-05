package tjsl

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"google.golang.org/api/iterator"
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
type RekapLeaderBoardDB struct {
	EmpNama   string `json:"emp_nama"`
	Nik       string `json:"nik"`
	Total     int    `json:"total"`
	Company   string `json:"company"`
	PosID     string `json:"pos_id"`
	PosTitle  string `json:"pos_title"`
	DeptID    string `json:"dept_id"`
	DeptTitle string `json:"dept_title"`
	KompID    string `json:"komp_id"`
	KompTitle string `json:"komp_title"`
	DirID     string `json:"dir_id"`
	DirTitle  string `json:"dir_title"`
}

type RekapLeaderBoard struct {
	EmpNama   string  `json:"emp_nama"`
	Nik       string  `json:"nik"`
	Total     int     `json:"total"`
	Company   string  `json:"company"`
	PosID     *string `json:"pos_id"`
	PosTitle  *string `json:"pos_title"`
	DeptID    *string `json:"dept_id"`
	DeptTitle *string `json:"dept_title"`
	KompID    *string `json:"komp_id"`
	KompTitle *string `json:"komp_title"`
	DirID     *string `json:"dir_id"`
	DirTitle  *string `json:"dir_title"`
	Photo     string  `json:"photo"`
}

func (KegiatanKaryawan) TableName() string {
	return "tjsl.kegiatan_karyawan"
}

type KegiatanKaryawanRepo struct {
	DB            *gorm.DB
	StorageClient *storage.Client
}

func NewKegiatanKaryawanRepo(db *gorm.DB, sc *storage.Client) *KegiatanKaryawanRepo {
	return &KegiatanKaryawanRepo{DB: db, StorageClient: sc}
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
	SELECT NULL AS bulan, NULL AS jumlah_perbulan,count(kk.status) AS total_pertahun
		FROM tjsl.kegiatan_karyawan kk
	WHERE tanggal_kegiatan >= CURRENT_DATE - INTERVAL '1' year AND status = ? AND nik = ?
	UNION ALL
	SELECT EXTRACT(MONTH FROM tanggal_kegiatan) AS bulan, COUNT(kk.status) AS jumlah_perbulan, NULL AS total_pertahun
		FROM tjsl.kegiatan_karyawan kk
	WHERE periode = ? AND nik = ? AND status = ?
	GROUP BY EXTRACT(MONTH FROM tanggal_kegiatan)
	ORDER BY bulan ASC`, status, nik, periode, nik, status).Scan(&data).Error

	if err != nil {
		return data, err
	}

	return data, nil
}

func (t KegiatanKaryawanRepo) FindLeaderBoardStatusPeriodeCompany(periode string, status string, company string) ([]RekapLeaderBoard, error) {
	var data []RekapLeaderBoardDB

	err := t.DB.Raw(`
	SELECT
		pmkr.nama as emp_nama, pmkr.emp_no as nik, COUNT(kk.status) AS total,
		pmkr.dept_title, pmkr.company as company,
		pmkr.pos_id as pos_id, pmkr.pos_title as pos_title,
		pmkr.dept_id as dept_id, pmkr.dept_title as dept_title, pmkr.komp_id as komp_id,
		pmkr.komp_title as komp_title, pmkr.dir_id as dir_id, pmkr.dir_title as dir_title
	from dbo.pihc_master_kary_rt pmkr
	left join tjsl.kegiatan_karyawan kk ON pmkr.emp_no = kk.nik
	where pmkr.company = ? and 
		(kk.status = ? or kk.status is null) and 
		(kk.periode = ? or kk.periode is null)
	GROUP BY
		pmkr.emp_no, pmkr.nama, pmkr.dept_title, pmkr.company, pmkr.pos_id,
		pmkr.pos_title, pmkr.dept_id, pmkr.dept_title, pmkr.komp_id,
		pmkr.komp_title, pmkr.dir_id, pmkr.dir_title
	ORDER BY
		total desc
  	limit 10`, company, status, periode).Scan(&data).Error

	var result []RekapLeaderBoard
	var namaFile string
	for _, datas := range data {
		files, err := t.FindPhotosKaryawan(datas.Nik, datas.Company)
		if err != nil {
			namaFile = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
		} else {
			namaFile = "https://storage.googleapis.com/" + files
		}
		data_rekap_leader_board_convert := convertRekapLeaderBoard(datas, namaFile)
		result = append(result, data_rekap_leader_board_convert)
	}
	if err != nil {
		return result, err
	}

	return result, nil
}

func (t KegiatanKaryawanRepo) FindPhotosKaryawan(objName string, company string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	ctx := context.Background() // Create a new context

	// bckt, err := t.StorageClient.Bucket(os.Getenv("GC_LUMEN_BUCKET")).Attrs(ctx)
	// if err != nil {
	// 	return "", fmt.Errorf("Bucket(%q).Attrs: %w", os.Getenv("GC_LUMEN_BUCKET"), err)
	// }

	// fmt.Println(bckt.Name)
	// fmt.Println(bckt.Location)
	// fmt.Println(bckt.LocationType)
	// fmt.Println(bckt.StorageClass)
	// fmt.Println(bckt.RPO)
	// fmt.Println(bckt.Created)
	// fmt.Println(bckt.MetaGeneration)
	// fmt.Println(bckt.PredefinedACL)

	buckets := os.Getenv("GC_LUMEN_BUCKET")
	bckt := t.StorageClient.Bucket(buckets)
	location := "DataKaryawan/Foto/" + company + "/" + objName
	query := &storage.Query{Prefix: location}

	objectIterator := bckt.Objects(ctx, query)
	found := false

	for {
		attrs, err := objectIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error iterating over objects: %v", err)
		}

		if location+filepath.Ext(attrs.Name) == attrs.Name {
			location = buckets + "/" + location + filepath.Ext(attrs.Name)
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("object not found")
	}

	return location, nil
}

func convertRekapLeaderBoard(source RekapLeaderBoardDB, files string) RekapLeaderBoard {
	var results RekapLeaderBoard
	results.EmpNama = source.EmpNama
	results.Nik = source.Nik
	results.Total = source.Total
	results.Company = source.Company

	if source.PosID != "" {
		results.PosID = &source.PosID
	}
	if source.PosTitle != "" {
		results.PosTitle = &source.PosTitle
	}
	if source.DeptID != "" {
		results.DeptID = &source.DeptID
	}
	if source.DeptTitle != "" {
		results.DeptTitle = &source.DeptTitle
	}
	if source.KompID != "" {
		results.KompID = &source.KompID
	}
	if source.KompTitle != "" {
		results.KompTitle = &source.KompTitle
	}
	if source.DirID != "" {
		results.DirID = &source.DirID
	}
	if source.DirTitle != "" {
		results.DirTitle = &source.DirTitle
	}
	results.Photo = files
	return results
}
