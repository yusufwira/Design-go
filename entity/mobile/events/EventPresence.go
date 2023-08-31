package events

import (
	"time"

	"gorm.io/gorm"
)

type EventPresence struct {
	IdEventPresence  int       `json:"id_event_presence" gorm:"primary_key"`
	IdEvent          int       `json:"id_event"`
	IdRoom           string    `json:"id_room"`
	EmpNo            string    `json:"emp_no"`
	PresenceDateTime time.Time `json:"presence_date_time"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Nama             string    `json:"nama"`
	Email            string    `json:"email"`
	Instansi         string    `json:"instansi"`
}

func (EventPresence) TableName() string {
	return "mobile.event_presence"
}

type EventPresenceRepo struct {
	DB *gorm.DB
}

func NewEventPresenceRepo(db *gorm.DB) *EventPresenceRepo {
	return &EventPresenceRepo{DB: db}
}

func (t EventPresenceRepo) FindPresenceIDNIK(id int, nik string) (bool, error) {
	var ev_rb int64
	err := t.DB.Table("mobile.event_presence").Where("id_event=? AND emp_no=?", id, nik).Count(&ev_rb).Error
	if err != nil {
		return true, err
	}
	return ev_rb != 0, nil
}
