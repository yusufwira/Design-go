package event_controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/mobile/events"
	"gorm.io/gorm"
)

type EventController struct {
	EventBookingRoomRepo  *events.EventBookingRoomRepo
	EventPersonRepo       *events.EventPersonRepo
	MainEventRepo         *events.MainEventRepo
	EventPresenceRepo     *events.EventPresenceRepo
	EventMsterRoomRepo    *events.EventMsterRoomRepo
	PihcMasterKaryRepo    *pihc.PihcMasterKaryRepo
	PihcMasterKaryRtRepo  *pihc.PihcMasterKaryRtRepo
	PihcMasterCompanyRepo *pihc.PihcMasterCompanyRepo
}

func NewEventController(db *gorm.DB) *EventController {
	return &EventController{EventBookingRoomRepo: events.NewEventBookingRoomRepo(db),
		EventPersonRepo:       events.NewEventPersonRepo(db),
		MainEventRepo:         events.NewMainEventRepo(db),
		EventPresenceRepo:     events.NewEventPresenceRepo(db),
		EventMsterRoomRepo:    events.NewEventMsterRoomRepo(db),
		PihcMasterKaryRepo:    pihc.NewPihcMasterKaryRepo(db),
		PihcMasterKaryRtRepo:  pihc.NewPihcMasterKaryRtRepo(db),
		PihcMasterCompanyRepo: pihc.NewPihcMasterCompanyRepo(db)}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return (fe.Field() + " wajib di isi")
	case "validyear":
		return ("Field has an invalid value: " + fe.Field() + fe.Tag())
	}
	return "Unknown error"
}

