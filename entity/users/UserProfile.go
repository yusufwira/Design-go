package users

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	Nik       string  `json:"nik"`
	Alamat    *string `json:"alamat"`
	Kelurahan *string `json:"kelurahan"`
	Kecamatan *string `json:"kecamatan"`
	Kabupaten *string `json:"kabupaten"`
	Provinsi  *string `json:"provinsi"`
	Kodepos   *string `json:"kodepos"`
	Domisili
	PosisiMap   *string   `json:"posisi_map"`
	Email2      *string   `json:"email2" gorm:"default:null"`
	UpdatedBy   *string   `json:"updated_by" gorm:"default:null"`
	NoTelp1     *string   `json:"no_telp1" gorm:"default:null"`
	NoTelp2     *string   `json:"no_telp2" gorm:"default:null"`
	Lat         *string   `json:"lat"`
	Long        *string   `json:"long"`
	Email1      *string   `json:"email1" gorm:"default:null"`
	UpdatedFrom *string   `json:"updated_from" gorm:"default:null"`
	UpdatedDate time.Time `json:"updated_date" gorm:"autoUpdateTime"`
	IsAdmin     *int      `json:"is_admin" gorm:"default:null"`
}

// type UserProfile struct {
// 	Nik       string `json:"nik"`
// 	Alamat    string `json:"alamat"`
// 	Kelurahan string `json:"kelurahan"`
// 	Kecamatan string `json:"kecamatan"`
// 	Kabupaten string `json:"kabupaten"`
// 	Provinsi  string `json:"provinsi"`
// 	Kodepos   string `json:"kodepos"`
// 	Domisili
// 	PosisiMap   string    `json:"posisi_map"`
// 	Email2      string    `json:"email2" gorm:"default:null"`
// 	UpdatedBy   string    `json:"updated_by" gorm:"default:null"`
// 	NoTelp1     string    `json:"no_telp1" gorm:"default:null"`
// 	NoTelp2     string    `json:"no_telp2" gorm:"default:null"`
// 	Lat         string    `json:"lat"`
// 	Long        string    `json:"long"`
// 	Email1      string    `json:"email1" gorm:"default:null"`
// 	UpdatedFrom string    `json:"updated_from" gorm:"default:null"`
// 	UpdatedDate time.Time `json:"updated_date" gorm:"autoUpdateTime" gorm:"autoCreateTime"`
// 	IsAdmin     int       `json:"is_admin" gorm:"default:null"`
// }

type DomisiliDB struct {
	IDDomisili string `json:"id_domisili"`
	Domisili
}

type Domisili struct {
	TipeDomisili *string `json:"tipe_domisili"`
	KetDomisili  *string `json:"ket_domisili"`
}

func (UserProfile) TableName() string {
	return "public.users_profil"
}

type UserProfileRepo struct {
	DB *gorm.DB
}

func NewUserProfileRepo(db *gorm.DB) *UserProfileRepo {
	return &UserProfileRepo{DB: db}
}

func (t UserProfileRepo) Create(up UserProfile) (UserProfile, error) {
	err := t.DB.Create(&up).Error
	if err != nil {
		return up, err
	}
	return up, nil
}

func (t UserProfileRepo) Update(up UserProfile) (UserProfile, error) {
	err := t.DB.Where("nik=?", up.Nik).Save(&up).Error
	if err != nil {
		fmt.Println("ERROR")
		return up, err
	}
	return up, nil
}

func (t UserProfileRepo) FindProfileUsers(nik string) (UserProfile, error) {
	var userProfile UserProfile
	err := t.DB.Where("nik=?", nik).First(&userProfile).Error
	if err != nil {
		fmt.Println("ERROR2")
		return userProfile, err
	}
	return userProfile, nil
}

func (t UserProfileRepo) FindDomisili() []DomisiliDB {
	var domisili []DomisiliDB
	t.DB.Table("public.users_profil_m_domisili").Order("ket_domisili,id_domisili desc").Find(&domisili)
	return domisili
}

func (t UserProfileRepo) FindKetDomisili(tipe string) DomisiliDB {
	var domisili DomisiliDB
	t.DB.Table("public.users_profil_m_domisili").Where("tipe_domisili=?", tipe).Take(&domisili)
	return domisili
}
