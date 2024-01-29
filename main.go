package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/yusufwira/lern-golang-gin/connection"
	"github.com/yusufwira/lern-golang-gin/controller"
	"github.com/yusufwira/lern-golang-gin/controller/cuti_karyawan_controller"
	"github.com/yusufwira/lern-golang-gin/controller/jobtender"
	"github.com/yusufwira/lern-golang-gin/controller/mobile_api/event_controller"
	"github.com/yusufwira/lern-golang-gin/controller/mobile_api/profile_controller"
	"github.com/yusufwira/lern-golang-gin/controller/tjsl_controller"
)

func main() {
	db := connection.Database()

	mstrKgtController := tjsl_controller.NewMstrKgtController(db.Db, db.StorageClient)
	kgtKrywnController := tjsl_controller.NewKgtKrywnController(db.Db, db.StorageClient)
	koorkgtController := tjsl_controller.NewKoorKgtController(db.Db, db.StorageClient)
	eventController := event_controller.NewEventController(db.Db, db.StorageClient)
	userProfileController := profile_controller.NewUsersProfileController(db.Db, db.StorageClient)
	UserController := controller.NewUserController(db.Db, db.StorageClient)
	cutiKrywnController := cuti_karyawan_controller.NewCutiKrywnController(db.Db)
	jobTenderController := jobtender.GetJobVacancyController(db.Db)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://127.0.0.1:8000"},
		AllowMethods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "x-csrf-token", "Authorization"},
	}))

	r.Use(gin.Recovery())
	r.POST("/login", UserController.Login)
	r.POST("/register", UserController.Register)

	auth := r.Group("/api")
	{
		auth.GET("/token", connection.Middleware)
		auth.GET("/testToken", connection.Validation)

		auth.GET("/getUserOuath", connection.Validation, func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.Index())
		})

		auth.POST("/getKaryawanNameAll", connection.Validation, UserController.GetDataKaryawanNameAll)
		auth.GET("/getKaryawanAll", connection.Validation, UserController.GetDataKaryawanAll)
		auth.GET("/getAtasanPegawai/:nik", UserController.GetAtasanPegawai)
		auth.POST("/getKaryawanNameIndiv", connection.Validation, UserController.GetDataKaryawanNameIndiv)

		auth.POST("/postUser", connection.Validation, func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.Store(c))
		})
		auth.DELETE("/delUserID/:id", connection.Validation, func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.DelData(c))
		})
		auth.GET("/getUserID/:id", connection.Validation, func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.GetData(c))
		})
		auth.PUT("/upUserID/:id", connection.Validation, func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.UpData(c))
		})
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	jobtender := auth.Group("job_tender")
	{
		jobtender.GET("/detailJobVacancy/:id", jobTenderController.GetDetailJob)
	}

	tjsl := auth.Group(os.Getenv("TJSL_API_URL"))
	{
		// Master Kegiatan
		tjsl.POST("/listKegiatan", connection.Validation, mstrKgtController.ListMasterKegiatan)
		tjsl.POST("/storeMasterKegiatan", connection.Validation, mstrKgtController.StoreMasterKegiatan)
		tjsl.GET("/getMasterKegiatan/:slug", connection.Validation, mstrKgtController.GetMasterKegiatan)
		tjsl.DELETE("/deleteMasterKegiatan/:slug", connection.Validation, mstrKgtController.DeleteMasterKegiatan)

		// Pengajuan Kegiatan
		tjsl.POST("/storePengajuan", connection.Validation, kgtKrywnController.StorePengajuanKegiatan)
		tjsl.GET("/showPengajuan/:slug", connection.Validation, kgtKrywnController.ShowDetailPengajuanKegiatan)
		tjsl.GET("/myTjsl", connection.Validation, kgtKrywnController.ShowPengajuanKegiatan)
		tjsl.DELETE("/deletePengajuan/:slug", connection.Validation, kgtKrywnController.DeletePengajuanKegiatan)

		tjsl.POST("/approve", connection.Validation, kgtKrywnController.StoreApprovePengajuanKegiatan)
		tjsl.POST("/listApprovalTjsl", connection.Validation, kgtKrywnController.ListApprvlKgtKrywn)
		tjsl.GET("/getChartSummary", connection.Validation, kgtKrywnController.GetChartSummary)
		tjsl.POST("/getLeaderBoard", connection.Validation, kgtKrywnController.GetLeaderBoardKgtKrywn)

		// Koordinator
		tjsl.POST("/storeKoordinator", connection.Validation, koorkgtController.StoreKoordinator)
		tjsl.GET("/showKoordinator/:slug", connection.Validation, koorkgtController.ShowDetailKoordinator)
		tjsl.DELETE("/deleteKoordinator/:slug", connection.Validation, koorkgtController.DeleteKoordinator)
		tjsl.GET("/listKoordinator", connection.Validation, koorkgtController.ListKoordinator)
	}

	event := auth.Group(os.Getenv("EVENT_API_URL"))
	{
		event.POST("/store_new", connection.Validation, eventController.StoreEvent)
		event.POST("/updateStatusEvent", connection.Validation, eventController.UpdateStatusEvent)
		event.GET("/getDataApproval/:nik", connection.Validation, eventController.GetDataApproval)
		event.POST("/konfirmasiKehadiran", connection.Validation, eventController.KonfirmasiKehadiran)
		event.POST("/getDataInFeed/:nik", connection.Validation, eventController.GetDataInFeed)

		// DISPOSE
		event.POST("/storeDispose", connection.Validation, eventController.StoreDispose)
		event.POST("/getDataDispose", connection.Validation, eventController.GetDataDispose)

		event.GET("/getDataEvent/:nik", connection.Validation, eventController.GetDataEvent)
		event.POST("/getDataByNik", connection.Validation, eventController.GetDataByNik)
		event.POST("/deleteEvent", connection.Validation, eventController.DeleteEvent)
		event.GET("/showEvent/:id/:nik", connection.Validation, eventController.ShowEvent)
		event.DELETE("/deleteEventBooking/:id_booking", connection.Validation, eventController.DeleteEventBooking)

		// GCS
		event.POST("/storeFileGCS", connection.Validation, eventController.StoreFileGCS)
		event.POST("/renameFileGCS", connection.Validation, eventController.RenameFileGCS)
		event.POST("/deleteFileGCS", connection.Validation, eventController.DeleteFileGCS)

		// NOTULEN
		event.POST("/storeNotulen", connection.Validation, eventController.StoreNotulen)
		event.GET("/getDataNotulen/:id", connection.Validation, eventController.GetDataNotulen)
		event.DELETE("/deleteFileNotulen/:id", connection.Validation, eventController.DeleteFileNotulen)

		// ROOM MASTER
		event.GET("/getCategoryRoom", connection.Validation, eventController.GetCategoryRoom)
		event.POST("/getDataRoom", connection.Validation, eventController.GetRoomEvent)
		event.POST("/getBookingRoom", connection.Validation, eventController.GetBookingRoom)

		event.POST("/storeBookingRoom", connection.Validation, eventController.StoreBookingRoomEvent)

		// PRESENCE
		event.POST("/storeEventPresence", connection.Validation, eventController.StoreEventPresence)
		event.GET("/printDaftarHadir/:id", connection.Validation, eventController.PrintDaftarHadir)
	}

	personalInformation := auth.Group(os.Getenv("PERSONALINFORMATION_API_URL"))
	{
		personalInformation.POST("/storeData", connection.Validation, userProfileController.StoreData)
		personalInformation.POST("/getData", connection.Validation, userProfileController.GetData)
		personalInformation.GET("/getCategory", connection.Validation, userProfileController.GetCategory)
	}

	profile := auth.Group(os.Getenv("MOBILE_API_URL"))
	{
		profile.POST("/storeProfile", connection.Validation, userProfileController.StoreProfile)
		profile.POST("/storeAboutUs", connection.Validation, userProfileController.StoreAboutUs)
		profile.GET("/showAboutUs/:nik", connection.Validation, userProfileController.GetShowAboutUs)
		profile.GET("/getSosialMediaInformation/:nik", connection.Validation, userProfileController.GetSocialMediaInformation)

		profile.POST("/storeInformationContact", connection.Validation, userProfileController.StoreInformationContact)
		profile.POST("/storeSkill", connection.Validation, userProfileController.StoreSkill)
		profile.POST("/updateSkill", connection.Validation, userProfileController.UpdateSkill)
		profile.POST("/deleteSkill", connection.Validation, userProfileController.DeleteSkill)
		profile.POST("/updatePhotoProfile", connection.Validation, userProfileController.UpdatePhotoProfile)
		profile.GET("/getSkill/:nik", connection.Validation, userProfileController.GetSkill)
		profile.GET("/getPengalamanKerja/:nik", connection.Validation, userProfileController.GetPengalamanKerja)
		profile.GET("/getContactInformation/:nik", connection.Validation, userProfileController.GetContactInformation)
		profile.GET("/showProfile/:nik", connection.Validation, userProfileController.ShowProfile)

		profile.GET("/dataPegawai", connection.Validation, userProfileController.DataPegawai)
		profile.GET("/getAtasanPegawai", connection.Validation, userProfileController.DataAtasanPegawai)
	}

	cuti := auth.Group(os.Getenv("CUTI_URL"))
	{
		// PENGAJUAN CUTI
		cuti.POST("/storeCuti", connection.Validation, cutiKrywnController.StoreCutiKaryawan)
		cuti.GET("/getTipeAbsenSaldoPengajuan", connection.Validation, cutiKrywnController.GetTipeAbsenSaldoPengajuan)
		cuti.GET("/myCuti", connection.Validation, cutiKrywnController.GetMyPengajuanCuti)
		cuti.GET("/showPengajuanCuti/:id_pengajuan_absen", connection.Validation, cutiKrywnController.ShowDetailPengajuanCuti)
		cuti.DELETE("/deletePengajuanCuti/:id_pengajuan_absen", connection.Validation, cutiKrywnController.DeletePengajuanCuti)

		// Approval
		cuti.POST("/listApprovalCuti", connection.Validation, cutiKrywnController.ListApprvlCuti)
		cuti.GET("/showApprovalPengajuanCuti/:id_pengajuan_absen", connection.Validation, cutiKrywnController.ShowDetailApprovalPengajuanCuti)
		cuti.POST("/approve", connection.Validation, cutiKrywnController.StoreApprovePengajuanAbsen)

		// SALDO CUTI
		cuti.POST("/storeAdminSaldo", connection.Validation, cutiKrywnController.StoreAdminSaldoCutiKaryawan)
		cuti.POST("/listAdminSaldo", connection.Validation, cutiKrywnController.ListAdminSaldoCutiKaryawan)
		cuti.GET("/getAdminSaldoCuti/:id_saldo_cuti", connection.Validation, cutiKrywnController.GetAdminSaldoCuti)
		cuti.GET("/getAdminTipeAbsen", connection.Validation, cutiKrywnController.GetAdminTipeAbsen)
		cuti.DELETE("/deleteAdminSaldoCuti/:id_saldo_cuti", connection.Validation, cutiKrywnController.DeleteAdminSaldoCuti)

		// Company
		cuti.GET("/getCompany", connection.Validation, cutiKrywnController.GetCompany)
		// Direktorat
		cuti.POST("/getDirektorat", connection.Validation, cutiKrywnController.GetDirektorat)
		// Kompartemen
		cuti.POST("/getKompartemen", connection.Validation, cutiKrywnController.GetKompartemen)
		// Departemen
		cuti.POST("/getDepartemen", connection.Validation, cutiKrywnController.GetDepartemen)
	}

	r.Run(os.Getenv("PORT_RUN")) // local
}
