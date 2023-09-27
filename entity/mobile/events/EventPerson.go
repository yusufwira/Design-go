package events

import (
	"time"

	"gorm.io/gorm"
)

type EventPerson struct {
	Id                int       `json:"id" gorm:"primary_key"`
	IdEvent           int       `json:"id_event"`
	IdParent          *int      `json:"id_parent" gorm:"default:null"`
	Nik               string    `json:"nik"`
	StatusKehadiran   string    `json:"status_kehadiran"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	KetKetidakhadiran *string   `json:"ket_ketidakhadiran" gorm:"default:null"`
}

type EventCounts struct {
	CountGuest      int64 `json:"count_guest"`
	CountHadir      int64 `json:"count_hadir"`
	CountMenunggu   int64 `json:"count_menunggu"`
	CountTidakHadir int64 `json:"count_tidak_hadir"`
}

type EventDetailEventPerson struct {
	Nik             string `json:"nik"`
	Nama            string `json:"nama"`
	DeptTitle       string `json:"dept_title"`
	Email           string `json:"email"`
	StatusKehadiran string `json:"status_kehadiran"`
	PhotoURL        string `json:"photo_url"`
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

func (t EventPersonRepo) Update(me EventPerson) (EventPerson, error) {
	err := t.DB.Save(&me).Error
	if err != nil {
		return me, err
	}
	return me, nil
}

func (t EventPersonRepo) FindEventPersonID(id int) []EventPerson {
	var event_person []EventPerson

	t.DB.Where("id_event=?", id).Find(&event_person)

	if len(event_person) != 0 {
		return event_person
	}

	return nil
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
        WHERE id_event = ? AND id_parent is NULL`, id).Scan(&counts).Error

	if err != nil {
		return counts, err
	}

	counts.CountGuest = counts.CountHadir + counts.CountMenunggu + counts.CountTidakHadir

	return counts, nil
}

func (t EventPersonRepo) FindDetailEventPerson(id int) ([]EventDetailEventPerson, error) {
	var detailEventPerson []EventDetailEventPerson

	err := t.DB.Raw(`
	select ep.nik as nik, pmk.nama as nama, pmk.dept_title as dept_title, pmk.email as email,
		ep.status_kehadiran as status_kehadiran, pp.url as photo_url
		from mobile.event_person ep 
		left join dbo.pihc_master_karyawan pmk on pmk.emp_no  = ep.nik
		left join mobile.profile_photo pp on pp.emp_no = ep.nik
	where ep.id_event = ? and (ep.id_parent is null or ep.id_parent = 0)`, id).Scan(&detailEventPerson).Error

	if err != nil {
		return detailEventPerson, err
	}

	return detailEventPerson, nil
}

func (t EventPersonRepo) GetDataDisposePerson(nik string, idEvent int) ([]EventDetailEventPerson, error) {
	var detailEventPerson []EventDetailEventPerson

	err := t.DB.Raw(`
	select ep.nik as nik, pmk.nama as nama, pmk.dept_title as dept_title, pmk.email as email,
		ep.status_kehadiran as status_kehadiran, pp.url as photo_url
		from mobile.event_person ep 
		left join dbo.pihc_master_karyawan pmk on pmk.emp_no  = ep.nik
		left join mobile.profile_photo pp on pp.emp_no = ep.nik
	where id_parent IN (
		SELECT ep1.nik::bigint
		FROM mobile.event_person ep1
		WHERE ep1.nik = ? and ep1.id_event = ?
	) and id_event = ? and pmk.emp_no is not null and pmk.dept_title is not null;`, nik, idEvent, idEvent).Scan(&detailEventPerson).Error

	if err != nil {
		return nil, err
	}

	return detailEventPerson, nil
}

func (t EventPersonRepo) DeleteEventPerson(id int) ([]EventPerson, []EventPerson, error) {
	var ev_person EventPerson
	var ev_person_isnotparent []EventPerson
	var ev_person_parent []EventPerson
	var data_isnotParent []EventPerson
	var data_isParent []EventPerson
	var list_nik []string

	err := t.DB.Where("id_event=? AND id_parent is NULL", id).Find(&ev_person_isnotparent).Error
	if err == nil {
		for _, dataPerson := range ev_person_isnotparent {
			list_nik = append(list_nik, dataPerson.Nik)
		}
		data_isnotParent = ev_person_isnotparent

		t.DB.Where("id_event=? AND id_parent in (?)", id, list_nik).Find(&ev_person_parent)
		data_isParent = ev_person_parent

		t.DB.Where("id_event= ?", id).Delete(&ev_person)

		return data_isnotParent, data_isParent, nil
	}
	return nil, nil, err
}
