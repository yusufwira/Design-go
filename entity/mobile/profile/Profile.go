package profile

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID            int       `json:"id"`
	Nik           string    `json:"nik"`
	Bio           *string   `json:"bio" gorm:"default:null"`
	LinkTwitter   *string   `json:"link_twitter" gorm:"default:null"`
	LinkInstagram *string   `json:"link_instagram" gorm:"default:null"`
	LinkWebsite   *string   `json:"link_website" gorm:"default:null"`
	LinkFacebook  *string   `json:"link_facebook" gorm:"default:null"`
	LinkTiktok    *string   `json:"link_tiktok" gorm:"default:null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Suku          *string   `json:"suku"`
	GolonganDarah *string   `json:"golongan_darah" gorm:"default:null"`
	VisiPribadi   *string   `json:"visi_pribadi" gorm:"default:null"`
	NilaiPribadi  *string   `json:"nilai_pribadi" gorm:"default:null"`
	Interest      *string   `json:"interest" gorm:"default:null"`
	LinkLinkedin  *string   `json:"link_linkedin" gorm:"default:null"`
}

type AboutUs struct {
	ID           int       `json:"id"`
	Nik          string    `json:"nik"`
	AboutUsDesc  *string   `json:"about_us_desc" gorm:"default:null"`
	AboutUsHobby *string   `json:"about_us_hobby" gorm:"default:null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Profile) TableName() string {
	return "mobile.profile"
}

func (AboutUs) TableName() string {
	return "mobile.about_us"
}

type ProfileRepo struct {
	DB *gorm.DB
}

type AboutUsRepo struct {
	DB *gorm.DB
}

func NewProfileRepo(db *gorm.DB) *ProfileRepo {
	return &ProfileRepo{DB: db}
}

func NewAboutUsRepo(db *gorm.DB) *AboutUsRepo {
	return &AboutUsRepo{DB: db}
}

func (t ProfileRepo) Create(p Profile) (Profile, error) {
	err := t.DB.Create(&p).Error
	if err != nil {
		return p, err
	}
	return p, nil
}

func (t ProfileRepo) Update(p Profile) (Profile, error) {
	err := t.DB.Where("nik=?", p.Nik).Save(&p).Error
	if err != nil {
		fmt.Println("ERROR")
		return p, err
	}
	return p, nil
}

func (t AboutUsRepo) Create(au AboutUs) (AboutUs, error) {
	err := t.DB.Create(&au).Error
	if err != nil {
		return au, err
	}
	return au, nil
}

func (t AboutUsRepo) Update(au AboutUs) (AboutUs, error) {
	err := t.DB.Where("nik=?", au.Nik).Save(&au).Error
	if err != nil {
		fmt.Println("ERROR")
		return au, err
	}
	return au, nil
}

func (t ProfileRepo) FindProfile(nik string) (Profile, error) {
	var profile Profile
	err := t.DB.Where("nik=?", nik).First(&profile).Error
	if err != nil {
		fmt.Println("ERROR2")
		return profile, err
	}
	return profile, nil
}
func (t AboutUsRepo) FindProfileAboutUs(nik string) (AboutUs, error) {
	var au AboutUs
	err := t.DB.Where("nik=?", nik).First(&au).Error
	// err := t.DB.Table("mobile.about_us").Where("nik=?", nik).First(&au).Error
	if err != nil {
		fmt.Println("ERROR2")
		return au, err
	}
	return au, nil
}
