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
	AssetsLogoFull     *string `json:"assets_logo_full"`
	BasePathUserPhoto  *string `json:"base_path_user_photo"`
	UserPhotoExtension *string `json:"user_photo_extension"`
	AssetsThumbnail    string  `json:"assets_thumbnail"`
	PosID              string  `json:"pos_id"`
}

type ViewOrganisasi struct {
	Unit1 string `json:"unit1"`
	Unit2 string `json:"unit2"`
	Org3  string `json:"org3"`
	Org4  string `json:"org4"`
}

func (PihcMasterCompany) TableName() string {
	return "dbo.pihc_master_company"
}

func (ViewOrganisasi) TableName() string {
	return `dbo."View_Organisasi"`
}

type PihcMasterCompanyRepo struct {
	DB *gorm.DB
}

func NewPihcMasterCompanyRepo(db *gorm.DB) *PihcMasterCompanyRepo {
	return &PihcMasterCompanyRepo{DB: db}
}

type ViewOrganisasiRepo struct {
	DB *gorm.DB
}

func NewViewOrganisasiRepo(db *gorm.DB) *ViewOrganisasiRepo {
	return &ViewOrganisasiRepo{DB: db}
}

func (t PihcMasterCompanyRepo) FindPihcMsterCompany(comp_code string) (PihcMasterCompany, error) {
	var pihc_mc PihcMasterCompany
	err := t.DB.Where("code=?", comp_code).Take(&pihc_mc).Error
	if err != nil {
		return pihc_mc, err
	}
	return pihc_mc, nil
}
func (t PihcMasterCompanyRepo) FindPihcMsterCompanyArray(comp_code []string) ([]PihcMasterCompany, error) {
	var pihc_mc []PihcMasterCompany
	err := t.DB.Where("code in(?)", comp_code).Find(&pihc_mc).Error
	if err != nil {
		return pihc_mc, err
	}
	return pihc_mc, nil
}

func (t ViewOrganisasiRepo) FindViewOrganization(nik string) (ViewOrganisasi, error) {
	var vo ViewOrganisasi
	err := t.DB.Raw(`
	select vo.unit1,vo.unit2 , vo.org3, vo.org4
		from dbo."View_Organisasi" vo
		join dbo.pihc_master_kary_rt pmkr on pmkr.pos_id = vo."position"
		where pmkr.emp_no =?`, nik).Scan(&vo).Error
	if err != nil {
		return vo, err
	}
	return vo, nil
}
