package events

import (
	"gorm.io/gorm"
)

type EventMsterRoom struct {
	IDRoom     string `json:"id_room"`
	Name       string `json:"name"`
	RoomQrCode string `json:"room_qr_code"`
	CompCode   string `json:"comp_code"`
	CategoryRoom
	Deskripsi string `json:"deskripsi"`
	Fasilitas *string `json:"fasilitas"`
	Foto      *string `json:"foto"`
	Kapasitas int    `json:"kapasitas"`
	Status    *string `json:"status"`
}

type CategoryRoom struct {
	CategoryRoom string `json:"category_room"`
}

func (EventMsterRoom) TableName() string {
	return "mobile.event_mstr_room"
}

type EventMsterRoomRepo struct {
	DB *gorm.DB
}

func NewEventMsterRoomRepo(db *gorm.DB) *EventMsterRoomRepo {
	return &EventMsterRoomRepo{DB: db}
}

func (t EventMsterRoomRepo) FindEventMasterRoom(room *string) (EventMsterRoom, error) {
	var mster_room EventMsterRoom
	err := t.DB.Where("id_room=?", room).Take(&mster_room).Error
	if err != nil {
		return mster_room, err
	}
	return mster_room, nil
}

func (t EventMsterRoomRepo) FindCategoryRoom(compCode string) ([]EventMsterRoom, error) {
	var category_room []EventMsterRoom

	err := t.DB.Select("DISTINCT(category_room)").Where("comp_code=?", compCode).Find(&category_room).Error
	if err != nil {
		return nil, err
	}

	return category_room, nil
}

func (t EventMsterRoomRepo) FindDefaultCategoryRoom(compCode string) ([]EventMsterRoom, error) {
	var category_room []EventMsterRoom

	err := t.DB.Select("DISTINCT(category_room)").Where("comp_code =?", compCode).Find(&category_room).Error
	if err != nil {
		return nil, err
	}

	return category_room, nil
}

func (t EventMsterRoomRepo) FindRoomEvent(compCode string, categoryRoom string) ([]EventMsterRoom, error) {
	var eventRoom []EventMsterRoom

	err := t.DB.Where("comp_code=? AND category_room=?", compCode, categoryRoom).Find(&eventRoom).Error
	if err != nil {
		return nil, err
	}

	return eventRoom, nil
}

func (t EventMsterRoomRepo) FindDefaultRoomEvent(categoryRoom string) ([]EventMsterRoom, error) {
	var eventRoom []EventMsterRoom

	err := t.DB.Where("category_room=?", categoryRoom).Find(&eventRoom).Error
	if err != nil {
		return nil, err
	}

	return eventRoom, nil
}

