package pihc

import "gorm.io/gorm"

type PihcKaryawanMutasiPI struct {
	PiEmpNo string `json:"pi_emp_no"`
	Nama    string `json:"nama"`
	EmpNo   string `json:"emp_no"`
	Company string `json:"company"`
}

func (PihcKaryawanMutasiPI) TableName() string {
	return "dbo.pihc_karyawan_mutasi_pi"
}

type PihcKaryawanMutasiPIRepo struct {
	DB *gorm.DB
}

func NewPihcKaryawanMutasiPIRepo(db *gorm.DB) *PihcKaryawanMutasiPIRepo {
	return &PihcKaryawanMutasiPIRepo{DB: db}
}

func (t PihcKaryawanMutasiPIRepo) FindPihcKaryawanMutasiPI(nik string) (PihcKaryawanMutasiPI, error) {
	var pihc_mutasi PihcKaryawanMutasiPI
	err := t.DB.Where("pi_emp_no=?", nik).Take(&pihc_mutasi).Error
	if err != nil {
		return pihc_mutasi, err
	}
	return pihc_mutasi, nil
}