func (c *EventController) StoreEvent(ctx *gin.Context) {
	var req Authentication.ValidasiEvent
	var ev_br events.EventBookingRoom
	var ev_person events.EventPerson
	var ev_main events.MainEvent

	if err := ctx.ShouldBind(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]Authentication.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = Authentication.ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": out})
		}
		return
	}

	PIHC_MSTR_KRY_RT, _ := c.PihcMasterKaryRtRepo.FindUserByNIK(strconv.Itoa(req.Nik))

	comp_code := PIHC_MSTR_KRY_RT.Company

	parsedTimeStart, _ := time.Parse(time.DateTime, req.Start)
	parsedTimeEnd, _ := time.Parse(time.DateTime, req.End)

	if req.ID != 0 {
		main_event, _ := c.MainEventRepo.FindEventMainID(req.ID)

		main_event.EventTitle = req.Title
		main_event.EventDesc = req.Desc
		main_event.EventType = req.Type
		main_event.EventStart = parsedTimeStart
		main_event.EventEnd = parsedTimeEnd
		main_event.Status = req.Status
		main_event.CreatedBy = strconv.Itoa(req.Nik)
		main_event.CompCode = comp_code
		main_event.EventLocation = req.Location

		if req.Type == "online" {
			// Type Online

			if main_event.EventRoom != "" {
				c.EventBookingRoomRepo.DeleteRoomBooking(main_event.Id)
				main_event.EventRoom = ""
			}

			main_event.EventUrl = req.URL
			// main_event.ApprovalPerson = "7222322"

			eventMain, err_eventMain := c.MainEventRepo.Update(main_event)

			if err_eventMain == nil {
				var list_id_person []int

				// Remove the "[" and "]" characters from the string
				trimmed := strings.Trim(req.Person, "[]")

				// Split the string into individual values
				values := strings.Split(trimmed, ",")

				for _, nikPerson := range values {
					ev_person.IdEvent = eventMain.Id
					ev_person.Nik = nikPerson
					ev_person.StatusKehadiran = "menunggu"
					eventPerson, _ := c.EventPersonRepo.Create(ev_person)
					list_id_person = append(list_id_person, eventPerson.Id)
				}

				c.EventPersonRepo.DelParticipationLama(eventMain.Id, list_id_person)

			}
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
				"data":    "Data berhasil diUpdate",
			})
		} else if req.Type == "offline" {
			// Type Offline
			if main_event.EventUrl != "" {
				main_event.EventUrl = ""
			}

			eventBookingRoom, errEventBookingRoom := c.EventBookingRoomRepo.FindRoomBooking(main_event.EventRoom, main_event.Id)

			if errEventBookingRoom == nil {

				existRoom, _ := c.EventBookingRoomRepo.FindExistRoom(req.IDRoom, parsedTimeStart, parsedTimeEnd)

				if !existRoom {
					main_event.EventRoom = req.IDRoom
					// main_event.ApprovalPerson = "7222322"

					eventMain, err_eventMain := c.MainEventRepo.Update(main_event)

					if err_eventMain == nil {
						var list_id_person []int
						eventBookingRoom.CodeRoom = eventMain.EventRoom
						eventBookingRoom.DateStart = eventMain.EventStart
						eventBookingRoom.DateEnd = eventMain.EventEnd

						c.EventBookingRoomRepo.Update(eventBookingRoom)

						// Remove the "[" and "]" characters from the string
						trimmed := strings.Trim(req.Person, "[]")

						// Split the string into individual values
						values := strings.Split(trimmed, ",")

						for _, nikPerson := range values {
							ev_person.IdEvent = eventMain.Id
							ev_person.Nik = nikPerson
							ev_person.StatusKehadiran = "menunggu"
							eventPerson, _ := c.EventPersonRepo.Create(ev_person)
							list_id_person = append(list_id_person, eventPerson.Id)
						}
						c.EventPersonRepo.DelParticipationLama(eventMain.Id, list_id_person)

						ctx.JSON(http.StatusOK, gin.H{
							"status":  http.StatusOK,
							"success": "Success",
							"data":    "Data berhasil diUpdate",
						})
					}
				}
			}
		} else if req.Type == "hybrid" {
			// Type Hybrid

			existRoom, _ := c.EventBookingRoomRepo.FindExistRoom(req.IDRoom, parsedTimeStart, parsedTimeEnd)

			eventBookingRoom, errEventBookingRoom := c.EventBookingRoomRepo.FindRoomBooking(main_event.EventRoom, main_event.Id)

			if errEventBookingRoom == nil {
				if !existRoom {
					main_event.EventRoom = req.IDRoom
					main_event.EventUrl = req.URL
					// main_event.ApprovalPerson = "7222322"

					eventMain, err_eventMain := c.MainEventRepo.Update(main_event)

					if err_eventMain == nil {
						var list_id_person []int
						eventBookingRoom.CodeRoom = eventMain.EventRoom
						eventBookingRoom.DateStart = eventMain.EventStart
						eventBookingRoom.DateEnd = eventMain.EventEnd

						c.EventBookingRoomRepo.Update(eventBookingRoom)

						// Remove the "[" and "]" characters from the string
						trimmed := strings.Trim(req.Person, "[]")

						// Split the string into individual values
						values := strings.Split(trimmed, ",")

						for _, nikPerson := range values {
							ev_person.IdEvent = eventMain.Id
							ev_person.Nik = nikPerson
							ev_person.StatusKehadiran = "menunggu"
							eventPerson, _ := c.EventPersonRepo.Create(ev_person)
							list_id_person = append(list_id_person, eventPerson.Id)
						}

						c.EventPersonRepo.DelParticipationLama(eventMain.Id, list_id_person)

						ctx.JSON(http.StatusOK, gin.H{
							"status":  http.StatusOK,
							"success": "Success",
							"data":    "Data berhasil diUpdate",
						})
					}
				}
			}
		}
	} else {
		ev_main.EventTitle = req.Title
		ev_main.EventDesc = req.Desc
		ev_main.EventType = req.Type
		ev_main.EventStart = parsedTimeStart
		ev_main.EventEnd = parsedTimeEnd
		ev_main.Status = req.Status
		ev_main.CreatedBy = strconv.Itoa(req.Nik)
		ev_main.CompCode = comp_code
		ev_main.EventLocation = req.Location

		if req.Type == "online" {
			// Type Online
			ev_main.EventUrl = req.URL
			// main_event.ApprovalPerson = "7222322"

			event_main, _ := c.MainEventRepo.Create(ev_main)

			// Remove the "[" and "]" characters from the string
			trimmed := strings.Trim(req.Person, "[]")

			// Split the string into individual values
			values := strings.Split(trimmed, ",")

			for _, nikPerson := range values {
				ev_person.IdEvent = event_main.Id
				ev_person.Nik = nikPerson
				ev_person.StatusKehadiran = "menunggu"
				c.EventPersonRepo.Create(ev_person)
			}
		} else if req.Type == "offline" {
			// Type Offline
			existRoom, _ := c.EventBookingRoomRepo.FindExistRoom(req.IDRoom, parsedTimeStart, parsedTimeEnd)

			if !existRoom {
				ev_main.EventRoom = req.IDRoom
				// main_event.ApprovalPerson = "7222322"

				event_main, _ := c.MainEventRepo.Create(ev_main)

				ev_br.CodeRoom = event_main.EventRoom
				ev_br.IdEvent = event_main.Id
				ev_br.DateStart = event_main.EventStart
				ev_br.DateEnd = event_main.EventEnd

				c.EventBookingRoomRepo.Create(ev_br)

				// Remove the "[" and "]" characters from the string
				trimmed := strings.Trim(req.Person, "[]")

				// Split the string into individual values
				values := strings.Split(trimmed, ",")

				for _, nikPerson := range values {
					ev_person.IdEvent = event_main.Id
					ev_person.Nik = nikPerson
					ev_person.StatusKehadiran = "menunggu"
					c.EventPersonRepo.Create(ev_person)
				}
			}
		} else if req.Type == "hybrid" {
			// Type Hybrid
			existRoom, _ := c.EventBookingRoomRepo.FindExistRoom(req.IDRoom, parsedTimeStart, parsedTimeEnd)

			if !existRoom {
				ev_main.EventRoom = req.IDRoom
				ev_main.EventUrl = req.URL
				// main_event.ApprovalPerson = "7222322"

				event_main, _ := c.MainEventRepo.Create(ev_main)

				ev_br.CodeRoom = event_main.EventRoom
				ev_br.IdEvent = event_main.Id
				ev_br.DateStart = event_main.EventStart
				ev_br.DateEnd = event_main.EventEnd

				c.EventBookingRoomRepo.Create(ev_br)

				// Remove the "[" and "]" characters from the string
				trimmed := strings.Trim(req.Person, "[]")

				// Split the string into individual values
				values := strings.Split(trimmed, ",")

				for _, nikPerson := range values {
					ev_person.IdEvent = event_main.Id
					ev_person.Nik = nikPerson
					ev_person.StatusKehadiran = "menunggu"
					c.EventPersonRepo.Create(ev_person)
				}
			}
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
	})
}

