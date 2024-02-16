package event_controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/mobile/events"
	"gorm.io/gorm"
)

type EventController struct {
	EventBookingRoomRepo   *events.EventBookingRoomRepo
	EventPersonRepo        *events.EventPersonRepo
	MainEventRepo          *events.MainEventRepo
	EventPresenceRepo      *events.EventPresenceRepo
	EventMsterRoomRepo     *events.EventMsterRoomRepo
	EventNotulenRepo       *events.EventNotulenRepo
	EventMateriFileRepo    *events.EventMateriFileRepo
	PihcMasterKaryDbRepo   *pihc.PihcMasterKaryDbRepo
	PihcMasterKaryRtDbRepo *pihc.PihcMasterKaryRtDbRepo
	PihcMasterCompanyRepo  *pihc.PihcMasterCompanyRepo
}

func NewEventController(Db *gorm.DB, StorageClient *storage.Client) *EventController {
	return &EventController{
		EventBookingRoomRepo:   events.NewEventBookingRoomRepo(Db),
		EventPersonRepo:        events.NewEventPersonRepo(Db),
		MainEventRepo:          events.NewMainEventRepo(Db),
		EventPresenceRepo:      events.NewEventPresenceRepo(Db),
		EventMsterRoomRepo:     events.NewEventMsterRoomRepo(Db),
		EventNotulenRepo:       events.NewEventNotulenRepo(Db, StorageClient),
		EventMateriFileRepo:    events.NewEventMateriFileRepo(Db),
		PihcMasterKaryDbRepo:   pihc.NewPihcMasterKaryDbRepo(Db),
		PihcMasterKaryRtDbRepo: pihc.NewPihcMasterKaryRtDbRepo(Db),
		PihcMasterCompanyRepo:  pihc.NewPihcMasterCompanyRepo(Db)}
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
	var ev_materi_file events.EventMateriFile

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

	PIHC_MSTR_KRY_RT, _ := c.PihcMasterKaryRtDbRepo.FindUserByNIK(strconv.Itoa(req.Nik))

	comp_code := PIHC_MSTR_KRY_RT.Company

	// parsedTimeStart, _ := time.Parse(time.DateTime, req.Start)
	// parsedTimeEnd, _ := time.Parse(time.DateTime, req.End)

	if req.ID != 0 {
		main_event, _ := c.MainEventRepo.FindEventMainID(req.ID)

		if req.FileMateri != nil {
			var list_id_file_materi []int
			for _, files_materi := range req.FileMateri {
				ev_materi_file.IdEvent = main_event.Id
				ev_materi_file.FileName = files_materi.FileName
				ev_materi_file.FileUrl = files_materi.FileURL
				materi_file, _ := c.EventMateriFileRepo.Create(ev_materi_file)
				list_id_file_materi = append(list_id_file_materi, materi_file.IdMateriFile)
			}
			c.EventMateriFileRepo.DelMateriFileLama(main_event.Id, list_id_file_materi)
		}

		main_event.EventTitle = req.Title
		main_event.EventDesc = req.Desc
		main_event.EventType = req.Type
		main_event.EventStart, _ = time.Parse(time.DateTime, req.Start)
		main_event.EventEnd, _ = time.Parse(time.DateTime, req.End)
		main_event.Status = req.Status
		main_event.CreatedBy = strconv.Itoa(req.Nik)
		main_event.CompCode = comp_code
		main_event.EventLocation = req.Location

		if req.ApprovalPerson != nil {
			main_event.ApprovalPerson = req.ApprovalPerson
		}

		if req.Type == "online" {
			// Type Online

			if main_event.EventRoom != nil {
				c.EventBookingRoomRepo.DeleteRoomBooking(main_event.Id)
				main_event.EventRoom = nil
			}

			main_event.EventUrl = req.URL

			eventMain, err_eventMain := c.MainEventRepo.Update(main_event)

			if err_eventMain == nil {
				var list_id_person []int

				if req.IsPublic == 0 {
					// Remove the "[" and "]" characters from the string
					trimmed := strings.Trim(*req.Person, "[]")

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
			}
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
				"data":    "Data berhasil diUpdate",
			})
			return
		} else if req.Type == "offline" {
			// Type Offline
			if main_event.EventUrl != nil {
				main_event.EventUrl = nil
			}

			eventBookingRoom, errEventBookingRoom := c.EventBookingRoomRepo.FindRoomBooking(req.IDRoom, main_event.Id)

			if errEventBookingRoom == nil {

				existRoom, _ := c.EventBookingRoomRepo.FindExistRoom(*req.IDRoom, main_event.Id, main_event.EventStart, main_event.EventEnd)
				fmt.Println(existRoom)

				if !existRoom {
					main_event.EventRoom = req.IDRoom
					// main_event.ApprovalPerson = "7222322"

					eventMain, err_eventMain := c.MainEventRepo.Update(main_event)

					if err_eventMain == nil {
						var list_id_person []int
						eventBookingRoom.IdEvent = eventMain.Id
						eventBookingRoom.CodeRoom = eventMain.EventRoom
						eventBookingRoom.DateStart = eventMain.EventStart
						eventBookingRoom.DateEnd = eventMain.EventEnd

						c.EventBookingRoomRepo.Update(eventBookingRoom)

						if req.IsPublic == 0 {
							// Remove the "[" and "]" characters from the string
							trimmed := strings.Trim(*req.Person, "[]")

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
						return
					}
				}
			}
		} else if req.Type == "hybrid" {
			// Type Hybrid

			existRoom, _ := c.EventBookingRoomRepo.FindExistRoom(*req.IDRoom, main_event.Id, main_event.EventStart, main_event.EventEnd)

			eventBookingRoom, errEventBookingRoom := c.EventBookingRoomRepo.FindRoomBooking(req.IDRoom, main_event.Id)

			if errEventBookingRoom == nil {
				if !existRoom {
					main_event.EventRoom = req.IDRoom

					if req.URL != nil {
						main_event.EventUrl = req.URL
					}

					eventMain, err_eventMain := c.MainEventRepo.Update(main_event)

					if err_eventMain == nil {
						var list_id_person []int
						eventBookingRoom.IdEvent = eventMain.Id
						eventBookingRoom.CodeRoom = eventMain.EventRoom
						eventBookingRoom.DateStart = eventMain.EventStart
						eventBookingRoom.DateEnd = eventMain.EventEnd

						c.EventBookingRoomRepo.Update(eventBookingRoom)

						if req.IsPublic == 0 {
							// Remove the "[" and "]" characters from the string
							trimmed := strings.Trim(*req.Person, "[]")

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
						return
					}
				}
			}
		}
	} else {
		ev_main.EventTitle = req.Title
		ev_main.EventDesc = req.Desc
		ev_main.EventType = req.Type
		ev_main.EventStart, _ = time.Parse(time.DateTime, req.Start)
		ev_main.EventEnd, _ = time.Parse(time.DateTime, req.End)
		ev_main.Status = req.Status
		ev_main.CreatedBy = strconv.Itoa(req.Nik)
		ev_main.CompCode = comp_code
		ev_main.EventLocation = req.Location

		if req.ApprovalPerson != nil {
			ev_main.ApprovalPerson = req.ApprovalPerson
		}

		if req.Type == "online" {
			// Type Online
			ev_main.EventUrl = req.URL

			event_main, _ := c.MainEventRepo.Create(ev_main)

			if req.IsPublic == 0 {
				if req.Person != nil && *req.Person != "" {
					// Remove the "[" and "]" characters from the string
					trimmed := strings.Trim(*req.Person, "[]")

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

			if req.FileMateri != nil {
				for _, files_materi := range req.FileMateri {
					ev_materi_file.IdEvent = event_main.Id
					ev_materi_file.FileName = files_materi.FileName
					ev_materi_file.FileUrl = files_materi.FileURL
					c.EventMateriFileRepo.Create(ev_materi_file)
				}
			}
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
			})
		} else if req.Type == "offline" {
			// Type Offline
			existRoom, _ := c.EventBookingRoomRepo.FindExistRoom(*req.IDRoom, ev_main.Id, ev_main.EventStart, ev_main.EventEnd)

			if !existRoom {
				ev_main.EventRoom = req.IDRoom
				// main_event.ApprovalPerson = "7222322"

				event_main, _ := c.MainEventRepo.Create(ev_main)

				ev_br.CodeRoom = event_main.EventRoom
				ev_br.IdEvent = event_main.Id
				ev_br.DateStart = event_main.EventStart
				ev_br.DateEnd = event_main.EventEnd

				c.EventBookingRoomRepo.Create(ev_br)

				if req.IsPublic == 0 {
					if req.Person != nil && *req.Person != "" {
						// Remove the "[" and "]" characters from the string
						trimmed := strings.Trim(*req.Person, "[]")

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

				if req.FileMateri != nil {
					for _, files_materi := range req.FileMateri {
						ev_materi_file.IdEvent = event_main.Id
						ev_materi_file.FileName = files_materi.FileName
						ev_materi_file.FileUrl = files_materi.FileURL
						c.EventMateriFileRepo.Create(ev_materi_file)
					}
				}
				ctx.JSON(http.StatusOK, gin.H{
					"status":  http.StatusOK,
					"success": "Success",
				})
			}
		} else if req.Type == "hybrid" {
			// Type Hybrid
			existRoom, _ := c.EventBookingRoomRepo.FindExistRoom(*req.IDRoom, ev_main.Id, ev_main.EventStart, ev_main.EventEnd)

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

				if req.IsPublic == 0 {
					if req.Person != nil && *req.Person != "" {
						// Remove the "[" and "]" characters from the string
						trimmed := strings.Trim(*req.Person, "[]")

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

				if req.FileMateri != nil {
					for _, files_materi := range req.FileMateri {
						ev_materi_file.IdEvent = event_main.Id
						ev_materi_file.FileName = files_materi.FileName
						ev_materi_file.FileUrl = files_materi.FileURL
						c.EventMateriFileRepo.Create(ev_materi_file)
					}
				}
				ctx.JSON(http.StatusOK, gin.H{
					"status":  http.StatusOK,
					"success": "Success",
				})
			}
		}
	}
	ctx.AbortWithStatus(http.StatusInternalServerError)
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
		main_event.ApprovalPerson = &req.Nik
		main_event.Status = req.Status
		main_event.EventKeterangan = &req.Keterangan

		c.MainEventRepo.Update(main_event)

		if req.Status == "Declined" {
			if main_event.EventRoom != nil {
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

func (c *EventController) GetDataApproval(ctx *gin.Context) {
	nik := ctx.Param("nik")
	var ev []Authentication.EventData

	main_event, _ := c.MainEventRepo.FindEventMainNikApprovalPerson(nik)

	data_list_event := make([]Authentication.EventData, len(main_event))
	for i, DataEvent := range main_event {
		data_list_event[i] = Authentication.EventData{
			EventID:        DataEvent.Id,
			EventTitle:     DataEvent.EventTitle,
			EventDesc:      DataEvent.EventDesc,
			EventStart:     DataEvent.EventStart.Format("2006-01-02 15:04:05"),
			EventEnd:       DataEvent.EventEnd.Format("2006-01-02 15:04:05"),
			EventType:      DataEvent.EventType,
			EventImgName:   DataEvent.EventImgName,
			EventImgURL:    DataEvent.EventImgUrl,
			EventDate:      DataEvent.EventStart.Format("2006-01-02"),
			EventTimeStart: DataEvent.EventStart.Format("15:04:05"),
			EventTimeEnd:   DataEvent.EventEnd.Format("15:04:05"),
			EventStatus:    DataEvent.Status,
			CompCode:       DataEvent.CompCode,
			AssetCompCode:  "https://pismart-dev.pupuk-indonesia.com/public/assets/media/logos/logo-pi-full.png",
		}

		data_list_event[i].EventURL = DataEvent.EventUrl

	}
	ev = data_list_event

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    ev,
	})
}

func (c *EventController) KonfirmasiKehadiran(ctx *gin.Context) {
	var req Authentication.ValidasiKonfirmasi
	var ev Authentication.EventShowEvent

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

	ev_person, err_ev_person := c.EventPersonRepo.FindEventPersonIDNIK(req.EventID, req.Nik)

	if err_ev_person != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		ev_person.StatusKehadiran = req.Status
		c.EventPersonRepo.Update(ev_person)

		main_event, err_main_event := c.MainEventRepo.FindEventMainID(ev_person.IdEvent)
		if err_main_event == nil {
			ev.EventID = main_event.Id
			ev.EventTitle = main_event.EventTitle
			ev.EventDesc = main_event.EventDesc
			ev.EventStart = main_event.EventStart.Format("2006-01-02 15:04:05")
			ev.EventDateStart = main_event.EventStart.Format("2006-01-02")
			ev.EventEnd = main_event.EventEnd.Format("2006-01-02 15:04:05")
			ev.EventDateEnd = main_event.EventEnd.Format("2006-01-02")
			ev.EventDateFormat = main_event.EventStart.Format("02 January 2006")
			ev.EventType = main_event.EventType
			ev.EventTimeStart = main_event.EventStart.Format("15:04:05")
			ev.EventTimeEnd = main_event.EventEnd.Format("15:04:05")
			ev.EventStatus = main_event.Status

			ev.EventURL = main_event.EventUrl

			ev.EventLocation = main_event.EventLocation

			ev.EventImgName = main_event.EventImgName
			ev.EventImgURL = main_event.EventImgUrl
			ev.EventDate = main_event.EventStart.Format("2006-01-02")

			eventCount, _ := c.EventPersonRepo.GetEventCounts(main_event.Id)
			persons, _ := c.EventPersonRepo.FindDetailEventPerson(main_event.Id)

			data_list := []Authentication.EventPersonDetail{}
			for _, dataPerson := range persons {
				person := Authentication.EventPersonDetail{
					Nik:             dataPerson.Nik,
					Nama:            dataPerson.Nama,
					DeptTitle:       dataPerson.DeptTitle,
					Email:           dataPerson.Email,
					StatusKehadiran: dataPerson.StatusKehadiran,
					PhotoURL:        dataPerson.PhotoURL,
				}
				if dataPerson.Nama == "" {
					person.Nama = "tidak ada"
				}
				if dataPerson.DeptTitle == "" {
					person.DeptTitle = "-"
				}
				if dataPerson.Email == "" {
					person.Email = "tidak ada"
				}
				if dataPerson.PhotoURL == "" {
					person.PhotoURL = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
				}

				data_list = append(data_list, person)
			}
			ev.Person = data_list

			notulen, _ := c.EventNotulenRepo.FindEventNotulenK(main_event.Id)
			if notulen == nil {
				ev.Notulen = nil
			} else {
				notulen_file, _ := c.EventNotulenRepo.GetDataNotulenFile(notulen.IdNotulen)
				var dataNotulenFile *Authentication.GetDataNotulenFiles
				dataNotulenFile.IdNotulen = notulen.IdNotulen
				dataNotulenFile.IdEvent = notulen.IdEvent
				dataNotulenFile.Deskripsi = notulen.Deskripsi
				dataNotulenFile.CreatedAt = notulen.CreatedAt
				dataNotulenFile.UpdatedAt = notulen.UpdatedAt
				dataNotulenFile.Files = notulen_file

				ev.Notulen = dataNotulenFile
			}

			ev_materi, _ := c.EventMateriFileRepo.FindEventMateriFile(main_event.Id)
			is_absent, _ := c.EventPresenceRepo.FindPresenceIDNIK(main_event.Id, main_event.CreatedBy)

			ev.CompCode = main_event.CompCode
			ev.StatusKehadiran = ev_person.StatusKehadiran
			ev.Count = eventCount
			ev.IsAbsent = is_absent
			ev.AssestCompCode = "https://pismart-dev.pupuk-indonesia.com/public/assets/media/logos/logo-pi-full.png"
			ev.Materi = ev_materi
			ev.TimeCreatedAt = main_event.CreatedAt.Format("15:04:05")
			ev.EventRoom = main_event.EventRoom

			if main_event.EventType == "offline" || main_event.EventType == "hybrid" {
				book_room, err := c.EventBookingRoomRepo.FindBookRoomShow(main_event.EventRoom, main_event.Id)
				if err != nil {
					ev.BookRoom = nil
				} else {
					br := &Authentication.DataBookRoomShow{
						IDBooking:    book_room.IDBooking,
						CodeRoom:     book_room.CodeRoom,
						DateStart:    book_room.DateStart.Format("2006-01-02 15:04:05"),
						DateEnd:      book_room.DateEnd.Format("2006-01-02 15:04:05"),
						TimeStart:    book_room.DateStart.Format("15:04:05"),
						TimeEnd:      book_room.DateEnd.Format("15:04:05"),
						RoomID:       book_room.RoomID,
						RoomName:     book_room.RoomName,
						RoomCategory: book_room.RoomCategory,
						RoomCompCode: book_room.RoomCompCode,
						RoomCompName: book_room.RoomCompName,
					}
					ev.BookRoom = br
				}
			} else {
				ev.BookRoom = nil
			}

			list_nik := []*string{&main_event.CreatedBy, main_event.ApprovalPerson}
			var list_nik_strings []string
			for _, nik := range list_nik {
				list_nik_strings = append(list_nik_strings, *nik)
			}
			result, _ := c.PihcMasterKaryDbRepo.FindUserByNIKArray(list_nik_strings)
			ev.EventCreatedByNik = &main_event.CreatedBy
			ev.EventApprovalPerson = main_event.ApprovalPerson

			for _, data := range result {
				if main_event.CreatedBy == data.EmpNo {
					ev.EventCreatedBy = data.Nama
					ev.EventCreatedDeptTitle = data.DeptTitle
				}
				if main_event.ApprovalPerson == &data.EmpNo {
					ev.EventApprovalPersonName = data.Nama
				}
			}

			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
				"data":    ev,
			})
		} else {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"succes": "Data Not Found",
			})
		}
	}
}
func (c *EventController) GetDataInFeed(ctx *gin.Context) {
	nik := ctx.Param("nik")
	var data []Authentication.DataInFeed
	typeHari := []string{"Hari Ini", "Besok", "Lusa"}

	for _, hari := range typeHari {
		dataInFeed, jumlahEvent, _ := c.MainEventRepo.FindDataInFeed(nik, hari)
		model_data := &Authentication.Model{
			ID:         dataInFeed.Id,
			EventTitle: dataInFeed.EventTitle,
			EventDesc:  dataInFeed.EventDesc,
			EventStart: dataInFeed.EventStart.Format("2006-01-02 15:04:05"),
			EventEnd:   dataInFeed.EventEnd.Format("2006-01-02 15:04:05"),
			EventType:  dataInFeed.EventType,
			CompCode:   dataInFeed.CompCode,
			Status:     dataInFeed.Status,
			CreatedBy:  dataInFeed.CreatedBy,
			CreatedAt:  dataInFeed.CreatedAt,
			UpdatedAt:  dataInFeed.UpdatedAt,
		}
		if dataInFeed.EventUrl != nil {
			model_data.EventURL = dataInFeed.EventUrl
		}
		if dataInFeed.EventImgName != nil {
			model_data.EventImgName = dataInFeed.EventImgName
		}
		if dataInFeed.EventImgUrl != nil {
			model_data.EventImgURL = dataInFeed.EventImgUrl
		}
		if dataInFeed.ApprovalPerson != nil {
			model_data.ApprovalPerson = dataInFeed.ApprovalPerson
		}
		if dataInFeed.EventRoom != nil {
			model_data.EventRoom = dataInFeed.EventRoom
		}
		if dataInFeed.EventLocation != nil {
			model_data.EventLocation = dataInFeed.EventLocation
		}
		if dataInFeed.EventKeterangan != nil {
			model_data.EventKeterangan = dataInFeed.EventKeterangan
		}
		list_data := Authentication.DataInFeed{
			Type:        hari,
			Model:       model_data,
			JumlahEvent: int(jumlahEvent),
		}
		if model_data.ID == 0 {
			list_data.Model = nil
		}
		data = append(data, list_data)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data,
	})
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

	events_person := c.EventPersonRepo.FindEventPersonID(req.EventID)

	if events_person != nil {
		// Remove the "[" and "]" characters from the string
		trimmed := strings.Trim(req.Dispose, "[]")

		// Split the string into individual values
		values := strings.Split(trimmed, ",")

		for _, nikPerson := range values {
			ev_person.IdEvent = req.EventID
			ev_person.Nik = nikPerson

			idParent, _ := strconv.Atoi(req.Nik)
			ev_person.IdParent = &idParent
			ev_person.StatusKehadiran = "menunggu"
			c.EventPersonRepo.Create(ev_person)
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *EventController) GetDataDispose(ctx *gin.Context) {
	var req Authentication.ValidasiGetDataDispose

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

	ev_person, _ := c.EventPersonRepo.GetDataDisposePerson(req.Nik, req.IdEvent)

	if ev_person == nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		data := []Authentication.EventPersonDetail{}
		for _, dataPerson := range ev_person {
			person := Authentication.EventPersonDetail{
				Nik:             dataPerson.Nik,
				Nama:            dataPerson.Nama,
				DeptTitle:       dataPerson.DeptTitle,
				Email:           dataPerson.Email,
				StatusKehadiran: dataPerson.StatusKehadiran,
				PhotoURL:        dataPerson.PhotoURL,
			}
			if dataPerson.DeptTitle == "" {
				person.DeptTitle = "-"
			}
			if dataPerson.PhotoURL == "" {
				person.PhotoURL = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
			}

			data = append(data, person)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    data,
		})
	}
}

