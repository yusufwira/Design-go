package events

import (
	"time"

	"gorm.io/gorm"
)

type EventBookingRoom struct {
	IdBooking int       `json:"id_booking" gorm:"primary_key"`
	CodeRoom  string    `json:"code_room"`
	IdEvent   int       `json:"id_event"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (EventBookingRoom) TableName() string {
	return "mobile.event_booking_room"
}

type EventBookingRoomRepo struct {
	DB *gorm.DB
}

func NewEventBookingRoomRepo(db *gorm.DB) *EventBookingRoomRepo {
	return &EventBookingRoomRepo{DB: db}
}

func (t EventBookingRoomRepo) Create(ev_br EventBookingRoom) (EventBookingRoom, error) {
	err := t.DB.Create(&ev_br).Error
	if err != nil {
		return ev_br, err
	}
	return ev_br, nil
}

func (t EventBookingRoomRepo) FindRoomBooking(codeRoom string, idEvent int) (EventBookingRoom, error) {
	var ev_rb EventBookingRoom
	err := t.DB.Where("code_room=? AND id_event=?", codeRoom, idEvent).Find(&ev_rb).Error
	if err != nil {
		return ev_rb, err
	}
	return ev_rb, nil
}

func (t EventBookingRoomRepo) FindExistRoom(codeRoom string, timeStart time.Time, timeEnd time.Time) (bool, error) {
	var ev_rb int64
	err := t.DB.Table("mobile.event_booking_room").Where("code_room=? AND tsrange(date_start , date_end) && tsrange(? , ?)", codeRoom, timeStart, timeEnd).Count(&ev_rb).Error
	if err != nil {
		return true, err
	}
	return ev_rb != 0, nil
}

func (t EventBookingRoomRepo) DeleteRoomBooking(id int) error {
	var ev_rb []EventBookingRoom
	err := t.DB.Where("id_event=?", id).Error
	if err == nil {
		t.DB.Where("id_event = ?", id).Delete(&ev_rb)
		return nil
	}
	return err
}

func (t EventBookingRoomRepo) Update(ev_rb EventBookingRoom) (EventBookingRoom, error) {
	err := t.DB.Save(&ev_rb).Error
	if err != nil {
		return ev_rb, err
	}
	return ev_rb, nil
}
