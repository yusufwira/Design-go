package pihc

import (
	"fmt"

	"gorm.io/gorm"
)

type PihcMasterKaryRt struct {
	EmpNo          string `json:"emp_no" gorm:"primary_key"`
	Nama           string `json:"nama"`
	Gender         string `json:"gender"`
	Agama          string `json:"agama"`
	StatusKawin    string `json:"status_kawin"`
	Anak           int8   `json:"anak"` //
	Mdg            int    `json:"mdg"`  //
	EmpGrade       string `json:"emp_grade"`
	EmpGradeTitle  string `json:"emp_grade_title"`
	Area           string `json:"area"`
	AreaTitle      string `json:"area_title"`
	SubArea        string `json:"sub_area"`
	SubAreaTtitle  string `json:"sub_area_title"`
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
	SubPosID       string `json:"sup_pos_id"`
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

func (PihcMasterKaryRt) TableName() string {
	return "dbo.pihc_master_kary_rt"
}

type PihcMasterKaryRtRepo struct {
	DB *gorm.DB
}

func NewPihcMasterKaryRtRepo(db *gorm.DB) *PihcMasterKaryRtRepo {
	return &PihcMasterKaryRtRepo{DB: db}
}

func (t PihcMasterKaryRtRepo) FindUserByNIK(nik string) (PihcMasterKaryRt, error) {
	var pihc_mkrt PihcMasterKaryRt
	err := t.DB.Where("emp_no=?", nik).Take(&pihc_mkrt).Error
	if err != nil {
		fmt.Println("Error retrieving user by NIK:", err)
		return pihc_mkrt, err
	}
	return pihc_mkrt, nil
}
