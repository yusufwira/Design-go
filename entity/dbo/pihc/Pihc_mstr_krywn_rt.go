package pihc

import (
	"fmt"
	"time"

	"github.com/yusufwira/lern-golang-gin/entity/mobile/profile"
	"github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type PihcMasterKaryRtDb struct {
	EmpNo          string    `json:"emp_no" gorm:"primary_key"`
	Nama           string    `json:"nama"`
	Gender         string    `json:"gender"`
	Agama          string    `json:"agama"`
	StatusKawin    string    `json:"status_kawin"`
	Anak           int8      `json:"anak"` //
	Mdg            int       `json:"mdg"`  //
	EmpGrade       string    `json:"emp_grade"`
	EmpGradeTitle  string    `json:"emp_grade_title"`
	Area           string    `json:"area"`
	AreaTitle      string    `json:"area_title"`
	SubArea        string    `json:"sub_area"`
	SubAreaTitle   string    `json:"sub_area_title"`
	Contract       string    `json:"contract"`
	Pendidikan     string    `json:"pendidikan"`
	Company        string    `json:"company"`
	Lokasi         string    `json:"lokasi"`
	EmployeeStatus string    `json:"employee_status"`
	Email          string    `json:"email"`
	HP             string    `json:"hp"`
	TglLahir       time.Time `json:"tgl_lahir"`
	PosID          string    `json:"pos_id"`
	PosTitle       string    `json:"pos_title"`
	SupPosID       string    `json:"sup_pos_id"`
	PosGrade       string    `json:"pos_grade"`
	PosKategori    string    `json:"pos_kategori"`
	OrgID          string    `json:"org_id"`
	OrgTitle       string    `json:"org_title"`
	DeptID         string    `json:"dept_id"`
	DeptTitle      string    `json:"dept_title"`
	KompID         string    `json:"komp_id"`
	KompTitle      string    `json:"komp_title"`
	DirID          string    `json:"dir_id"`
	DirTitle       string    `json:"dir_title"`
	PosLevel       string    `json:"pos_level"`
	SupEmpNo       string    `json:"sup_emp_no"`
	BagID          string    `json:"bag_id"`
	BagTitle       string    `json:"bag_title"`
	SeksiID        *string   `json:"seksi_id"`
	SeksiTitle     *string   `json:"seksi_title"`
	PreNameTitle   string    `json:"pre_name_title"`
	PostNameTitle  string    `json:"post_name_title"`
	NoNPWP         string    `json:"no_npwp"`
	BankAccount    string    `json:"bank_account"`
	BankName       string    `json:"bank_name"`
	MdgDate        string    `json:"mdg_date"`
	PayScale       *string   `json:"PayScale"`
	CCCode         string    `json:"cc_code"`
	Nickname       string    `json:"nickname"`
	JobGrade       *string   `json:"job_grade"`
}

type PihcMasterKaryRt struct {
	EmpNo          string  `json:"emp_no" gorm:"primary_key"`
	Nama           string  `json:"nama"`
	Gender         string  `json:"gender"`
	Agama          string  `json:"agama"`
	StatusKawin    string  `json:"status_kawin"`
	Anak           int8    `json:"anak"` //
	Mdg            string  `json:"mdg"`  //
	EmpGrade       string  `json:"emp_grade"`
	EmpGradeTitle  string  `json:"emp_grade_title"`
	Area           string  `json:"area"`
	AreaTitle      string  `json:"area_title"`
	SubArea        string  `json:"sub_area"`
	SubAreaTitle   string  `json:"sub_area_title"`
	Contract       string  `json:"contract"`
	Pendidikan     string  `json:"pendidikan"`
	Company        string  `json:"company"`
	Lokasi         string  `json:"lokasi"`
	EmployeeStatus string  `json:"employee_status"`
	Email          string  `json:"email"`
	HP             string  `json:"hp"`
	TglLahir       string  `json:"tgl_lahir"`
	PosID          string  `json:"pos_id"`
	PosTitle       string  `json:"pos_title"`
	SupPosID       string  `json:"sup_pos_id"`
	PosGrade       string  `json:"pos_grade"`
	PosKategori    string  `json:"pos_kategori"`
	OrgID          string  `json:"org_id"`
	OrgTitle       string  `json:"org_title"`
	DeptID         string  `json:"dept_id"`
	DeptTitle      string  `json:"dept_title"`
	KompID         string  `json:"komp_id"`
	KompTitle      string  `json:"komp_title"`
	DirID          string  `json:"dir_id"`
	DirTitle       string  `json:"dir_title"`
	PosLevel       string  `json:"pos_level"`
	SupEmpNo       string  `json:"sup_emp_no"`
	BagID          string  `json:"bag_id"`
	BagTitle       string  `json:"bag_title"`
	SeksiID        *string `json:"seksi_id"`
	SeksiTitle     *string `json:"seksi_title"`
	PreNameTitle   string  `json:"pre_name_title"`
	PostNameTitle  string  `json:"post_name_title"`
	NoNPWP         string  `json:"no_npwp"`
	BankAccount    string  `json:"bank_account"`
	BankName       string  `json:"bank_name"`
	MdgDate        string  `json:"mdg_date"`
	PayScale       *string `json:"PayScale"`
	CCCode         string  `json:"cc_code"`
	Nickname       string  `json:"nickname"`
	JobGrade       *string `json:"job_grade"`
}