func (c *EventController) GetDataEvent(ctx *gin.Context) {
	nik := ctx.Param("nik")
	var ev []Authentication.EventDataByNik

	main_event, _ := c.MainEventRepo.FindEventMainNikCreatedBy(nik)

	data_list := make([]Authentication.EventDataByNik, len(main_event))

	for i, DataEvent := range main_event {
		data_list[i].EventData.EventID = DataEvent.Id
		data_list[i].EventData.EventTitle = DataEvent.EventTitle
		data_list[i].EventData.EventDesc = DataEvent.EventDesc
		data_list[i].EventData.EventStart = DataEvent.EventStart.Format("2006-01-02 15:04:05")
		data_list[i].EventData.EventEnd = DataEvent.EventEnd.Format("2006-01-02 15:04:05")
		data_list[i].EventData.EventType = DataEvent.EventType
		data_list[i].EventData.EventURL = DataEvent.EventUrl
		data_list[i].EventData.EventImgName = DataEvent.EventImgName
		data_list[i].EventData.EventImgURL = DataEvent.EventImgUrl
		data_list[i].EventData.EventDate = DataEvent.EventStart.Format("2006-01-02")
		data_list[i].EventData.EventTimeStart = DataEvent.EventStart.Format("15:04:05")
		data_list[i].EventData.EventTimeEnd = DataEvent.EventEnd.Format("15:04:05")
		data_list[i].EventData.EventStatus = DataEvent.Status
		data_list[i].EventLocation = DataEvent.EventLocation
		data_list[i].CompCode = DataEvent.CompCode
		data_list[i].AssetCompCode = "https://pismart-dev.pupuk-indonesia.com/public/assets/media/logos/logo-pi-full.png"

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
			EventType:       DataEvent.EventType,
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
	var ev Authentication.EventShowEvent

	mainEvent, err := c.MainEventRepo.FindEventMainIDNIK(idEvent, nik)
	if err == nil {
		ev.EventID = mainEvent.Id
		ev.EventTitle = mainEvent.EventTitle
		ev.EventDesc = mainEvent.EventDesc
		ev.EventStart = mainEvent.EventStart.Format("2006-01-02 15:04:05")
		ev.EventDateStart = mainEvent.EventStart.Format("2006-01-02")
		ev.EventEnd = mainEvent.EventEnd.Format("2006-01-02 15:04:05")
		ev.EventDateEnd = mainEvent.EventEnd.Format("2006-01-02")
		ev.EventDateFormat = mainEvent.EventStart.Format("02 January 2006")
		ev.EventType = mainEvent.EventType
		ev.EventTimeStart = mainEvent.EventStart.Format("15:04:05")
		ev.EventTimeEnd = mainEvent.EventEnd.Format("15:04:05")
		ev.EventStatus = mainEvent.Status
		ev.EventURL = mainEvent.EventUrl
		ev.EventLocation = mainEvent.EventLocation
		ev.EventImgName = mainEvent.EventImgName
		ev.EventImgURL = mainEvent.EventImgUrl
		ev.EventDate = mainEvent.EventStart.Format("2006-01-02")

		eventCount, _ := c.EventPersonRepo.GetEventCounts(mainEvent.Id)
		ev_person, _ := c.EventPersonRepo.FindDetailEventPerson(mainEvent.Id)

		data_list := []Authentication.EventPersonDetail{}
		for _, dataPerson := range ev_person {
			person := Authentication.EventPersonDetail{
				Nik:             dataPerson.Nik,
				Nama:            dataPerson.Nama,
				DeptTitle:       dataPerson.DeptTitle,
				Email:           dataPerson.Email,
				StatusKehadiran: dataPerson.StatusKehadiran,
				PhotoURL:        dataPerson.PhotoURL,
			}
			if dataPerson.Nama == "" {
				person.Nama = "tidak ada"
			}
			if dataPerson.DeptTitle == "" {
				person.DeptTitle = "-"
			}
			if dataPerson.Email == "" {
				person.Email = "tidak ada"
			}
			if dataPerson.PhotoURL == "" {
				person.PhotoURL = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
			}

			data_list = append(data_list, person)
		}
		ev.Person = data_list

		notulen, _ := c.EventNotulenRepo.FindEventNotulenK(mainEvent.Id)
		if notulen == nil {
			ev.Notulen = nil
		} else {
			notulen_file, _ := c.EventNotulenRepo.GetDataNotulenFile(notulen.IdNotulen)
			var dataNotulenFile Authentication.GetDataNotulenFiles
			dataNotulenFile.IdNotulen = notulen.IdNotulen
			dataNotulenFile.IdEvent = notulen.IdEvent
			dataNotulenFile.Deskripsi = notulen.Deskripsi
			dataNotulenFile.CreatedAt = notulen.CreatedAt
			dataNotulenFile.UpdatedAt = notulen.UpdatedAt
			dataNotulenFile.Files = notulen_file

			ev.Notulen = &dataNotulenFile
		}

		ev_materi, _ := c.EventMateriFileRepo.FindEventMateriFile(mainEvent.Id)
		is_absent, _ := c.EventPresenceRepo.FindPresenceIDNIK(mainEvent.Id, mainEvent.CreatedBy)
		person, _ := c.EventPersonRepo.FindEventPersonIDNIK(mainEvent.Id, mainEvent.CreatedBy)

		ev.CompCode = mainEvent.CompCode
		// -------------------- BELUM -----------------------------
		ev.StatusKehadiran = person.StatusKehadiran
		// --------------------------------------------------------
		ev.Count = eventCount
		ev.IsAbsent = is_absent
		ev.AssestCompCode = "https://pismart-dev.pupuk-indonesia.com/public/assets/media/logos/logo-pi-full.png"
		ev.Materi = ev_materi
		ev.TimeCreatedAt = mainEvent.CreatedAt.Format("15:04:05")
		ev.EventRoom = mainEvent.EventRoom

		if mainEvent.EventType == "offline" || mainEvent.EventType == "hybrid" {
			book_room, err := c.EventBookingRoomRepo.FindBookRoomShow(mainEvent.EventRoom, mainEvent.Id)
			if err != nil {
				fmt.Println("A")
				ev.BookRoom = nil
			} else {
				fmt.Println("B")
				br := &Authentication.DataBookRoomShow{
					IDBooking:    book_room.IDBooking,
					CodeRoom:     book_room.CodeRoom,
					DateStart:    book_room.DateStart.Format("2006-01-02 15:04:05"),
					DateEnd:      book_room.DateEnd.Format("2006-01-02 15:04:05"),
					TimeStart:    book_room.DateStart.Format("15:04:05"),
					TimeEnd:      book_room.DateEnd.Format("15:04:05"),
					RoomID:       book_room.RoomID,
					RoomName:     book_room.RoomName,
					RoomCategory: book_room.RoomCategory,
					RoomCompCode: book_room.RoomCompCode,
					RoomCompName: book_room.RoomCompName,
				}
				ev.BookRoom = br
			}
		} else {
			fmt.Println("C")
			ev.BookRoom = nil
		}

		list_nik := []*string{&mainEvent.CreatedBy, mainEvent.ApprovalPerson}
		var list_nik_strings []string
		for _, nik := range list_nik {
			list_nik_strings = append(list_nik_strings, *nik)
		}
		result, _ := c.PihcMasterKaryDbRepo.FindUserByNIKArray(list_nik_strings)
		ev.EventCreatedByNik = &mainEvent.CreatedBy
		ev.EventApprovalPerson = mainEvent.ApprovalPerson

		for _, data := range result {
			if mainEvent.CreatedBy == data.EmpNo {
				ev.EventCreatedBy = data.Nama
				ev.EventCreatedDeptTitle = data.DeptTitle
			}
			if mainEvent.ApprovalPerson == &data.EmpNo {
				ev.EventApprovalPersonName = data.Nama
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    ev,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"succes": "Data Not Found",
		})
	}
}

func (c *EventController) GetBookingRoom(ctx *gin.Context) {
	var req Authentication.ValidasiGetBookingRoom
	var data_br []Authentication.DataBookRoomDate

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

	result, _ := c.EventBookingRoomRepo.GetBookingRoomDate(req.IdRoom, req.Date)

	data_list := []Authentication.DataBookRoomDate{}
	for _, data := range result {
		br := Authentication.DataBookRoomDate{
			IDBooking:         data.IDBooking,
			NamaEvent:         data.NamaEvent,
			NamaPembuat:       data.NamaPembuat,
			KompatemenPembuat: data.KompatemenPembuat,
			DateStart:         data.DateStart.Format("2006-01-02 15:04:05"),
			DateEnd:           data.DateEnd.Format("2006-01-02 15:04:05"),
			TimeStart:         data.DateStart.Format("15:04:05"),
			TimeEnd:           data.DateEnd.Format("15:04:05"),
			Date:              data.DateStart.Format("2006-01-02"),
		}
		data_list = append(data_list, br)
	}
	data_br = data_list

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data_br,
	})
}

