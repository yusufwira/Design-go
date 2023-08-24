package pihc

import (
	"fmt"

	"gorm.io/gorm"
)

type PihcMasterPosition struct {
	CompanyCode     string `json:"company_code" gorm:"primary_key"`
	Position        string `json:"position"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	PosTitle        string `json:"pos_title"`
	ManagerPos      string `json:"manager_pos"`
	ManagerTitle    string `json:"manager_title"`
	KategoriJabatan string `json:"kategori_jabatan"`
	Grade           string `json:"grade"`
	Lokasi          string `json:"lokasi"`
	LokasiDesc      string `json:"lokasi_desc"`
	OrgUnit         string `json:"org_unit"`
	OrgUnitDesc     string `json:"org_unit_desc"`
	CostCenter      string `json:"costcenter"`
	CostCenterDesc  string `json:"costcenter_desc"`
	DateSent        string `json:"date_sent"`
	TimeSent        string `json:"time_sent"`
	DateGet         string `json:"date_get"`
	TimeGet         string `json:"time_get"`
	Short           string `json:"short"`
	IsHead          string `json:"is_head"`
}

func (PihcMasterPosition) TableName() string {
	return "dbo.pihc_master_position"
}

type PihcMasterPositionRepo struct {
	DB *gorm.DB
}

func NewPihcMasterPositionRepo(db *gorm.DB) *PihcMasterPositionRepo {
	return &PihcMasterPositionRepo{DB: db}
}

func (t PihcMasterPositionRepo) FindUserByPosID(pos_id string) (PihcMasterPosition, error) {
	var pihc_mster_position PihcMasterPosition
	err := t.DB.Where("position=?", pos_id).Take(&pihc_mster_position).Error
	if err != nil {
		fmt.Println("Error retrieving user by NIK:", err)
		return pihc_mster_position, err
	}
	return pihc_mster_position, nil
}
