package authentication

import (
	"time"

	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/mobile/events"
	"gopkg.in/guregu/null.v4"
)

type ValidasiEvent struct {
	ID         int     `json:"id"`
	Nik        int     `json:"nik" binding:"required"`
	Title      string  `json:"title" binding:"required"`
	Desc       string  `json:"desc"`
	Type       string  `json:"type" binding:"required"`
	URL        *string `json:"url"`
	Location   *string `json:"location"`
	Image      string  `json:"image"`
	Start      string  `json:"start" binding:"required"`
	End        string  `json:"end" binding:"required"`
	IsPublic   int     `json:"is_public"`
	FileMateri []struct {
		IdMateriFile int    `json:"id_materi_file"`
		IdEvent      int    `json:"id_event"`
		FileName     string `json:"file_name"`
		FileURL      string `json:"file_url"`
	} `json:"file_materi"`
	Status         string  `json:"status" binding:"required"`
	IDRoom         *string `json:"id_room"`
	Person         *string `json:"person"`
	ApprovalPerson *string `json:"approval_person"`
}

type ValidasiStoreBookingRoom struct {
	DateStart string `form:"date_start" binding:"required"`
	DateEnd   string `form:"date_end" binding:"required"`
	IDRoom    string `form:"id_room" binding:"required"`
	IDEvent   int    `form:"id_event"`
}

type ValidasiStoreEventPresence struct {
	IdPresence int `form:"id_presence" binding:"required"`
	ValidasiKonfirmasiNik
	Type             string `form:"type" binding:"required"`
	PresenceDateTime string `form:"presence_date_time"`
}

type ValidasiUpdateStatusEvent struct {
	ValidasiKonfirmasi
	Keterangan string `form:"keterangan"`
}

type ValidasiKonfirmasi struct {
	ValidasiKonfirmasiNik
	ValidasiKonfirmasiEventID
	Status string `form:"status" binding:"required"`
}

type ValidasiKonfirmasiNik struct {
	Nik string `form:"nik" binding:"required"`
}
type ValidasiKonfirmasiCategoryRoom struct {
	CategoryRoom string `form:"category_room" binding:"required"`
}

type ValidasiKonfirmasiGetEventRoom struct {
	ValidasiKonfirmasiNik
	ValidasiKonfirmasiCategoryRoom
}

type ValidasiKonfirmasiEventID struct {
	EventID int `form:"event_id" binding:"required"`
}

type ValidasiStoreDispose struct {
	ValidasiKonfirmasiNik
	ValidasiKonfirmasiEventID
	Dispose string `form:"dispose" binding:"required"`
}

type ValidasiGetDataDispose struct {
	Nik     string `form:"nik" binding:"required"`
	IdEvent int    `form:"id_event" binding:"required"`
}

type ValidasiGetDataByNik struct {
	ValidasiKonfirmasiNik
	Month int `form:"month" binding:"required"`
	Year  int `form:"year" binding:"required"`
}

type ValidasiDeleteEventByID struct {
	Id         int    `form:"id" binding:"required"`
	Keterangan string `form:"keterangan"`
}

type ValidasiGetBookingRoom struct {
	IdRoom string `form:"id_room" binding:"required"`
	Date   string `form:"date" binding:"required"`
}

type ValidasiStoreNotulen struct {
	IdEvent   int    `form:"id_event" binding:"required"`
	Deskripsi string `form:"deskripsi" binding:"required"`
}
type ValidasiRenameFileNotulen struct {
	ValidasiStoreNotulen
	OldNameFile string `form:"old_name_file" binding:"required"`
	NewNameFile string `form:"new_name_file" binding:"required"`
}

type ValidasiDeleteFileNotulen struct {
	ValidasiStoreNotulen
	NameFile string `form:"name_file" binding:"required"`
}