func (c *EventController) DeleteEvent(ctx *gin.Context) {
	var req Authentication.ValidasiDeleteEventByID
	var data Authentication.ListEventDelete
	var ev_hist events.EventHistory

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

	mainEvent, errMain := c.MainEventRepo.DeleteMainEvent(req.Id)
	if errMain == nil {
		if mainEvent.EventType == "offline" || mainEvent.EventType == "hybrid" {
			c.EventBookingRoomRepo.DeleteRoomBooking(mainEvent.Id)
		}

		ev_person, ev_person_parent, _ := c.EventPersonRepo.DeleteEventPerson(mainEvent.Id)
		c.EventPersonRepo.DeleteEventPerson(mainEvent.Id)
		c.EventNotulenRepo.DeleteEventNotulen(mainEvent.Id)
		c.EventMateriFileRepo.DeleteEventMateriFile(mainEvent.Id)
		c.EventPresenceRepo.DeleteEventPresence(mainEvent.Id)

		data.ID = mainEvent.Id
		data.EventTitle = mainEvent.EventTitle
		data.EventDesc = mainEvent.EventDesc
		data.EventStart = mainEvent.EventStart.Format("2006-01-02 15:04:05")
		data.EventEnd = mainEvent.EventEnd.Format("2006-01-02 15:04:05")
		data.EventType = mainEvent.EventType
		if mainEvent.EventUrl != nil {
			data.EventURL = mainEvent.EventUrl
		}
		if mainEvent.EventImgName != nil {
			data.EventImgName = mainEvent.EventImgName
		}
		if mainEvent.EventImgUrl != nil {
			data.EventImgURL = mainEvent.EventImgUrl
		}
		data.CompCode = mainEvent.CompCode
		data.Status = mainEvent.Status
		data.CreatedBy = mainEvent.CreatedBy
		data.CreatedAt = mainEvent.CreatedAt
		data.UpdatedAt = mainEvent.UpdatedAt
		if mainEvent.ApprovalPerson != nil {
			data.ApprovalPerson = mainEvent.ApprovalPerson
		}
		if mainEvent.EventRoom != nil {
			data.EventRoom = mainEvent.EventRoom
		}
		if mainEvent.EventLocation != nil {
			data.EventLocation = mainEvent.EventLocation
		}
		if mainEvent.EventKeterangan != nil {
			data.EventKeterangan = mainEvent.EventKeterangan
		}

		data_list := []Authentication.Persons{}
		list_dispose := []Authentication.Disposes{}
		for _, dataPersonIsNotParent := range ev_person {
			br := Authentication.Persons{
				EventPerson: dataPersonIsNotParent,
			}
			for _, dataPersonIsParent := range ev_person_parent {
				if dataPersonIsNotParent.Nik == strconv.Itoa(*dataPersonIsParent.IdParent) {
					fmt.Println("ISPARENT")
					dispose := Authentication.Disposes{
						EventPerson: dataPersonIsParent,
						Profile:     nil,
					}
					list_dispose = append(list_dispose, dispose)
				}
			}
			br.Dispose = list_dispose
			data_list = append(data_list, br)
		}
		data.Person = data_list

		// parsedTimeStart, _ := time.Parse(time.DateTime, data.EventStart)
		// parsedTimeEnd, _ := time.Parse(time.DateTime, data.EventEnd)

		ev_hist.EventTitle = data.EventTitle
		ev_hist.EventDesc = data.EventDesc
		ev_hist.EventStart, _ = time.Parse(time.DateTime, data.EventStart)
		ev_hist.EventEnd, _ = time.Parse(time.DateTime, data.EventEnd)
		ev_hist.EventType = data.EventType
		if data.EventURL != nil {
			ev_hist.EventURL = *data.EventURL
		}
		if data.EventImgName != nil {
			ev_hist.EventImgName = *data.EventImgName
		}
		if data.EventImgURL != nil {
			ev_hist.EventImgURL = *data.EventImgURL
		}
		ev_hist.CompCode = data.CompCode
		ev_hist.Status = data.Status
		ev_hist.CreatedBy = data.CreatedBy
		ev_hist.CreatedAt = data.CreatedAt
		if data.ApprovalPerson != nil {
			ev_hist.ApprovalPerson = *data.ApprovalPerson
		}
		if data.EventRoom != nil {
			ev_hist.EventRoom = *data.EventRoom
		}
		if data.EventLocation != nil {
			ev_hist.EventLocation = *data.EventLocation
		}
		ev_hist.EventKeterangan = req.Keterangan

		fmt.Println(data.UpdatedAt)
		fmt.Println(ev_hist.UpdatedAt)

		c.MainEventRepo.History(ev_hist)

		ctx.JSON(http.StatusOK, gin.H{
			"Status:":  http.StatusOK,
			"Success:": "Delete Success",
			"Data:":    data,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *EventController) DeleteEventBooking(ctx *gin.Context) {
	idBooking := ctx.Param("id_booking")
	int_idBooking, _ := strconv.Atoi(idBooking)

	data_booking_room, err := c.EventBookingRoomRepo.FindRoomBookingByIdBooking(int_idBooking)

	if err == nil {
		c.EventBookingRoomRepo.DeleteRoomBooking(data_booking_room.IdEvent)
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   data_booking_room,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

// func (c *EventController) DeleteFileNotulen(ctx *gin.Context) {
// 	idEvent := ctx.Param("id")
// 	int_idEvent, _ := strconv.Atoi(idEvent)

// 	mainEvent, errMainEvent := c.MainEventRepo.FindEventMainID(int_idEvent)

// 	if errMainEvent == nil {
// 		errNotulen := c.EventNotulenRepo.DeleteEventNotulen(mainEvent.Id)
// 		if errNotulen == nil {
// 			fmt.Println("XXXXX")

// 			ctx.JSON(http.StatusOK, gin.H{
// 				"status": http.StatusOK,
// 				"info":   "Success",
// 			})
// 		} else {
// 			ctx.AbortWithStatus(http.StatusInternalServerError)
// 		}
// 	} else {
// 		ctx.AbortWithStatus(http.StatusInternalServerError)
// 	}
// }

// func (c *EventController) DeleteFileMateri(ctx *gin.Context) {
// 	idEvent := ctx.Param("id")
// 	int_idEvent, _ := strconv.Atoi(idEvent)

// 	mainEvent, errMainEvent := c.MainEventRepo.FindEventMainID(int_idEvent)

// 	if errMainEvent == nil {
// 		errDeleteEventMateri := c.EventMateriFileRepo.DeleteEventMateriFile(mainEvent.Id)
// 		if errDeleteEventMateri == nil {
// 			ctx.JSON(http.StatusOK, gin.H{
// 				"status": http.StatusOK,
// 				"info":   "Success",
// 			})
// 		} else {
// 			ctx.AbortWithStatus(http.StatusInternalServerError)
// 		}
// 	} else {
// 		ctx.AbortWithStatus(http.StatusInternalServerError)
// 	}
// }

func (c *EventController) DeleteFileNotulen(ctx *gin.Context) {
	IdNotulenFile := ctx.Param("id")
	int_IdNotulenFile, _ := strconv.Atoi(IdNotulenFile)

	data, errNotulenFile := c.EventNotulenRepo.DeleteEventNotulenFiles(int_IdNotulenFile)
	if errNotulenFile == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Delete Success",
			"data":    data,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *EventController) DeleteFileMateri(ctx *gin.Context) {
	idMateriFile := ctx.Param("id")
	int_idMateriFile, _ := strconv.Atoi(idMateriFile)

	data, errMateriFile := c.EventMateriFileRepo.DeleteMateriFiles(int_idMateriFile)
	if errMateriFile == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Delete Success",
			"data":    data,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *EventController) StoreNotulen(ctx *gin.Context) {
	var req Authentication.ValidasiStoreNotulen
	var notulen events.EventNotulen
	var notulen_files events.EventNotulenFile

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

	// idEvent, _ := strconv.Atoi(req.IdEvent)

	notulen.IdEvent = req.IdEvent
	notulen.Deskripsi = req.Deskripsi

	notulens, _ := c.EventNotulenRepo.CreateNotulen(notulen)
	form, _ := ctx.MultipartForm()
	files := form.File["file"]

	for _, file := range files {
		originalFileName := file.Filename

		fmt.Println(originalFileName)

		// Open the multipart.FileHeader to get a multipart.File
		fileToUpload, err := file.Open()
		if err != nil {
			// Handle the error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not open file",
			})
			return
		}

		file_url, file_name, err := c.EventNotulenRepo.UploadFile(originalFileName, fileToUpload)
		if err == nil {
			notulen_files.IdNotulen = notulens.IdNotulen
			notulen_files.FileName = file_name
			notulen_files.FileUrl = file_url

			c.EventNotulenRepo.CreateNotulenFiles(notulen_files)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"info":   "Success",
	})

	// for _, file := range files {
	// 	fileExt := filepath.Ext(file.Filename)
	// 	originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
	// 	now := time.Now()
	// 	filename := fmt.Sprintf("%v", now.Unix()) + "_" + strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + fileExt
	// 	filePath := "https://storage.googleapis.com/lumen-oauth-storage/Event/Notulen/2023/" + filename // Change this to a valid local directory path

	// 	// if err := ctx.SaveUploadedFile(file, filePath); err != nil {
	// 	// 	ctx.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
	// 	// 	return
	// 	// }

	// 	filePaths = append(filePaths, filePath)
	// }

	// ctx.JSON(http.StatusOK, gin.H{"filepath": filePaths})
}

func (c *EventController) StoreFileGCS(ctx *gin.Context) {
	var req Authentication.ValidasiStoreNotulen

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

	form, _ := ctx.MultipartForm()
	files := form.File["file"]
	filePaths := []string{}

	for _, file := range files {
		originalFileName := file.Filename

		fmt.Println(originalFileName)

		// Open the multipart.FileHeader to get a multipart.File
		fileToUpload, err := file.Open()
		if err != nil {
			// Handle the error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not open file",
			})
			return
		}

		imageURL, _, err := c.EventNotulenRepo.UploadFile(originalFileName, fileToUpload)
		if err != nil {
			// Handle the error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not upload file",
			})
			return
		}

		filePaths = append(filePaths, imageURL)
	}
	// Handle the uploaded file paths, e.g., return them in the response
	ctx.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"filePaths": filePaths,
	})
}