func (c *EventController) UpdateStatusEvent(ctx *gin.Context) {
	var req Authentication.ValidasiUpdateStatusEvent

	if err := ctx.ShouldBind(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]Authentication.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = Authentication.ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": out})
		}
		return
	}

	main_event, err := c.MainEventRepo.FindEventMainID(req.EventID)

	if err == nil {
		main_event.ApprovalPerson = req.Nik
		main_event.Status = req.Status
		main_event.EventKeterangan = req.Keterangan

		c.MainEventRepo.Update(main_event)

		if req.Status == "Decline" {
			if main_event.EventRoom != "" {
				c.EventBookingRoomRepo.DeleteRoomBooking(main_event.Id)
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *EventController) StoreDispose(ctx *gin.Context) {
	var req Authentication.ValidasiStoreDispose
	var ev_person events.EventPerson

	if err := ctx.ShouldBind(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]Authentication.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = Authentication.ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": out})
		}
		return
	}

	_, errPerson := c.EventPersonRepo.FindEventPersonID(req.EventID)

	if errPerson == nil {
		// Remove the "[" and "]" characters from the string
		trimmed := strings.Trim(req.Dispose, "[]")

		// Split the string into individual values
		values := strings.Split(trimmed, ",")

		for _, nikPerson := range values {
			ev_person.IdEvent = req.EventID
			ev_person.Nik = nikPerson
			ev_person.IdParent = req.Nik
			ev_person.StatusKehadiran = "menunggu"
			c.EventPersonRepo.Create(ev_person)
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
		})
	}
}

