package events

import (
	"gorm.io/gorm"
)

type EventMsterRoom struct {
	IDRoom       string `json:"id_room"`
	Name         string `json:"name"`
	RoomQrCode   string `json:"room_qr_code"`
	CompCode     string `json:"comp_code"`
	CategoryRoom string `json:"category_room"`
	Deskripsi    string `json:"deskripsi"`
	Fasilitas    string `json:"fasilitas"`
	Foto         string `json:"foto"`
	Kapasitas    int    `json:"kapasitas"`
	Status       string `json:"status"`
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

func (t EventMsterRoomRepo) FindEventMasterRoom(room string) (EventMsterRoom, error) {
	var mster_room EventMsterRoom
	err := t.DB.Where("id_room=?", room).Take(&mster_room).Error
	if err != nil {
		return mster_room, err
	}
	return mster_room, nil
}