func (c *EventController) RenameFileGCS(ctx *gin.Context) {
	var req Authentication.ValidasiRenameFileNotulen

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

	url, err := c.EventNotulenRepo.RenameFileGCS(req.OldNameFile, req.NewNameFile)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"filePaths": url,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"notice": "Gagal Mengganti Nama File",
		})
	}
}

func (c *EventController) DeleteFileGCS(ctx *gin.Context) {
	var req Authentication.ValidasiDeleteFileNotulen

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

	err := c.EventNotulenRepo.DeleteFileGCS(req.NameFile)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"notice": "success",
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"notice": "Gagal Menghapus File",
		})
	}
}

func (c *EventController) GetDataNotulen(ctx *gin.Context) {
	id := ctx.Param("id")
	idEvent, _ := strconv.Atoi(id)

	notulen, _ := c.EventNotulenRepo.FindEventNotulenK(idEvent)
	notulen_file, _ := c.EventNotulenRepo.GetDataNotulenFile(notulen.IdNotulen)

	var data Authentication.GetDataNotulenFiles
	data.IdNotulen = notulen.IdNotulen
	data.IdEvent = notulen.IdEvent
	data.Deskripsi = notulen.Deskripsi
	data.CreatedAt = notulen.CreatedAt
	data.UpdatedAt = notulen.UpdatedAt
	data.Files = notulen_file

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data,
	})
}

