package events

import (
	"time"

	"gorm.io/gorm"
)

type EventMateriFile struct {
	IdMateriFile int       `json:"id_materi_file" gorm:"primary_key"`
	IdEvent      int       `json:"id_event"`
	FileName     string    `json:"file_name" gorm:"default:null"`
	FileUrl      string    `json:"file_url" gorm:"default:null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (EventMateriFile) TableName() string {
	return "mobile.event_materi_file"
}

type EventMateriFileRepo struct {
	DB *gorm.DB
}

func NewEventMateriFileRepo(db *gorm.DB) *EventMateriFileRepo {
	return &EventMateriFileRepo{DB: db}
}

func (t EventMateriFileRepo) FindEventMateriFile(idEvent int) ([]EventMateriFile, error) {
	var ev_materi_file []EventMateriFile
	err := t.DB.Where("id_event=?", idEvent).Find(&ev_materi_file).Error
	if err != nil {
		return nil, err
	}
	return ev_materi_file, nil
}

func (t EventMateriFileRepo) DeleteEventMateriFile(id int) error {
	var ev_materi_file EventMateriFile
	err := t.DB.Where("id_event=?", id).Error
	if err == nil {
		t.DB.Where("id_event= ?", id).Delete(&ev_materi_file)
		return nil
	}
	return err
}