type Event struct {
	EventID         int     `json:"event_id"`
	EventTitle      string  `json:"event_title"`
	EventDesc       string  `json:"event_desc"`
	EventStart      string  `json:"event_start"`
	EventEnd        string  `json:"event_end"`
	EventType       string  `json:"event_type"`
	EventURL        *string `json:"event_url"`
	EventImgName    *string `json:"event_img_name"`
	EventImgURL     *string `json:"event_img_url"`
	EventDate       string  `json:"event_date"`
	EventTimeStart  string  `json:"event_time_start"`
	EventTimeEnd    string  `json:"event_time_end"`
	IsAbsent        bool    `json:"is_absent"`
	CompCode        string  `json:"comp_code"`
	StatusKehadiran string  `json:"status_kehadiran"`
	AssetCompCode   string  `json:"asset_comp_code"`
}

type ListEventDelete struct {
	Model
	Person []Persons `json:"person"`
}

type Persons struct {
	events.EventPerson
	Dispose []Disposes `json:"dispose"`
}

type Disposes struct {
	events.EventPerson
	Profile *string `json:"profile"`
}

type EventDataByNik struct {
	EventData
	EventLocation *string        `json:"event_location"`
	BookRoom      *EventBookRoom `json:"book_room"`
}

type EventData struct {
	EventID        int     `json:"event_id"`
	EventTitle     string  `json:"event_title"`
	EventDesc      string  `json:"event_desc"`
	EventStart     string  `json:"event_start"`
	EventEnd       string  `json:"event_end"`
	EventType      string  `json:"event_type"`
	EventURL       *string `json:"event_url"`
	EventImgName   *string `json:"event_img_name"`
	EventImgURL    *string `json:"event_img_url"`
	EventDate      string  `json:"event_date"`
	EventTimeStart string  `json:"event_time_start"`
	EventTimeEnd   string  `json:"event_time_end"`
	EventStatus    string  `json:"event_status"`
	CompCode       string  `json:"comp_code"`
	AssetCompCode  string  `json:"asset_comp_code"`
}

type EventBookRoom struct {
	events.EventMsterRoom
	Companys pihc.PihcMasterCompany `json:"companys"`
}

type EventPersonDetail struct {
	Nik             string `json:"nik"`
	Nama            string `json:"nama"`
	DeptTitle       string `json:"dept_title"`
	Email           string `json:"email"`
	StatusKehadiran string `json:"status_kehadiran"`
	PhotoURL        string `json:"photo_url"`
}

type EventShowEvent struct {
	EventID                 int                      `json:"event_id"`
	EventTitle              string                   `json:"event_title"`
	EventDesc               string                   `json:"event_desc"`
	EventStart              string                   `json:"event_start"`
	EventDateStart          string                   `json:"event_date_start"`
	EventEnd                string                   `json:"event_end"`
	EventDateEnd            string                   `json:"event_date_end"`
	EventDateFormat         string                   `json:"event_date_format"`
	EventType               string                   `json:"event_type"`
	EventURL                *string                  `json:"event_url"`
	EventLocation           *string                  `json:"event_location"`
	EventImgName            *string                  `json:"event_img_name"`
	EventImgURL             *string                  `json:"event_img_url"`
	EventDate               string                   `json:"event_date"`
	EventTimeStart          string                   `json:"event_time_start"`
	EventTimeEnd            string                   `json:"event_time_end"`
	EventStatus             string                   `json:"event_status"`
	EventCreatedBy          *string                  `json:"event_created_by"`
	EventCreatedByNik       *string                  `json:"event_created_by_nik"`
	EventCreatedDeptTitle   *string                  `json:"event_created_dept_title"`
	EventApprovalPerson     *string                  `json:"event_approval_person"`
	EventRoom               *string                  `json:"event_room"`
	EventApprovalPersonName *string                  `json:"event_approval_person_name"`
	CompCode                string                   `json:"comp_code"`
	StatusKehadiran         string                   `json:"status_kehadiran"`
	Count                   events.EventCounts       `json:"count"`
	IsAbsent                bool                     `json:"is_absent"`
	AssestCompCode          string                   `json:"assest_comp_code"`
	Person                  []EventPersonDetail      `json:"person"`
	Notulen                 *GetDataNotulenFiles     `json:"notulen"`
	Materi                  []events.EventMateriFile `json:"materi"`
	BookRoom                *DataBookRoomShow        `json:"book_room"`
	TimeCreatedAt           string                   `json:"time_created_at"`
}