func (c *EventController) GetCategoryRoom(ctx *gin.Context) {
	var req Authentication.ValidasiKonfirmasiNik

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

	pihc_mstr_krywn, err_pihc_mstr_krywn := c.PihcMasterKaryDbRepo.FindUserByNIK(req.Nik)

	var data []string
	if err_pihc_mstr_krywn == nil {
		category_room, _ := c.EventMsterRoomRepo.FindCategoryRoom(pihc_mstr_krywn.Company)

		for _, ctr_data := range category_room {
			data = append(data, ctr_data.CategoryRoom.CategoryRoom)
		}

		if data == nil {
			data = []string{}
		}
	} else {
		category_room, _ := c.EventMsterRoomRepo.FindDefaultCategoryRoom("A000")

		for _, ctr_data := range category_room {
			data = append(data, ctr_data.CategoryRoom.CategoryRoom)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data,
	})
}

func (c *EventController) GetRoomEvent(ctx *gin.Context) {
	var req Authentication.ValidasiKonfirmasiGetEventRoom

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

	pihc_mstr_krywn, err_pihc_mstr_krywn := c.PihcMasterKaryDbRepo.FindUserByNIK(req.Nik)

	var data []events.EventMsterRoom
	if err_pihc_mstr_krywn == nil {
		category_room, _ := c.EventMsterRoomRepo.FindRoomEvent(pihc_mstr_krywn.Company, req.CategoryRoom)
		data = category_room
	} else {
		category_room, _ := c.EventMsterRoomRepo.FindDefaultRoomEvent("A000")
		data = category_room
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data,
	})
}

