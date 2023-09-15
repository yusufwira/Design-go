package pihc

import (
	"gorm.io/gorm"
)

type PihcMasterCompany struct {
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	Address            *string `json:"address"`
	BusinessSector     *string `json:"business_sector"`
	CompanyWebsite     *string `json:"company_website"`
	SortID             int     `json:"sort_id"`
	AssetsLogo         string  `json:"assets_logo"`
	IspupukProducen    int     `json:"ispupuk_producen"`
	UserDefaultLogo    string  `json:"user_default_logo"`
	OrgUnit            string  `json:"org_unit"`
	AssetsLogoFull     string  `json:"assets_logo_full"`
	BasePathUserPhoto  string  `json:"base_path_user_photo"`
	UserPhotoExtension string  `json:"user_photo_extension"`
	AssetsThumbnail    string  `json:"assets_thumbnail"`
	PosID              string  `json:"pos_id"`
}

func (PihcMasterCompany) TableName() string {
	return "dbo.pihc_master_company"
}

type PihcMasterCompanyRepo struct {
	DB *gorm.DB
}

func NewPihcMasterCompanyRepo(db *gorm.DB) *PihcMasterCompanyRepo {
	return &PihcMasterCompanyRepo{DB: db}
}

func (t PihcMasterCompanyRepo) FindPihcMsterCompany(comp_code string) (PihcMasterCompany, error) {
	var pihc_mc PihcMasterCompany
	err := t.DB.Where("code=?", comp_code).Take(&pihc_mc).Error
	if err != nil {
		return pihc_mc, err
	}
	return pihc_mc, nil
}
