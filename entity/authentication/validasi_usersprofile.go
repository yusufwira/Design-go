package authentication

import "time"

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

type ValidationStoreAboutUs struct {
	ValidationRequiredNIK
	Desc  string `json:"desc" form:"desc"`
	Hobby string `json:"hobby" form:"hobby"`
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