func (c *EventController) StoreBookingRoomEvent(ctx *gin.Context) {
	var req Authentication.ValidasiStoreBookingRoom
	var ev_br events.EventBookingRoom

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

	parsedTimeStart, _ := time.Parse(time.DateTime, req.DateStart)
	parsedTimeEnd, _ := time.Parse(time.DateTime, req.DateEnd)

	existRoom, _ := c.EventBookingRoomRepo.FindExistRoom(req.IDRoom, req.IDEvent, parsedTimeStart, parsedTimeEnd)

	if !existRoom {
		br_mine, _ := c.EventBookingRoomRepo.FindRoomBooking(&req.IDRoom, req.IDEvent)
		if br_mine.IdEvent == req.IDEvent {
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"status": 503,
				"info":   "Maaf Ruangan sudah terisi",
			})
		} else {
			if req.IDEvent != 0 {
				ev_br.IdEvent = req.IDEvent
			} else {
				ev_br.IdEvent = 0
			}
			ev_br.CodeRoom = &req.IDRoom
			ev_br.DateStart = parsedTimeStart
			ev_br.DateEnd = parsedTimeEnd

			c.EventBookingRoomRepo.Create(ev_br)

			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
			})
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"status": 503,
			"info":   "Maaf Ruangan sudah terisi",
		})
	}
}