func (c *EventController) GetDataEvent(ctx *gin.Context) {
	nik := ctx.Param("nik")
	var ev []Authentication.EventDataByNik

	main_event, _ := c.MainEventRepo.FindEventMainNik(nik)

	data_list := make([]Authentication.EventDataByNik, len(main_event))
	for i, DataEvent := range main_event {
		data_list[i] = Authentication.EventDataByNik{
			EventID:        DataEvent.Id,
			EventTitle:     DataEvent.EventTitle,
			EventDesc:      DataEvent.EventDesc,
			EventStart:     DataEvent.EventStart.Format("2006-01-02 15:04:05"),
			EventEnd:       DataEvent.EventEnd.Format("2006-01-02 15:04:05"),
			EventType:      DataEvent.EventType,
			EventURL:       DataEvent.EventUrl,
			EventImgName:   DataEvent.EventImgName,
			EventImgURL:    DataEvent.EventImgUrl,
			EventDate:      DataEvent.EventStart.Format("2006-01-02"),
			EventTimeStart: DataEvent.EventStart.Format("15:04:05"),
			EventTimeEnd:   DataEvent.EventEnd.Format("15:04:05"),
			EventStatus:    DataEvent.Status,
			EventLocation:  DataEvent.EventLocation,
			CompCode:       DataEvent.CompCode,
			AssetCompCode:  "https://pismart-dev.pupuk-indonesia.com/public/assets/media/logos/logo-pi-full.png",
		}
		if DataEvent.EventType == "offline" || DataEvent.EventType == "hybrid" {
			ev_rb, err := c.EventBookingRoomRepo.FindRoomBooking(DataEvent.EventRoom, DataEvent.Id)
			if err == nil {
				ev_mstr_room, err := c.EventMsterRoomRepo.FindEventMasterRoom(ev_rb.CodeRoom)
				if err != nil {
					data_list[i].BookRoom = nil
				} else {
					pihc_mster_company, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(ev_mstr_room.CompCode)
					data_list[i].BookRoom = &Authentication.EventBookRoom{
						EventMsterRoom: ev_mstr_room,
						Companys:       pihc_mster_company,
					}
				}
			}
		} else {
			data_list[i].BookRoom = nil
		}
	}
	ev = data_list

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    ev,
	})
}

func (c *EventController) GetDataByNik(ctx *gin.Context) {
	var req Authentication.ValidasiGetDataByNik
	var ev []Authentication.Event

	if err := ctx.ShouldBind(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]Authentication.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = Authentication.ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": out})
		}
		return
	}
	// fmt.Println(time.Month(req.Month))
	status := "Drafted"
	main_event, _ := c.MainEventRepo.FindEventMainNikMonthYear(req.Nik, req.Month, req.Year, status)

	data_list := make([]Authentication.Event, len(main_event))
	for i, DataEvent := range main_event {
		ev_person, _ := c.EventPersonRepo.FindEventPersonIDNIK(DataEvent.Id, DataEvent.CreatedBy)
		is_absent, _ := c.EventPresenceRepo.FindPresenceIDNIK(DataEvent.Id, DataEvent.CreatedBy)

		data_list[i] = Authentication.Event{
			EventID:         DataEvent.Id,
			EventTitle:      DataEvent.EventTitle,
			EventDesc:       DataEvent.EventDesc,
			EventStart:      DataEvent.EventStart.Format("2006-01-02 15:04:05"),
			EventEnd:        DataEvent.EventEnd.Format("2006-01-02 15:04:05"),
			EventURL:        DataEvent.EventUrl,
			EventImgName:    DataEvent.EventImgName,
			EventImgURL:     DataEvent.EventImgUrl,
			EventDate:       DataEvent.EventStart.Format("2006-01-02"),
			EventTimeStart:  DataEvent.EventStart.Format("15:04:05"),
			EventTimeEnd:    DataEvent.EventEnd.Format("15:04:05"),
			IsAbsent:        is_absent,
			CompCode:        DataEvent.CompCode,
			StatusKehadiran: ev_person.StatusKehadiran,
			AssetCompCode:   "https://pismart-dev.pupuk-indonesia.com/public/assets/media/logos/logo-pi-full.png",
		}
	}
	ev = data_list

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    ev,
	})
}