type DataBookRoomShow struct {
	IDBooking    int    `json:"id_booking"`
	CodeRoom     string `json:"code_room"`
	DateStart    string `json:"date_start"`
	DateEnd      string `json:"date_end"`
	TimeStart    string `json:"time_start"`
	TimeEnd      string `json:"time_end"`
	RoomID       string `json:"room_id"`
	RoomName     string `json:"room_name"`
	RoomCategory string `json:"room_category"`
	RoomCompCode string `json:"room_comp_code"`
	RoomCompName string `json:"room_comp_name"`
}

type DataBookRoomDate struct {
	IDBooking         int    `json:"id_booking"`
	NamaEvent         string `json:"nama_event"`
	NamaPembuat       string `json:"nama_pembuat"`
	KompatemenPembuat string `json:"kompatemen_pembuat"`
	DateStart         string `json:"date_start"`
	DateEnd           string `json:"date_end"`
	TimeStart         string `json:"time_start"`
	TimeEnd           string `json:"time_end"`
	Date              string `json:"date"`
}

type DataCreatedByEvent struct {
	DataCreatedBy pihc.PihcMasterKary
}

type DataApprovalPersonEvent struct {
	DataApprovalPerson pihc.PihcMasterKary
}

type GetDataNotulenFiles struct {
	IdNotulen int                       `json:"id_notulen" gorm:"primary_key"`
	IdEvent   int                       `json:"id_event"`
	Deskripsi string                    `json:"deskripsi" gorm:"default:null"`
	CreatedAt time.Time                 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time                 `json:"updated_at" gorm:"autoUpdateTime"`
	Files     []events.EventNotulenFile `json:"files"`
}

type EventPrintDaftarHadir struct {
	Title         string               `json:"title"`
	Deskripsi     string               `json:"deskripsi"`
	Type          string               `json:"type"`
	Tanggal       string               `json:"tanggal"`
	JamMulai      string               `json:"jam_mulai"`
	JamSelesai    string               `json:"jam_selesai"`
	EventURL      *string              `json:"event_url"`
	EventLocation *string              `json:"event_location"`
	CompIcon      string               `json:"comp_icon"`
	CreatedBy     string               `json:"created_by"`
	CreatedByName *string              `json:"created_by_name"`
	CreatedByDept *string              `json:"created_by_dept"`
	BookRoom      *BookRoomDaftarHadir `json:"book_room"`
	Presence      []Presences          `json:"presence"`
}

type BookRoomDaftarHadir struct {
	RoomName     string `json:"room_name"`
	RoomCategory string `json:"room_category"`
	RoomCompName string `json:"room_comp_name"`
}

type Presences struct {
	Nomer    int         `json:"nomer"`
	Nik      string      `json:"nik"`
	Nama     null.String `json:"nama"`
	Email    null.String `json:"email"`
	Jabatan  null.String `json:"jabatan"`
	NoTelp   null.String `json:"no_telp"`
	Dept     string      `json:"dept"`
	Time     string      `json:"time"`
	Date     string      `json:"date"`
	Datetime string      `json:"datetime"`
	Instansi null.String `json:"instansi"`
}

type DataInFeed struct {
	Type        string `json:"type"`
	*Model      `json:"model"`
	JumlahEvent int `json:"jumlah_event"`
}

type Model struct {
	ID              int       `json:"id"`
	EventTitle      string    `json:"event_title"`
	EventDesc       string    `json:"event_desc"`
	EventStart      string    `json:"event_start"`
	EventEnd        string    `json:"event_end"`
	EventType       string    `json:"event_type"`
	EventURL        *string   `json:"event_url"`
	EventImgName    *string   `json:"event_img_name"`
	EventImgURL     *string   `json:"event_img_url"`
	CompCode        string    `json:"comp_code"`
	Status          string    `json:"status"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ApprovalPerson  *string   `json:"approval_person"`
	EventRoom       *string   `json:"event_room"`
	EventLocation   *string   `json:"event_location"`
	EventKeterangan *string   `json:"event_keterangan"`
}