func (c *EventController) StoreEventPresence(ctx *gin.Context) {
	var req Authentication.ValidasiStoreEventPresence
	var ev_presence events.EventPresence

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
	ev_main, err_ev_main := c.MainEventRepo.FindEventMainIDType(req.IdPresence, req.Type)

	if err_ev_main != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Event tidak Ditemukan",
		})
	} else {
		ev_person, err_ev_person := c.EventPersonRepo.FindEventPersonIDNIK(ev_main.Id, req.Nik)
		if err_ev_person != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"info":   "Anda tidak terdaftar dalam list Event",
			})
		} else {
			ev_presence.IdEvent = ev_main.Id
			ev_presence.EmpNo = ev_person.Nik
			parsedPresenceDateTime, _ := time.Parse(time.DateTime, req.PresenceDateTime)
			ev_presence.PresenceDateTime = parsedPresenceDateTime
			ev_presence_check, _ := c.EventPresenceRepo.FindPresenceIDNIK(ev_main.Id, req.Nik)

			if !ev_presence_check {
				if ev_main.EventType == "online" {
					c.EventPresenceRepo.Create(ev_presence)
				} else if ev_main.EventType == "offline" || ev_main.EventType == "hybrid" {
					if ev_main.EventRoom != nil {
						ev_presence.IdRoom = *ev_main.EventRoom
					}
					c.EventPresenceRepo.Create(ev_presence)
				}
				ctx.JSON(http.StatusOK, gin.H{
					"status": http.StatusOK,
					"info":   "Success",
				})
			} else {
				ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"status": http.StatusServiceUnavailable,
					"info":   "Anda sudah absent pada event ini",
				})
			}
		}
	}
}

