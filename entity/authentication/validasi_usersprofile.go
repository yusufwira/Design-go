package authentication

import (
	"time"

	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/mobile/profile"
	"github.com/yusufwira/lern-golang-gin/entity/users"
)

type ValidationSavePersonalInformationEmployee struct {
	EmployeeId   string `form:"employee_id" binding:"required"`
	Alamat       string `form:"alamat" binding:"required"`
	Kelurahan    string `form:"kelurahan" binding:"required"`
	Kecamatan    string `form:"kecamatan" binding:"required"`
	Kabupaten    string `form:"kabupaten" binding:"required"`
	Provinsi     string `form:"provinsi" binding:"required"`
	KodePos      string `form:"kodepos" binding:"required"`
	TipeDomisili string `form:"tipe_domisili" binding:"required"`
	Long         string `form:"long"`
	Lat          string `form:"lat"`
	PosisiMap    string `form:"posisi_map"`
}
type PersonalInformationEmployee struct {
	Nik          string  `json:"nik"`
	Alamat       *string `json:"alamat"`
	Kelurahan    *string `json:"kelurahan"`
	Kecamatan    *string `json:"kecamatan"`
	Kabupaten    *string `json:"kabupaten"`
	Provinsi     *string `json:"provinsi"`
	KodePos      *string `json:"kodepos"`
	TipeDomisili *string `json:"tipe_domisili"`
	KetDomisili  *string `json:"ket_domisili"`
	PosisiMap    *string `json:"posisi_map"`
	Lat          *string `json:"lat"`
	Long         *string `json:"long"`
	UpdatedFrom  *string `json:"updated_from"`
	UpdatedDate  string  `json:"updated_date"`
}

type ValidationGetName struct {
	Nik  string `json:"nik" form:"nik"`
	Name string `json:"name" form:"name"`
}

type ValidationStoreContactInformation struct {
	Nik     string `form:"nik" binding:"required"`
	NoTelp1 string `form:"no_telp_1"`
	NoTelp2 string `form:"no_telp_2"`
	Email1  string `form:"email_1"`
	Email2  string `form:"email_2"`
}

type ValidationGetPersonalInformation struct {
	ValidationRequiredNIK
}

type ValidationRequiredNIK struct {
	NIK string `json:"nik" form:"nik" binding:"required"`
}

type ValidationStoreProfile struct {
	ValidationRequiredNIK
	Bio           *string `json:"bio"`
	LinkTwitter   *string `json:"link_twitter"`
	LinkInstagram *string `json:"link_instagram"`
	LinkWebsite   *string `json:"link_website"`
	LinkFacebook  *string `json:"link_facebook"`
	LinkTiktok    *string `json:"link_tiktok"`
}

type ValidationDataPegawai struct {
	Key string `form:"key" json:"key" binding:"required"`
}

type DataPegawai struct {
	Nik                 string  `json:"nik"`
	Nama                string  `json:"nama"`
	DeptTitle           string  `json:"dept_title"`
	CompanyName         string  `json:"company_name"`
	Skill               *string `json:"skill"`
	PhotoProfile        string  `json:"photo_profile"`
	PhotoProfileDefault string  `json:"photo_profile_default"`
	CompanyLogo         string  `json:"company_logo"`
}

type ValidationStoreAboutUs struct {
	ValidationRequiredNIK
	Desc  string `json:"desc" form:"desc"`
	Hobby string `json:"hobby" form:"hobby"`
}

type ValidationStoreSkill struct {
	ValidationRequiredNIK
	Category []Category `json:"category"`
}

type ValidationUpdateSkill struct {
	ID            int `json:"id" form:"id"`
	IdParentSkill int `json:"id_parent_skill" form:"id_parent_skill"`
	ValidationRequiredNIK
	Type string `json:"type" form:"type"`
	Name string `json:"name" form:"name" binding:"required"`
}

type ValidationDeleteSkill struct {
	ID   int    `json:"id" form:"id"`
	Type string `json:"type" form:"type"`
}

type Category struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Skill []Skill `json:"skill"`
}

type Skill struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	SubSkill []SubSkill `json:"sub_skill"`
}

type SubSkill struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetStoreProfile struct {
	NIK           string    `json:"nik"`
	Bio           *string   `json:"bio"`
	LinkTwitter   *string   `json:"link_twitter"`
	LinkInstagram *string   `json:"link_instagram"`
	LinkWebsite   *string   `json:"link_website"`
	LinkFacebook  *string   `json:"link_facebook"`
	LinkTiktok    *string   `json:"link_tiktok"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ID            int       `json:"id"`
}

type GetSocialMedia struct {
	NIK           string  `json:"nik"`
	LinkTwitter   *string `json:"link_twitter"`
	LinkFacebook  *string `json:"link_facebook"`
	LinkInstagram *string `json:"link_instagram"`
	LinkTiktok    *string `json:"link_tiktok"`
	LinkWebsite   *string `json:"link_website"`
	Bio           *string `json:"bio"`
}

type ContactInformation struct {
	NoTelp1 *string `json:"no_telp_1"`
	NoTelp2 *string `json:"no_telp_2"`
	Email1  *string `json:"email_1"`
	Email2  *string `json:"email_2"`
}

type PengalamanKerja struct {
	ValidFrom    string `json:"valid_from"`
	ValidTo      string `json:"valid_to"`
	Grade        string `json:"grade"`
	PositionId   string `json:"position_id"`
	PositionName string `json:"position_name"`
	Unit1        string `json:"unit_1"`
	Unit2        string `json:"unit_2"`
}

type ProfilePribadi struct {
	pihc.PihcMasterKary
	Domisili            *users.UserProfile     `json:"domisili"`
	ProfileMobile       *MobileProfile         `json:"profile_mobile"`
	AboutUs             *profile.AboutUs       `json:"about_us"`
	Companys            pihc.PihcMasterCompany `json:"companys"`
	Skill               []ShowSkills           `json:"skill"`
	CompanyLogo         string                 `json:"company_logo"`
	PhotoProfile        string                 `json:"photo_profile"`
	PhotoProfileDefault string                 `json:"photo_profile_default"`
	Organisasi          []string               `json:"organisasi"`
}
type ValidationPhotoProfile struct {
	ValidationRequiredNIK
}

type MobileProfile struct {
	profile.Profile
	UserProfile users.UserProfile `json:"user_profile"`
}

type ShowSkills struct {
	profile.ProfileSkill
	Skill []ProfileMainSkill `json:"skill"`
}

type ProfileMainSkill struct {
	profile.ProfileSkill
	SubSkill []ProfileSubSkill `json:"sub_skill"`
}

type ProfileSubSkill struct {
	profile.ProfileSkill
}

type ShowAboutUs struct {
	Id           int       `json:"id"`
	Nik          string    `json:"nik"`
	AboutUsDesc  *string   `json:"about_us_desc"`
	AboutUsHobby *string   `json:"about_us_hobby"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DataDomisili struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name"`
}
