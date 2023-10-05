package pihc

import (
	"fmt"
	"time"

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
	SubAreaTtitle  string    `json:"sub_area_title"`
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
	SubPosID       string    `json:"sup_pos_id"`
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
	SubAreaTtitle  string  `json:"sub_area_title"`
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
	SubPosID       string  `json:"sup_pos_id"`
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

func (t PihcMasterKaryRtDbRepo) FindUserByName(name string, nik string) ([]PihcMasterKaryRtDb, error) {
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
