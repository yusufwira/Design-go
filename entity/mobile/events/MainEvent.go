package events

import (
	"time"

	"gorm.io/gorm"
)

type MainEvent struct {
	Id              int       `json:"id" gorm:"primary_key"`
	EventTitle      string    `json:"event_title"`
	EventDesc       string    `json:"event_desc"`
	EventStart      time.Time `json:"event_start"`
	EventEnd        time.Time `json:"event_end"`
	EventType       string    `json:"event_type"`
	EventUrl        *string   `json:"event_url" gorm:"default:null"`
	EventImgName    *string   `json:"event_img_name" gorm:"default:null"`
	EventImgUrl     *string   `json:"event_img_url" gorm:"default:null"`
	CompCode        string    `json:"comp_code" gorm:"default:null"`
	Status          string    `json:"status"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	ApprovalPerson  *string   `json:"approval_person" gorm:"default:null"`
	EventRoom       *string   `json:"event_room" gorm:"default:null"`
	EventLocation   *string   `json:"event_location" gorm:"default:null"`
	EventKeterangan *string   `json:"event_keterangan" gorm:"default:null"`
}

type EventHistory struct {
	EventTitle      string    `json:"event_title"`
	EventDesc       string    `json:"event_desc"`
	EventStart      time.Time `json:"event_start"`
	EventEnd        time.Time `json:"event_end"`
	EventType       string    `json:"event_type"`
	EventURL        string    `json:"event_url" gorm:"default:null"`
	EventImgName    string    `json:"event_img_name" gorm:"default:null"`
	EventImgURL     string    `json:"event_img_url" gorm:"default:null"`
	CompCode        string    `json:"comp_code"`
	Status          string    `json:"status"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime;default:null"`
	ApprovalPerson  string    `json:"approval_person" gorm:"default:null"`
	EventRoom       string    `json:"event_room" gorm:"default:null"`
	EventLocation   string    `json:"event_location" gorm:"default:null"`
	EventKeterangan string    `json:"event_keterangan" gorm:"default:null"`
}

func (MainEvent) TableName() string {
	return "mobile.event"
}

type MainEventRepo struct {
	DB *gorm.DB
}

func NewMainEventRepo(db *gorm.DB) *MainEventRepo {
	return &MainEventRepo{DB: db}
}

func (t MainEventRepo) Create(me MainEvent) (MainEvent, error) {
	err := t.DB.Create(&me).Error
	if err != nil {
		return me, err
	}
	return me, nil
}

func (t MainEventRepo) History(eh EventHistory) (EventHistory, error) {
	err := t.DB.Table("mobile.event_hist").Create(&eh).Error
	if err != nil {
		return eh, err
	}
	return eh, nil
}

func (t MainEventRepo) FindEventMainID(id int) (MainEvent, error) {
	var main_event MainEvent
	err := t.DB.Where("id=?", id).Take(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}
func (t MainEventRepo) FindEventMainDrafted(stats string) ([]MainEvent, error) {
	var main_event []MainEvent
	err := t.DB.Where("status=?", stats).Find(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}

func (t MainEventRepo) FindEventMainIDType(Id int, Type string) (MainEvent, error) {
	var main_event MainEvent
	err := t.DB.Where("id=? AND event_type=?", Id, Type).Take(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}

func (t MainEventRepo) FindEventMainIDNIK(id int, nik string) (MainEvent, error) {
	var main_event MainEvent
	err := t.DB.Where("id=? AND created_by=?", id, nik).Take(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}

func (t MainEventRepo) Update(me MainEvent) (MainEvent, error) {
	err := t.DB.Save(&me).Error
	if err != nil {
		return me, err
	}
	return me, nil
}

func (t MainEventRepo) FindEventMainNikCreatedBy(nik string) ([]MainEvent, error) {
	var main_event []MainEvent
	err := t.DB.Where("created_by=?", nik).Find(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}

func (t MainEventRepo) FindEventMainNikApprovalPerson(nik string) ([]MainEvent, error) {
	var main_event []MainEvent
	err := t.DB.Where("approval_person=? AND status='Drafted'", nik).Find(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}

func (t MainEventRepo) FindEventMainNikMonthYear(nik string, month int, year int, status string) ([]MainEvent, error) {
	var main_event []MainEvent
	err := t.DB.Where("created_by=? AND DATE_PART('month', event_start)=? AND DATE_PART('year', event_start)=? AND status!=?", nik, month, year, status).Order("id asc").Find(&main_event).Error
	if err != nil {
		return main_event, err
	}
	return main_event, nil
}

func (t MainEventRepo) DeleteMainEvent(id int) (MainEvent, error) {
	var main_event MainEvent
	err := t.DB.Where("id=?", id).Take(&main_event).Error
	if err == nil {
		t.DB.Where("id= ?", id).Delete(&main_event)
		return main_event, nil
	}
	return main_event, err
}

func (t MainEventRepo) FindDataInFeed(nik string, typeHari string) (MainEvent, int64, error) {
	var main_event MainEvent
	var jmlah_event int64

	if typeHari == "Hari Ini" {
		err_ev := t.DB.Where("date(event_start) = current_date and status != 'Drafted' and created_by=?", nik).Order("event_start asc").Take(&main_event).Error
		t.DB.Table("mobile.event").Where("date(event_start) = current_date and status != 'Drafted' and created_by=?", nik).Count(&jmlah_event)
		if err_ev != nil {
			return main_event, jmlah_event, err_ev
		}
	}
	if typeHari == "Besok" {
		err_ev := t.DB.Where("date(event_start) = current_date + interval '1' day and status != 'Drafted' and created_by=?", nik).Order("event_start asc").Take(&main_event).Error
		t.DB.Table("mobile.event").Where("date(event_start) = current_date + interval '1' day and status != 'Drafted' and created_by=?", nik).Count(&jmlah_event)
		if err_ev != nil {
			return main_event, jmlah_event, err_ev
		}
	}
	if typeHari == "Lusa" {
		err_ev := t.DB.Where("date(event_start) = current_date + interval '2' day and status != 'Drafted' and created_by=?", nik).Order("event_start asc").Take(&main_event).Error
		t.DB.Table("mobile.event").Where("date(event_start) = current_date + interval '2' day and status != 'Drafted' and created_by=?", nik).Count(&jmlah_event)
		if err_ev != nil {
			return main_event, jmlah_event, err_ev
		}
	}
	return main_event, jmlah_event, nil
}

// func (t EventPresenceRepo) FindPresenceIDNIK(id int, nik string) (bool, error) {
// 	var ev_rb int64
// 	err := t.DB.Table("mobile.event_presence").Where("id_event=? AND emp_no=?", id, nik).Count(&ev_rb).Error
// 	if err != nil {
// 		return true, err
// 	}
// 	return ev_rb != 0, nil
// }
// func (t EventPersonRepo) GetEventCounts(id int) (EventCounts, error) {
// 	var counts EventCounts

// 	err := t.DB.Raw(`
//         SELECT
//             COUNT(*) FILTER (WHERE status_kehadiran = 'hadir') AS count_hadir,
//             COUNT(*) FILTER (WHERE status_kehadiran = 'menunggu') AS count_menunggu,
//             COUNT(*) FILTER (WHERE status_kehadiran = 'tidak_hadir') AS count_tidak_hadir
//         FROM mobile.event_person
//         WHERE id_event = ? AND id_parent is NULL`, id).Scan(&counts).Error

// 	if err != nil {
// 		return counts, err
// 	}

// 	counts.CountGuest = counts.CountHadir + counts.CountMenunggu + counts.CountTidakHadir

// 	return counts, nil
// }