type DataPegawaiRtDb struct {
	EmpNo          string `json:"emp_no" gorm:"primary_key"`
	Nama           string `json:"nama"`
	Gender         string `json:"gender"`
	Agama          string `json:"agama"`
	StatusKawin    string `json:"status_kawin"`
	Anak           int8   `json:"anak"` //
	Mdg            string `json:"mdg"`  //
	EmpGrade       string `json:"emp_grade"`
	EmpGradeTitle  string `json:"emp_grade_title"`
	Area           string `json:"area"`
	AreaTitle      string `json:"area_title"`
	SubArea        string `json:"sub_area"`
	SubAreaTitle   string `json:"sub_area_title"`
	Contract       string `json:"contract"`
	Pendidikan     string `json:"pendidikan"`
	Company        string `json:"company"`
	Lokasi         string `json:"lokasi"`
	EmployeeStatus string `json:"employee_status"`
	Email          string `json:"email"`
	HP             string `json:"hp"`
	TglLahir       string `json:"tgl_lahir"`
	PosID          string `json:"pos_id"`
	PosTitle       string `json:"pos_title"`
	SupPosID       string `json:"sup_pos_id"`
	PosGrade       string `json:"pos_grade"`
	PosKategori    string `json:"pos_kategori"`
	OrgID          string `json:"org_id"`
	OrgTitle       string `json:"org_title"`
	DeptID         string `json:"dept_id"`
	DeptTitle      string `json:"dept_title"`
	KompID         string `json:"komp_id"`
	KompTitle      string `json:"komp_title"`
	DirID          string `json:"dir_id"`
	DirTitle       string `json:"dir_title"`
	PosLevel       string `json:"pos_level"`
	SupEmpNo       string `json:"sup_emp_no"`
	BagID          string `json:"bag_id"`
	BagTitle       string `json:"bag_title"`
	SeksiID        string `json:"seksi_id"`
	SeksiTitle     string `json:"seksi_title"`
	PreNameTitle   string `json:"pre_name_title"`
	PostNameTitle  string `json:"post_name_title"`
	NoNPWP         string `json:"no_npwp"`
	BankAccount    string `json:"bank_account"`
	BankName       string `json:"bank_name"`
	MdgDate        string `json:"mdg_date"`
	PayScale       string `json:"PayScale"`
	CCCode         string `json:"cc_code"`
	Nickname       string `json:"nickname"`
	JobGrade       string `json:"job_grade"`
}

