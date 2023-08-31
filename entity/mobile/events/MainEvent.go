package events

import (
	"time"

	"gorm.io/gorm"
)

type MainEvent struct {
	Id              int       `json:"id" gorm:"primary_key"`
	EventTitle      string    `json:"event_title"`
	EventDesc       string    `json:"event_desc"`
	EventStart      time.Time `json:"event_start"`
	EventEnd        time.Time `json:"event_end"`
	EventType       string    `json:"event_type"`
	EventUrl        string    `json:"event_url" gorm:"default:null"`
	EventImgName    string    `json:"event_img_name" gorm:"default:null"`
	EventImgUrl     string    `json:"event_img_url" gorm:"default:null"`
	CompCode        string    `json:"comp_code"`
	Status          string    `json:"status"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	ApprovalPerson  string    `json:"approval_person"`
	EventRoom       string    `json:"event_room" gorm:"default:null"`
	EventLocation   string    `json:"event_location" gorm:"default:null"`
	EventKeterangan string    `json:"event_keterangan" gorm:"default:null"`
}

func (MainEvent) TableName() string {
	return "mobile.event"
}

type MainEventRepo struct {
	DB *gorm.DB
}

func NewMainEventRepo(db *gorm.DB) *MainEventRepo {
	return &MainEventRepo{DB: db}
}

func (t MainEventRepo) Create(me MainEvent) (MainEvent, error) {
	err := t.DB.Create(&me).Error
	if err != nil {
		return me, err
	}
	return me, nil
}

func (t MainEventRepo) FindEventMainID(id int) (MainEvent, error) {
	var main_event MainEvent
	err := t.DB.Where("id=?", id).Take(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}

func (t MainEventRepo) FindEventMainIDNIK(id int, nik string) (MainEvent, error) {
	var main_event MainEvent
	err := t.DB.Where("id=? AND created_by=?", id, nik).Take(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}

func (t MainEventRepo) Update(me MainEvent) (MainEvent, error) {
	err := t.DB.Save(&me).Error
	if err != nil {
		return me, err
	}
	return me, nil
}

func (t MainEventRepo) FindEventMainNik(nik string) ([]MainEvent, error) {
	var main_event []MainEvent
	err := t.DB.Where("created_by=?", nik).Find(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}

func (t MainEventRepo) FindEventMainNikMonthYear(nik string, month int, year int, status string) ([]MainEvent, error) {
	var main_event []MainEvent
	err := t.DB.Where("created_by=? AND DATE_PART('month', event_start)=? AND DATE_PART('year', event_start)=? AND status!=?", nik, month, year, status).Find(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}
