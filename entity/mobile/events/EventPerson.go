package events

import (
	"time"

	"gorm.io/gorm"
)

type EventPerson struct {
	Id                int       `json:"id" gorm:"primary_key"`
	IdEvent           int       `json:"id_event"`
	IdParent          int       `json:"id_parent" gorm:"default:null"`
	Nik               string    `json:"nik"`
	StatusKehadiran   string    `json:"status_kehadiran"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	KetKetidakhadiran string    `json:"ket_ketidakhadiran" gorm:"default:null"`
}

type EventCounts struct {
	CountGuest      int64 `json:"count_guest"`
	CountHadir      int64 `json:"count_hadir"`
	CountMenunggu   int64 `json:"count_menunggu"`
	CountTidakHadir int64 `json:"count_tidak_hadir"`
}

func (EventPerson) TableName() string {
	return "mobile.event_person"
}

type EventPersonRepo struct {
	DB *gorm.DB
}

func NewEventPersonRepo(db *gorm.DB) *EventPersonRepo {
	return &EventPersonRepo{DB: db}
}

func (t EventPersonRepo) Create(me EventPerson) (EventPerson, error) {
	err := t.DB.Create(&me).Error
	if err != nil {
		return me, err
	}
	return me, nil
}

func (t EventPersonRepo) FindEventPersonID(id int) ([]EventPerson, error) {
	var event_person []EventPerson
	err := t.DB.Where("id_event=?", id).Find(&event_person).Error
	if err != nil {
		return event_person, err
	}
	return event_person, nil
}

func (t EventPersonRepo) FindEventPersonIDNIK(id int, nik string) (EventPerson, error) {
	var event_person EventPerson
	err := t.DB.Where("id_event=? AND nik=?", id, nik).Take(&event_person).Error
	if err != nil {
		return event_person, err
	}
	return event_person, nil
}

func (t EventPersonRepo) DelParticipationLama(event_id int, list_id []int) {
	t.DB.Where("id_event = ? AND id not in(?)", event_id, list_id).Delete(&EventPerson{})
}

func (t EventPersonRepo) GetEventCounts(id int) (EventCounts, error) {
	var counts EventCounts

	err := t.DB.Raw(`
        SELECT
            COUNT(*) FILTER (WHERE status_kehadiran = 'hadir') AS count_hadir,
            COUNT(*) FILTER (WHERE status_kehadiran = 'menunggu') AS count_menunggu,
            COUNT(*) FILTER (WHERE status_kehadiran = 'tidak_hadir') AS count_tidak_hadir
        FROM mobile.event_person
        WHERE id_event = ?`, id).Scan(&counts).Error

	if err != nil {
		return counts, err
	}

	counts.CountGuest = counts.CountHadir + counts.CountMenunggu + counts.CountTidakHadir

	return counts, nil
}
