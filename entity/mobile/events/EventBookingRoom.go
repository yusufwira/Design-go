package events

import (
	"time"

	"gorm.io/gorm"
)

type EventBookingRoom struct {
	IdBooking int       `json:"id_booking" gorm:"primary_key"`
	CodeRoom  *string    `json:"code_room"`
	IdEvent   int       `json:"id_event"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type DataBookRoomShowDb struct {
	IDBooking    int       `json:"id_booking"`
	CodeRoom     string    `json:"code_room"`
	DateStart    time.Time `json:"date_start"`
	DateEnd      time.Time `json:"date_end"`
	RoomID       string    `json:"room_id"`
	RoomName     string    `json:"room_name"`
	RoomCategory string    `json:"room_category"`
	RoomCompCode string    `json:"room_comp_code"`
	RoomCompName string    `json:"room_comp_name"`
}

type DataBookRoomDateDb struct {
	IDBooking         int       `json:"id_booking"`
	NamaEvent         string    `json:"nama_event"`
	NamaPembuat       string    `json:"nama_pembuat"`
	KompatemenPembuat string    `json:"kompatemen_pembuat"`
	DateStart         time.Time `json:"date_start"`
	DateEnd           time.Time `json:"date_end"`
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

func (t EventBookingRoomRepo) FindRoomBooking(codeRoom *string, idEvent int) (EventBookingRoom, error) {
	var ev_rb EventBookingRoom
	err := t.DB.Where("code_room=? AND id_event=?", codeRoom, idEvent).Find(&ev_rb).Error
	if err != nil {
		return ev_rb, err
	}
	return ev_rb, nil
}

func (t EventBookingRoomRepo) FindRoomBookingByIdBooking(idBooking int) (EventBookingRoom, error) {
	var ev_rb EventBookingRoom
	err := t.DB.Where("id_booking=?", idBooking).Take(&ev_rb).Error
	if err != nil {
		return ev_rb, err
	}
	return ev_rb, nil
}

func (t EventBookingRoomRepo) FindBookRoomShow(codeRoom *string, idEvent int) (DataBookRoomShowDb, error) {
	var data_brs DataBookRoomShowDb
	err := t.DB.Raw(`select ebr.id_booking as id_booking, ebr.code_room,
	ebr.date_start as date_start, ebr.date_end as date_end, emr.id_room as room_id,
	emr.name as room_name, emr.category_room as room_category, emr.comp_code as room_comp_code,
	 pmc.name as room_comp_name
	   from mobile.event_booking_room ebr 
	   join mobile.event_mstr_room emr on emr.id_room = ebr.code_room 
	   join dbo.pihc_master_company pmc on pmc.code = emr.comp_code 
   where code_room = ? and id_event = ?`, codeRoom, idEvent).Scan(&data_brs).Error
	if err != nil {
		return data_brs, err
	}
	return data_brs, nil
}

func (t EventBookingRoomRepo) GetBookingRoomDate(codeRoom string, date string) ([]DataBookRoomDateDb, error) {
	var data_br_date []DataBookRoomDateDb
	err := t.DB.Raw(`select ebr.id_booking as id_booking, e.event_title  as nama_event,
	pmk.nama as nama_pembuat, pmk.komp_title as kompatemen_pembuat,
	e.event_start date_start, e.event_end date_end
   --select *
	   from mobile.event_booking_room ebr 
	   join mobile."event" e on e.id = ebr.id_event
	   join dbo.pihc_master_karyawan pmk on e.created_by = pmk.emp_no 
   where code_room = ? and date(date_start) = ?`, codeRoom, date).Scan(&data_br_date).Error
	if err != nil {
		return data_br_date, err
	}
	return data_br_date, nil
}

func (t EventBookingRoomRepo) FindExistRoom(codeRoom string, idEvent int, timeStart time.Time, timeEnd time.Time) (bool, error) {
	var ev_rb int64
	err := t.DB.Table("mobile.event_booking_room").Where("code_room=? AND id_event!=? AND tsrange(date_start , date_end) && tsrange(? , ?)", codeRoom, idEvent, timeStart, timeEnd).Count(&ev_rb).Error
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
