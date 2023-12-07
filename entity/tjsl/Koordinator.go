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
		SELECT kk.id_koordinator, kk.kegiatan_parent_id, kk.nama, kk.created_by, kk.created_at, kk.updated_at, kk.comp_code, kk.slug, kk.periode
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

	var karyawan []pihc.PihcMasterKaryRtDb
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
		data_karyawan_convert := convertSourceTargetDataKaryawanRt(karyawan[i])
		data.Employee = data_karyawan_convert
		results[i].Employee = data.Employee
	}

	return results, nil
}

func convertSourceTargetDataKaryawanRt(source pihc.PihcMasterKaryRtDb) pihc.PihcMasterKaryRt {
	return pihc.PihcMasterKaryRt{
		EmpNo:          source.EmpNo,
		Nama:           source.Nama,
		Gender:         source.Gender,
		Agama:          source.Agama,
		StatusKawin:    source.StatusKawin,
		Anak:           source.Anak,
		Mdg:            "0",
		EmpGrade:       source.EmpGrade,
		EmpGradeTitle:  source.EmpGradeTitle,
		Area:           source.Area,
		AreaTitle:      source.AreaTitle,
		SubArea:        source.SubArea,
		SubAreaTtitle:  source.SubAreaTtitle,
		Contract:       source.Contract,
		Pendidikan:     source.Pendidikan,
		Company:        source.Company,
		Lokasi:         source.Lokasi,
		EmployeeStatus: source.EmployeeStatus,
		Email:          source.Email,
		HP:             source.HP,
		TglLahir:       source.TglLahir.Format("2006-01-02"),
		PosID:          source.PosID,
		PosTitle:       source.PosTitle,
		SubPosID:       source.SupPosID,
		PosGrade:       source.PosGrade,
		PosKategori:    source.PosKategori,
		OrgID:          source.OrgID,
		OrgTitle:       source.OrgTitle,
		DeptID:         source.DeptID,
		DeptTitle:      source.DeptTitle,
		KompID:         source.KompID,
		KompTitle:      source.KompTitle,
		DirID:          source.DirID,
		DirTitle:       source.DirTitle,
		PosLevel:       source.PosLevel,
		SupEmpNo:       source.SupEmpNo,
		BagID:          source.BagID,
		BagTitle:       source.BagTitle,
		SeksiID:        source.SeksiID,
		SeksiTitle:     source.SeksiTitle,
		PreNameTitle:   source.PreNameTitle,
		PostNameTitle:  source.PostNameTitle,
		NoNPWP:         source.NoNPWP,
		BankAccount:    source.BankAccount,
		BankName:       source.BankName,
		MdgDate:        source.MdgDate,
		PayScale:       source.PayScale,
		CCCode:         source.CCCode,
		Nickname:       source.Nickname,
	}
}
func (t KegiatanKoordinatorRepo) ListKoordinatorDalamKegiatan(slug string, nik string) ([]Result, error) {
	results := []Result{}
	err := t.DB.Raw(`SELECT id_koordinator,kk.kegiatan_parent_id,kk.nama,kk.created_by ,kk.created_at ,kk.updated_at ,kk.comp_code ,kk.slug,kk.periode 
							FROM tjsl.kegiatan_mstr km 
						JOIN tjsl.kegiatan_koordinator kk on kk.kegiatan_parent_id = km.id_kegiatan
						JOIN tjsl.koordinator_person kp on kp.koordinator_id = kk.id_koordinator 
						JOIN dbo.pihc_master_kary_rt pmkr on kk.created_by = pmkr.emp_no
						WHERE km.slug = ? and kp.nik = ?
					ORDER BY nama asc`, slug, nik).
		Scan(&results).Error

	if err != nil {
		return results, err
	}

	var karyawan []pihc.PihcMasterKaryRtDb
	t.DB.Raw(`SELECT pmkr.*
				FROM tjsl.kegiatan_mstr km
			JOIN tjsl.kegiatan_koordinator kk ON kk.kegiatan_parent_id = km.id_kegiatan
			JOIN dbo.pihc_master_kary_rt pmkr ON kk.created_by = pmkr.emp_no
			JOIN tjsl.koordinator_person kp on kp.koordinator_id = kk.id_koordinator 
				WHERE km.slug = ? and kp.nik = ?
			ORDER BY nama ASC`, slug, nik).
		Scan(&karyawan)

	for i, data := range results {
		data_karyawan_convert := convertSourceTargetDataKaryawanRt(karyawan[i])
		data.Employee = data_karyawan_convert
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