type SpesifikasiRekap struct {
	EmpNama   string `json:"emp_nama"`
	Nik       string `json:"nik"`
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

type ViewDirektorat struct {
	DirId    string `json:"dir_id"`
	DirTitle string `json:"dir_title"`
}
type ViewKompartemen struct {
	KompID    string `json:"komp_id"`
	KompTitle string `json:"komp_title"`
}
type ViewDepartemen struct {
	DeptID    string `json:"dept_id"`
	DeptTitle string `json:"dept_title"`
}

func (PihcMasterKaryRtDb) TableName() string {
	return "dbo.pihc_master_kary_rt"
}

type PihcMasterKaryRtDbRepo struct {
	DB *gorm.DB
}

type PihcMasterKaryRtRepo struct {
	DB *gorm.DB
}

func NewPihcMasterKaryRtDbRepo(db *gorm.DB) *PihcMasterKaryRtDbRepo {
	return &PihcMasterKaryRtDbRepo{DB: db}
}

func NewPihcMasterKaryRtRepo(db *gorm.DB) *PihcMasterKaryRtRepo {
	return &PihcMasterKaryRtRepo{DB: db}
}

func (t PihcMasterKaryRtDbRepo) FindUserByNIK(nik string) (PihcMasterKaryRtDb, error) {
	var pihc_mkrt PihcMasterKaryRtDb
	err := t.DB.Where("emp_no=?", nik).Take(&pihc_mkrt).Error
	if err != nil {
		return pihc_mkrt, err
	}
	return pihc_mkrt, nil
}

type DataKaryawans struct {
	PihcMasterKaryRtDb
	PihcMasterCompany
	users.UserProfileDB
	profile.Profile
	profile.AboutUs
	profile.PhotoProfile
}

func (t PihcMasterKaryRtDbRepo) FindUserProfileKaryawan(nik string) (DataKaryawans, error) {
	var pihc_mkrt DataKaryawans
	err := t.DB.Table("dbo.pihc_master_kary_rt").
		Select("dbo.pihc_master_kary_rt.*, pmc.*, up.*, p.*,pus.*,pp.*").
		Joins(`INNER JOIN dbo.pihc_master_company pmc ON pmc.code = dbo.pihc_master_kary_rt.company
			   LEFT JOIN public.users_profil up on up.nik = dbo.pihc_master_kary_rt.emp_no
			   LEFT JOIN mobile.profile p on p.nik = dbo.pihc_master_kary_rt.emp_no
			   LEFT JOIN mobile.about_us pus on pus.nik = dbo.pihc_master_kary_rt.emp_no
			   LEFT JOIN mobile.profile_photo pp on pp.emp_no = dbo.pihc_master_kary_rt.emp_no`).
		Where("dbo.pihc_master_kary_rt.emp_no=?", nik).Take(&pihc_mkrt).Error
	if err != nil {
		return pihc_mkrt, err
	}
	return pihc_mkrt, nil
}

func (t PihcMasterKaryRtDbRepo) IsSuperior(pos_id string) bool {
	var count int64
	t.DB.Table("dbo.pihc_master_kary_rt").Where("sup_pos_id=?", pos_id).Count(&count)
	return count != 0
}

func (t PihcMasterKaryRtDbRepo) FindUserByNameArr(name string, nik string) ([]PihcMasterKaryRtDb, error) {
	var pihc_mkrt []PihcMasterKaryRtDb
	if name != "" && nik == "" {
		err := t.DB.Where("lower(nama) like lower(?)", "%"+name+"%").Find(&pihc_mkrt).Error
		if err != nil {
			fmt.Println("ERROR")
			return pihc_mkrt, err
		}
	}
	if nik != "" && name == "" {
		err := t.DB.Where("emp_no like ?", "%"+nik+"%").Find(&pihc_mkrt).Error
		if err != nil {
			fmt.Println("ERROR")
			return pihc_mkrt, err
		}
	}

	return pihc_mkrt, nil
}
func (t PihcMasterKaryRtDbRepo) FindUserArr() ([]PihcMasterKaryRtDb, error) {
	var pihc_mkrt []PihcMasterKaryRtDb

	err := t.DB.Find(&pihc_mkrt).Error
	if err != nil {
		fmt.Println("ERROR")
		return pihc_mkrt, err
	}

	return pihc_mkrt, nil
}

func (t PihcMasterKaryRtDbRepo) FindUserByNameIndiv(name string, nik string) (PihcMasterKaryRtDb, error) {
	var pihc_mkrt PihcMasterKaryRtDb
	if name != "" && nik == "" {
		err := t.DB.Where("lower(nama) like lower(?)", "%"+name+"%").First(&pihc_mkrt).Error
		if err != nil {
			fmt.Println("ERROR")
			return pihc_mkrt, err
		}
	}
	if nik != "" && name == "" {
		err := t.DB.Where("emp_no like ?", "%"+nik+"%").First(&pihc_mkrt).Error
		if err != nil {
			fmt.Println("ERROR")
			return pihc_mkrt, err
		}
	}

	return pihc_mkrt, nil
}

func (t PihcMasterKaryRtRepo) FindUserRekapByNIK(nik string) (*SpesifikasiRekap, error) {
	var pihc_mkrt *SpesifikasiRekap

	err := t.DB.Raw(`
	select nama as emp_nama, emp_no as nik,company as company, pos_id as pos_id,
	   pos_title as pos_title , dept_id as dept_id,
	   dept_title as dept_title, komp_id as komp_id,
	   komp_title as komp_title, dir_id as dir_id, dir_title as dir_title
	from dbo.pihc_master_kary_rt pmkr 
	where emp_no = ?`, nik).Scan(&pihc_mkrt).Error

	if err != nil {
		return nil, err
	}
	return pihc_mkrt, nil
}

func (t PihcMasterKaryRtDbRepo) FindUserByNIKTahunCompCodePeriode(nik string, tahun string, comp_code string, status string) ([]PihcMasterKaryRtDb, error) {
	var pihc_mkrt []PihcMasterKaryRtDb
	err := t.DB.Where("emp_no (IN SELECT manager from tjsl.kegiatan_karyawan where manager=? AND periode=? AND comp_code=? AND status=?)", nik, tahun, comp_code, status).Find(&pihc_mkrt).Error
	if err != nil {
		return pihc_mkrt, err
	}
	return pihc_mkrt, nil
}

func (t PihcMasterKaryRtDbRepo) FindUserAtasanBySupPosID(sup_pos_id string) (PihcMasterKaryRtDb, error) {
	var pihc_mk PihcMasterKaryRtDb
	err := t.DB.Where("pos_id=?", sup_pos_id).Take(&pihc_mk).Error
	if err != nil {
		return pihc_mk, err
	}
	return pihc_mk, nil
}

func (t PihcMasterKaryRtDbRepo) FindDirektoratCompany(company string) ([]ViewDirektorat, error) {
	var vdir []ViewDirektorat
	err := t.DB.Raw(`
		select pmkrt.dir_id, pmkrt.dir_title
			from dbo.pihc_master_kary_rt pmkrt
		where company = ? and (dir_id is not null and dir_title is not null) and (dir_id != '' and dir_title != '')
		group by dir_id ,dir_title 
		order by dir_id
	`, company).Scan(&vdir).Error

	if err != nil {
		return vdir, err
	}
	return vdir, nil
}

func (t PihcMasterKaryRtDbRepo) FindKompartemenCompany(company string, dir_id string) ([]ViewKompartemen, error) {
	var vk []ViewKompartemen
	err := t.DB.Raw(`
		select pmkrt.komp_id, pmkrt.komp_title
			from dbo.pihc_master_kary_rt pmkrt
		where company = ? and dir_id = ? 
			and (komp_id is not null and komp_title is not null) and (komp_id != '' and komp_title != '')
		group by komp_id ,komp_title 
		order by komp_id
	`, company, dir_id).Scan(&vk).Error

	if err != nil {
		return vk, err
	}
	return vk, nil
}
func (t PihcMasterKaryRtDbRepo) FindDepartemenCompany(company string, komp_id string) ([]ViewDepartemen, error) {
	var vdept []ViewDepartemen
	err := t.DB.Raw(`
		select pmkrt.dept_id, pmkrt.dept_title
			from dbo.pihc_master_kary_rt pmkrt
		where company = ? and komp_id = ? 
			and (dept_id is not null and dept_title is not null) and (dept_id != '' and dept_title != '')
		group by dept_id ,dept_title 
		order by dept_id
	`, company, komp_id).Scan(&vdept).Error

	if err != nil {
		return vdept, err
	}
	return vdept, nil
}

//	func (t ViewOrganisasiRepo) FindViewOrganization(nik string) (ViewOrganisasi, error) {
//		var vo ViewOrganisasi
//		err := t.DB.Raw(`
//		select vo.unit1,vo.unit2 , vo.org3, vo.org4
//			from dbo."View_Organisasi" vo
//			join dbo.pihc_master_kary_rt pmkr on pmkr.pos_id = vo."position"
//			where pmkr.emp_no =?`, nik).Scan(&vo).Error
//		if err != nil {
//			return vo, err
//		}
//		return vo, nil
//	}
func (t PihcMasterKaryRtDbRepo) FindUserByNIKArray(nik []string) ([]PihcMasterKaryRtDb, error) {
	var pihc_mk []PihcMasterKaryRtDb
	err := t.DB.Where("emp_no in(?)", nik).Find(&pihc_mk).Error
	if err != nil {
		return pihc_mk, err
	}
	return pihc_mk, nil
}

func (t PihcMasterKaryRtDbRepo) FindUserByKeyArr(key string) ([]DataPegawaiRtDb, error) {
	var pihc_mk []DataPegawaiRtDb

	err := t.DB.Table("dbo.pihc_master_karyawan").Where("lower(nama) like lower(?) OR emp_no like ?", "%"+key+"%", "%"+key+"%").Find(&pihc_mk).Error
	if err != nil {
		fmt.Println("ERROR")
		return pihc_mk, err
	}

	return pihc_mk, nil
}
