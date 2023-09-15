package events

import (
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type EventPresence struct {
	IdEventPresence  int       `json:"id_event_presence" gorm:"primary_key"`
	IdEvent          int       `json:"id_event"`
	IdRoom           string    `json:"id_room" gorm:"default:null"`
	EmpNo            string    `json:"emp_no"`
	PresenceDateTime time.Time `json:"presence_date_time"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Nama             string    `json:"nama" gorm:"default:null"`
	Email            string    `json:"email" gorm:"default:null"`
	Instansi         string    `json:"instansi" gorm:"default:null"`
}

type EventDetailEventPresence struct {
	Nik      string      `json:"nik"`
	Nama     null.String `json:"nama"`
	Email    null.String `json:"email"`
	Jabatan  null.String `json:"jabatan"`
	NoTelp   null.String `json:"no_telp"`
	Dept     string      `json:"dept"`
	Presence time.Time   `json:"presence"`
	Instansi null.String `json:"instansi"`
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

func (t EventPresenceRepo) Create(ePresence EventPresence) (EventPresence, error) {
	err := t.DB.Create(&ePresence).Error
	if err != nil {
		return ePresence, err
	}
	return ePresence, nil
}

func (t EventPresenceRepo) FindPresenceIDNIK(id int, nik string) (bool, error) {
	var ev_rb int64
	err := t.DB.Table("mobile.event_presence").Where("id_event=? AND emp_no=?", id, nik).Count(&ev_rb).Error
	if err != nil {
		return true, err
	}
	return ev_rb != 0, nil
}

func (t EventPresenceRepo) FindDetailEventPresence(id int) ([]EventDetailEventPresence, error) {
	var detailEventPresence []EventDetailEventPresence

	err := t.DB.Raw(`
	select pmk.emp_no as nik, pmk.nama as nama, pmk.email as email,
	   	   pmk.pos_title as jabatan, pmk.hp as no_telp, pmk.dept_title as dept,
	   	   ep.presence_date_time as presence, pmc.name as instansi
		from mobile.event_presence ep
		left join dbo.pihc_master_karyawan pmk on pmk.emp_no = ep.emp_no
		left join dbo.pihc_master_company pmc on pmc.code = pmk.company 
	where ep.id_event = ?`, id).Scan(&detailEventPresence).Error

	if err != nil {
		return detailEventPresence, err
	}

	return detailEventPresence, nil
}

func (t EventPresenceRepo) DeleteEventPresence(id int) error {
	var ev_presence EventPresence
	err := t.DB.Where("id_event=?", id).Error
	if err == nil {
		t.DB.Where("id_event= ?", id).Delete(&ev_presence)
		return nil
	}
	return err
}