func (c *EventController) ShowEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	nik := ctx.Param("nik")
	idEvent, _ := strconv.Atoi(id)

	mainEvent, err := c.MainEventRepo.FindEventMainIDNIK(idEvent, nik)
	if err == nil {
		list_nik := []string{mainEvent.CreatedBy, mainEvent.ApprovalPerson}
		pihcCreatedByApprovalperson, _ := c.PihcMasterKaryRepo.FindUserByNIKArray(list_nik)
		ctx.JSON(http.StatusOK, gin.H{
			"data": pihcCreatedByApprovalperson,
		})
		eventCount, _ := c.EventPersonRepo.GetEventCounts(mainEvent.Id)
		fmt.Println(eventCount.CountGuest, eventCount.CountHadir, eventCount.CountMenunggu, eventCount.CountTidakHadir)
		// is_absent, _ := c.EventPresenceRepo.FindPresenceIDNIK(DataEvent.Id, DataEvent.CreatedBy)
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"succes": "Data Not Found",
		})
	}
}

func (c *EventController) DeleteEvent(ctx *gin.Context) {
	var req Authentication.ValidasiDeleteEventByID
	// var ev []Authentication.ListEventDelete

	if err := ctx.ShouldBind(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]Authentication.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = Authentication.ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": out})
		}
		return
	}
	// // fmt.Println(time.Month(req.Month))
	// status := "Drafted"
	// main_event, _ := c.MainEventRepo.FindEventMainNikMonthYear(req.Nik, req.Month, req.Year, status)

	// data_list := make([]Authentication.Event, len(main_event))
	// for i, DataEvent := range main_event {
	// 	ev_person, _ := c.EventPersonRepo.FindEventPersonIDNIK(DataEvent.Id, DataEvent.CreatedBy)
	// 	is_absent, _ := c.EventPresenceRepo.FindPresenceIDNIK(DataEvent.Id, DataEvent.CreatedBy)

	// 	data_list[i] = Authentication.Event{
	// 		EventID:         DataEvent.Id,
	// 		EventTitle:      DataEvent.EventTitle,
	// 		EventDesc:       DataEvent.EventDesc,
	// 		EventStart:      DataEvent.EventStart.Format("2006-01-02 15:04:05"),
	// 		EventEnd:        DataEvent.EventEnd.Format("2006-01-02 15:04:05"),
	// 		EventURL:        DataEvent.EventUrl,
	// 		EventImgName:    DataEvent.EventImgName,
	// 		EventImgURL:     DataEvent.EventImgUrl,
	// 		EventDate:       DataEvent.EventStart.Format("2006-01-02"),
	// 		EventTimeStart:  DataEvent.EventStart.Format("15:04:05"),
	// 		EventTimeEnd:    DataEvent.EventEnd.Format("15:04:05"),
	// 		IsAbsent:        is_absent,
	// 		CompCode:        DataEvent.CompCode,
	// 		StatusKehadiran: ev_person.StatusKehadiran,
	// 		AssetCompCode:   "https://pismart-dev.pupuk-indonesia.com/public/assets/media/logos/logo-pi-full.png",
	// 	}
	// }
	// ev = data_list

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"status":  http.StatusOK,
	// 	"success": "Success",
	// 	"data":    ev,
	// })
}
