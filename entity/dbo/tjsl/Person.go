package tjsl

import (
	"time"

	"gorm.io/gorm"
)

type KoordinatorPerson struct {
	Id            int       `json:"id" gorm:"primary_key"`
	KoordinatorId int       `json:"koordinator_id"`
	NIK           string    `json:"nik"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (KoordinatorPerson) TableName() string {
	return "tjsl.koordinator_person"
}

type KoordinatorPersonRepo struct {
	DB *gorm.DB
}

func NewKoordinatorPersonRepo(db *gorm.DB) *KoordinatorPersonRepo {
	return &KoordinatorPersonRepo{DB: db}
}

func (t KoordinatorPersonRepo) Create(koorprson KoordinatorPerson) (KoordinatorPerson, error) {
	err := t.DB.Create(&koorprson).Error
	if err != nil {
		return koorprson, err
	}
	return koorprson, nil
}

func (t KoordinatorPersonRepo) FindDataKoorPersonID(id int) []KoordinatorPerson {
	var koor_person []KoordinatorPerson
	t.DB.Where("koordinator_id=?", id).Find(&koor_person)
	return koor_person
}

func (t KoordinatorPersonRepo) FindDataKoorPersonNIK(NIK string) []KoordinatorPerson {
	var koor_person []KoordinatorPerson
	t.DB.Distinct("koordinator_id").Where("nik=?", NIK).Find(&koor_person)
	return koor_person
}

func (t KoordinatorPersonRepo) DelPersonID(koor_id int) ([]KoordinatorPerson, error) {
	var data []KoordinatorPerson
	err := t.DB.Where("koordinator_id = ?", koor_id).First(&data).Error
	if err == nil {
		t.DB.Where("koordinator_id = ?", koor_id).Delete(&data)
		return data, nil
	}
	return data, err
}
