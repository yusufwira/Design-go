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

func (t PihcMasterCompanyRepo) FindAllCompany() ([]PihcMasterCompany, error) {
	var pihc_mc []PihcMasterCompany
	err := t.DB.Order("code ASC").Find(&pihc_mc).Error
	if err != nil {
		return pihc_mc, err
	}
	return pihc_mc, nil
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

//	func (t ViewOrganisasiRepo) FindViewOrganization(nik string) (ViewOrganisasi, error) {
//		var vo ViewOrganisasi
//		err := t.DB.Table(`dbo."View_Organisasi"`).
//			Select(`dbo."View_Organisasi".unit1, dbo."View_Organisasi".unit2,
//				dbo."View_Organisasi".org3, dbo."View_Organisasi".org4`).
//			Joins(`inner join dbo.pihc_master_kary_rt pmkr on pmkr.pos_id = dbo."View_Organisasi"."position"`).
//			Where("pmkr.emp_no =?", nik).Take(&vo).Error
//		if err != nil {
//			return vo, err
//		}
//		return vo, nil
//	}
func (t PihcMasterPositionRepo) FindViewOrganization(nik string) (ViewOrganisasi, error) {
	var vo ViewOrganisasi
	err := t.DB.Table("dbo.pihc_master_position").
		Select(`
			CASE
				WHEN LEFT(dbo.pihc_master_position.grade::text, 1) >= '3'::text THEN
					COALESCE(b4.org_unit_desc, b3.org_unit_desc, b2.org_unit_desc, dbo.pihc_master_position.org_unit_desc)
				ELSE
					dbo.pihc_master_position.org_unit_desc
			END AS unit1,
			CASE
				WHEN LEFT(dbo.pihc_master_position.grade::text, 1) >= '3'::text THEN
					COALESCE(b4.org_unit_desc, b3.org_unit_desc, b2.org_unit_desc, dbo.pihc_master_position.org_unit_desc)
				ELSE
					b2.org_unit_desc
			END AS unit2,
			b3.org_unit_desc AS org3,
			b4.org_unit_desc AS org4
		`).
		Joins(`
        LEFT JOIN dbo.pihc_master_position b2 ON dbo.pihc_master_position.manager_pos::text = b2."position"::text
        LEFT JOIN dbo.pihc_master_position b3 ON b2.manager_pos::text = b3."position"::text
        LEFT JOIN dbo.pihc_master_position b4 ON b3.manager_pos::text = b4."position"::text
        INNER JOIN dbo.pihc_master_kary_rt pmkr ON pmkr.pos_id = dbo.pihc_master_position."position"
    	`).
		Where("pmkr.emp_no = ?", nik).
		Order(`dbo.pihc_master_position."position"`).
		Take(&vo).Error

	if err != nil {
		return vo, err
	}
	return vo, nil
}