func (c *EventController) PrintDaftarHadir(ctx *gin.Context) {
	id := ctx.Param("id")
	idEvent, _ := strconv.Atoi(id)
	var ev Authentication.EventPrintDaftarHadir

	mainEvent, errMainEvent := c.MainEventRepo.FindEventMainID(idEvent)
	if errMainEvent == nil {
		pihc_krywn, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(mainEvent.CreatedBy)
		ev.Title = mainEvent.EventTitle
		ev.Deskripsi = mainEvent.EventDesc
		ev.Type = mainEvent.EventType
		ev.Tanggal = mainEvent.EventStart.Format("02 January 2006")
		ev.JamMulai = mainEvent.EventStart.Format("15:04:05")
		ev.JamSelesai = mainEvent.EventEnd.Format("15:04:05")
		if mainEvent.EventUrl != nil {
			ev.EventURL = mainEvent.EventUrl
		}
		if mainEvent.EventLocation != nil {
			ev.EventLocation = mainEvent.EventLocation
		}
		ev.CompIcon = "public/assets/media/logos/logo-pi-full.png"
		ev.CreatedBy = mainEvent.CreatedBy
		ev.CreatedByName = pihc_krywn.Nama
		ev.CreatedByDept = pihc_krywn.DeptTitle

		if mainEvent.EventType == "offline" || mainEvent.EventType == "hybrid" {
			book_room, err := c.EventBookingRoomRepo.FindBookRoomShow(mainEvent.EventRoom, mainEvent.Id)
			if err != nil {
				ev.BookRoom = nil
			} else {
				br := &Authentication.BookRoomDaftarHadir{
					RoomName:     book_room.RoomName,
					RoomCategory: book_room.RoomCategory,
					RoomCompName: book_room.RoomCompName,
				}
				ev.BookRoom = br
			}
		} else {
			ev.BookRoom = nil
		}

		presence, _ := c.EventPresenceRepo.FindDetailEventPresence(mainEvent.Id)

		data_list := []Authentication.Presences{}
		for i, detailPresence := range presence {
			ev_presence := Authentication.Presences{
				Nomer:    i + 1,
				Nik:      detailPresence.Nik,
				Nama:     detailPresence.Nama,
				Email:    detailPresence.Email,
				Jabatan:  detailPresence.Jabatan,
				NoTelp:   detailPresence.NoTelp,
				Dept:     detailPresence.Dept,
				Time:     detailPresence.Presence.Format("02 January 2006 15:04:05"),
				Date:     detailPresence.Presence.Format("02 January 2006"),
				Datetime: detailPresence.Presence.Format("15:04:05"),
				Instansi: detailPresence.Instansi,
			}

			if detailPresence.Nik == "" {
				ev_presence.Nik = "Non App User"
			}
			if detailPresence.Dept == "" {
				ev_presence.Dept = "-"
			}

			data_list = append(data_list, ev_presence)
		}
		ev.Presence = data_list

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    ev,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
