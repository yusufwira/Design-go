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

type ProfileSkill struct {
	ID            int       `json:"id"`
	IdParentSkill *int      `json:"id_parent_skill"`
	Nik           string    `json:"nik"`
	Type          string    `json:"type"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Profile) TableName() string {
	return "mobile.profile"
}

func (AboutUs) TableName() string {
	return "mobile.about_us"
}

func (ProfileSkill) TableName() string {
	return "mobile.profile_skill"
}

type ProfileRepo struct {
	DB *gorm.DB
}

type AboutUsRepo struct {
	DB *gorm.DB
}

type ProfileSkillRepo struct {
	DB *gorm.DB
}

func NewProfileRepo(db *gorm.DB) *ProfileRepo {
	return &ProfileRepo{DB: db}
}

func NewProfileSkillRepo(db *gorm.DB) *ProfileSkillRepo {
	return &ProfileSkillRepo{DB: db}
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
		return p, err
	}
	return p, nil
}
func (t ProfileSkillRepo) Create(p ProfileSkill) (ProfileSkill, error) {
	err := t.DB.Create(&p).Error
	if err != nil {
		return p, err
	}
	return p, nil
}

func (t ProfileSkillRepo) CreateC(p []ProfileSkill) ([]ProfileSkill, error) {
	err := t.DB.Create(&p).Error
	if err != nil {
		return p, err
	}

	var indexA []int
	var indexB []int
	iterasiA := 0
	iterasiB := 0

	for i, data := range p {
		if data.Type == "category_skill" {
			if len(indexA) > 0 {
				iterasiA++
			}
			indexA = append(indexA, data.ID)
		}
		if data.Type == "main_skill" {
			if len(indexA) > 0 {
				data.IdParentSkill = &indexA[iterasiA]
			}
			if len(indexB) > 0 {
				iterasiB++
			}
			indexB = append(indexB, data.ID)
		}
		if data.Type == "sub_skill" {
			if len(indexB) > 0 {
				data.IdParentSkill = &indexB[iterasiB]
			}
		}
		p[i] = data
	}

	err1 := t.DB.Save(&p).Error
	if err1 != nil {
		return p, err1
	}

	return p, nil
}

func (t ProfileSkillRepo) Update(p ProfileSkill) (ProfileSkill, error) {
	err := t.DB.Where("id=?", p.ID).Save(&p).Error
	if err != nil {
		fmt.Println("ERROR")
		return p, err
	}

	return p, nil
}

func (t ProfileSkillRepo) UpdateC(p []ProfileSkill) ([]ProfileSkill, error) {
	var indexA []int
	var indexB []int
	iterasiA := 0
	iterasiB := 0

	for i, data := range p {
		if data.Type == "category_skill" {
			if len(indexA) > 0 {
				iterasiA++
			}
			indexA = append(indexA, data.ID)
		}
		if data.Type == "main_skill" {
			if len(indexA) > 0 {
				data.IdParentSkill = &indexA[iterasiA]
			}
			if len(indexB) > 0 {
				iterasiB++
			}
			indexB = append(indexB, data.ID)
		}
		if data.Type == "sub_skill" {
			if len(indexB) > 0 {
				data.IdParentSkill = &indexB[iterasiB]
			}
		}
		p[i] = data
	}

	err := t.DB.Save(&p).Error
	if err != nil {
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
func (t ProfileSkillRepo) FindProfileCategorySkill(id int) (ProfileSkill, error) {
	var profile ProfileSkill
	err := t.DB.Where("id=?", id).First(&profile).Error
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (t ProfileSkillRepo) FindProfileCategorySkillArr(nik string) ([]ProfileSkill, error) {
	var profile []ProfileSkill
	err := t.DB.Where("nik=? AND id_parent_skill IS NULL", nik).Find(&profile).Error
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (t ProfileSkillRepo) FindProfileSkillArr(nik string, id int) ([]ProfileSkill, error) {
	var profile []ProfileSkill
	err := t.DB.Where("nik=? AND id_parent_skill=?", nik, id).Find(&profile).Error
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (t ProfileSkillRepo) GetProfileSkillArr(nik string, typeSkill string) ([]ProfileSkill, error) {
	results := []ProfileSkill{}
	var ps string
	if typeSkill == "category_skill" {
		ps = "ps"
	}
	if typeSkill == "main_skill" {
		ps = "ps2"
	}
	if typeSkill == "sub_skill" {
		ps = "ps3"
	}
	query := fmt.Sprintf(`select distinct(%s.*)
                          from mobile.profile_skill ps
                          left join mobile.profile_skill ps2 on ps2.id_parent_skill = ps.id
                          left join mobile.profile_skill ps3 on ps3.id_parent_skill = ps2.id
                          where ps.nik = ? and ps.type ='category_skill'`, ps)

	err := t.DB.Raw(query, nik).
		Scan(&results).Error

	if err != nil {
		return results, err
	}

	return results, nil
}

func (t ProfileSkillRepo) FindProfileSkill(id int, parent int) (ProfileSkill, error) {
	var profile ProfileSkill
	err := t.DB.Where("id=? AND id_parent_skill=?", id, parent).First(&profile).Error
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (t AboutUsRepo) FindProfileAboutUs(nik string) (AboutUs, error) {
	var au AboutUs
	err := t.DB.Where("nik=?", nik).First(&au).Error
	if err != nil {
		return au, err
	}
	return au, nil
}
