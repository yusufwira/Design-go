package events

import (
	"time"

	"gorm.io/gorm"
)

type EventNotulen struct {
	IdNotulen int       `json:"id_notulen" gorm:"primary_key"`
	IdEvent   int       `json:"id_event"`
	Deskripsi string    `json:"deskripsi" gorm:"default:null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type EventNotulenFile struct {
	IdNotulenFile int       `json:"id_notulen_file" gorm:"primary_key"`
	IdNotulen     int       `json:"id_notulen"`
	FileName      string    `json:"file_name" gorm:"default:null"`
	FileUrl       string    `json:"file_url" gorm:"default:null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (EventNotulen) TableName() string {
	return "mobile.event_notulen"
}

type EventNotulenRepo struct {
	DB *gorm.DB
}

func NewEventNotulenRepo(db *gorm.DB) *EventNotulenRepo {
	return &EventNotulenRepo{DB: db}
}

func (t EventNotulenRepo) FindEventNotulenK(idEvent int) (*EventNotulen, error) {
	var ev_notulen EventNotulen
	err := t.DB.Where("id_event=?", idEvent).Take(&ev_notulen).Error
	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	// Return nil and nil error to indicate that no record was found
		// 	ev_notulen.Deskripsi = nil
		// 	return nil, nil
		// }
		return nil, err
	}
	return &ev_notulen, nil
}

func (t EventNotulenRepo) GetDataNotulenFile(idNotulen int) ([]EventNotulenFile, error) {
	var ev_notulen_file []EventNotulenFile
	err := t.DB.Table("mobile.event_notulen_file").Where("id_notulen=?", idNotulen).Find(&ev_notulen_file).Error
	if err != nil {
		return nil, err
	}
	return ev_notulen_file, nil
}

func (t EventNotulenRepo) DeleteEventNotulen(id int) error {
	var ev_rb EventNotulen
	err := t.DB.Table("mobile.event_notulen").Where("id_event=?", id).Take(ev_rb).Error
	if err == nil {
		t.DB.Table("mobile.event_notulen").Where("id_event= ?", id).Delete(&ev_rb)
		t.DeleteEventNotulenFile(ev_rb.IdNotulen)
		return nil
	}
	return err
}

func (t EventNotulenRepo) DeleteEventNotulenFile(id int) error {
	var ev_rb EventNotulenFile
	err := t.DB.Table("mobile.event_notulen_file").Where("id_notulen=?", id).Error
	if err == nil {
		t.DB.Table("mobile.event_notulen_file").Where("id_notulen= ?", id).Delete(&ev_rb)
		return nil
	}
	return err
}
