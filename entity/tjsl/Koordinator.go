package tjsl

import (
	"time"

	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"gorm.io/gorm"
)

type KegiatanKoordinator struct {
	IdKoordinator    int       `json:"id_koordinator" gorm:"primary_key"`
	KegiatanParentId *int      `json:"kegiatan_parent_id" gorm:"default:null"`
	Nama             string    `json:"nama"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CompCode         string    `json:"comp_code"`
	Slug             string    `json:"slug"`
	Periode          string    `json:"periode"`
}

type Result struct {
	KegiatanKoordinator
	Employee pihc.PihcMasterKaryRt `json:"employee" gorm:"foreignkey:EmpNo;association_foreignkey:CreatedBy"`
}

func (KegiatanKoordinator) TableName() string {
	return "tjsl.kegiatan_koordinator"
}

type KegiatanKoordinatorRepo struct {
	DB *gorm.DB
}

func NewKegiatanKoordinatorRepo(db *gorm.DB) *KegiatanKoordinatorRepo {
	return &KegiatanKoordinatorRepo{DB: db}
}

func (t KegiatanKoordinatorRepo) Create(koor_kgt KegiatanKoordinator) (KegiatanKoordinator, error) {
	err := t.DB.Create(&koor_kgt).Error
	if err != nil {
		return koor_kgt, err
	}
	return koor_kgt, nil
}

func (t KegiatanKoordinatorRepo) Update(koor_kgt KegiatanKoordinator) (KegiatanKoordinator, error) {
	err := t.DB.Save(&koor_kgt).Error
	if err != nil {
		return koor_kgt, err
	}
	return koor_kgt, nil
}

func (t KegiatanKoordinatorRepo) FindDataParentID(id int, nik string) ([]KegiatanKoordinator, error) {
	var koor_kgt []KegiatanKoordinator
	err := t.DB.Where("kegiatan_parent_id=? AND created_by=?", id, nik).Order("id_koordinator ASC").Find(&koor_kgt).Error
	if err != nil {
		return koor_kgt, err
	}
	return koor_kgt, nil
}

func (t KegiatanKoordinatorRepo) FindDataID(id int) (KegiatanKoordinator, error) {
	var koor_kgt KegiatanKoordinator
	err := t.DB.Where("id_koordinator=?", id).First(&koor_kgt).Error
	if err != nil {
		return koor_kgt, err
	}
	return koor_kgt, nil
}

func (t KegiatanKoordinatorRepo) FindDataKoorIDLuarKegiatan(nik string) ([]KegiatanKoordinator, error) {
	var koor_kgt []KegiatanKoordinator

	err := t.DB.Where("id_koordinator IN (select distinct koordinator_id from tjsl.koordinator_person where nik=?) AND kegiatan_parent_id IS NULL", nik).
		Order("id_koordinator ASC").Find(&koor_kgt).Error
	if err != nil {
		//Joins("inner join dbo.pihc_master_kary_rt as b on tjsl.kegiatan_koordinator.created_by = b.emp_no").
		return koor_kgt, err
	}
	return koor_kgt, nil
}

func (t KegiatanKoordinatorRepo) ListKoordinatorLuarKegiatan(nik string) ([]Result, error) {
	results := []Result{}
	err := t.DB.Raw(`
		SELECT kk.id_koordinator, kk.nama, kk.created_by, kk.created_at, kk.updated_at, kk.comp_code, kk.slug, kk.periode
		FROM tjsl.kegiatan_koordinator kk
		JOIN dbo.pihc_master_kary_rt pmkr ON kk.created_by = pmkr.emp_no
		WHERE kk.id_koordinator IN (
			SELECT DISTINCT koordinator_id
			FROM tjsl.koordinator_person
			WHERE nik = ?
		) AND kk.kegiatan_parent_id IS NULL
		ORDER BY kk.id_koordinator ASC
	`, nik).Scan(&results).Error

	if err != nil {
		return results, err
	}

	var karyawan []pihc.PihcMasterKaryRt
	t.DB.Raw(`
		SELECT pmkr.*
		FROM tjsl.kegiatan_koordinator kk
		JOIN dbo.pihc_master_kary_rt pmkr ON kk.created_by = pmkr.emp_no
		WHERE kk.id_koordinator IN (
			SELECT DISTINCT koordinator_id
			FROM tjsl.koordinator_person
			WHERE nik = ?
		) AND kk.kegiatan_parent_id IS NULL
		ORDER BY kk.id_koordinator ASC
	`, nik).Scan(&karyawan)

	for i, data := range results {
		data.Employee = karyawan[i]
		results[i].Employee = data.Employee
	}

	return results, nil
}
func (t KegiatanKoordinatorRepo) ListKoordinatorDalamKegiatan(slug string, nik string) ([]Result, error) {
	results := []Result{}
	err := t.DB.Raw(`SELECT kk.id_koordinator, kk.nama, kk.created_by, kk.created_at, kk.updated_at, kk.comp_code, kk.slug, kk.periode
						FROM tjsl.kegiatan_mstr km
					JOIN tjsl.kegiatan_koordinator kk ON kk.kegiatan_parent_id = km.id_kegiatan
					JOIN dbo.pihc_master_kary_rt pmkr ON kk.created_by = pmkr.emp_no
						WHERE km.slug = ? AND kk.created_by = ?
					ORDER BY kk.id_koordinator ASC`, slug, nik).
		Scan(&results).Error

	if err != nil {
		return results, err
	}

	var karyawan []pihc.PihcMasterKaryRt
	t.DB.Raw(`SELECT pmkr.*
				FROM tjsl.kegiatan_mstr km
			JOIN tjsl.kegiatan_koordinator kk ON kk.kegiatan_parent_id = km.id_kegiatan
			JOIN dbo.pihc_master_kary_rt pmkr ON kk.created_by = pmkr.emp_no
				WHERE km.slug = ? AND kk.created_by = ?
			ORDER BY kk.id_koordinator ASC`, slug, nik).
		Scan(&karyawan)

	for i, data := range results {
		data.Employee = karyawan[i]
		results[i].Employee = data.Employee
	}

	return results, nil
}

func (t KegiatanKoordinatorRepo) FindDataSlug(slug string) (KegiatanKoordinator, error) {
	var koor_kgt KegiatanKoordinator
	err := t.DB.Where("slug=?", slug).First(&koor_kgt).Error
	if err != nil {
		return koor_kgt, err
	}
	return koor_kgt, nil
}

func (t KegiatanKoordinatorRepo) DelKegiatanKoordinatorID(slug string) error {
	var data []KegiatanKoordinator
	err := t.DB.Where("slug = ?", slug).First(&data).Error
	if err == nil {
		t.DB.Where("slug = ?", slug).Delete(&data)
		return nil
	}
	return err
}
