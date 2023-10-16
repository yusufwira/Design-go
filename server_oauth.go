package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"
	"github.com/joho/godotenv"

	"github.com/yusufwira/lern-golang-gin/connection"
	"github.com/yusufwira/lern-golang-gin/controller"
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

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://127.0.0.1:8000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "x-csrf-token"},
	}))

	r.Use(gin.Recovery())

	connection.Middleware()

	auth := r.Group("/oauth2")
	{
		auth.GET("/token", ginserver.HandleTokenRequest)
	}

	api := auth.Group("/api")
	{
		fmt.Println("masuk")
		api.Use(ginserver.HandleTokenVerify())
		fmt.Println("masuk2")
		api.GET("/test", func(c *gin.Context) {
			ti, exists := c.Get(ginserver.DefaultConfig.TokenKey)
			if exists {
				c.JSON(http.StatusOK, ti)
				return
			}
			c.String(http.StatusOK, "not found")
		})

		r.GET("/getUserOuath", func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.Index())
		})

		r.GET("/getUserID/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.GetData(c))
		})

		r.POST("/getKaryawanNameAll", UserController.GetDataKaryawanNameAll)
		r.POST("/getKaryawanNameIndiv", UserController.GetDataKaryawanNameIndiv)

		r.POST("/postUser", func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.Store(c))
		})

		r.DELETE("/delUserID/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.DelData(c))
		})
		r.PUT("/upUserID/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, UserController.UpData(c))
		})

		r.POST("/login", UserController.Login)
		r.POST("/register", UserController.Register)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	tjsl := r.Group(os.Getenv("TJSL_API_URL"))
	{
		// Master Kegiatan
		tjsl.POST("/listKegiatan", mstrKgtController.ListMasterKegiatan)
		tjsl.POST("/storeMasterKegiatan", mstrKgtController.StoreMasterKegiatan)
		tjsl.GET("/getMasterKegiatan/:slug", mstrKgtController.GetMasterKegiatan)
		tjsl.DELETE("/deleteMasterKegiatan/:slug", mstrKgtController.DeleteMasterKegiatan)

		// Pengajuan Kegiatan
		tjsl.POST("/storePengajuan", kgtKrywnController.StorePengajuanKegiatan)
		tjsl.GET("/showPengajuan/:slug", kgtKrywnController.ShowDetailPengajuanKegiatan)
		tjsl.GET("/myTjsl", kgtKrywnController.ShowPengajuanKegiatan)
		tjsl.DELETE("/deletePengajuan/:slug", kgtKrywnController.DeletePengajuanKegiatan)

		tjsl.POST("/approve", kgtKrywnController.StoreApprovePengajuanKegiatan)
		tjsl.POST("/listApprovalTjsl", kgtKrywnController.ListApprvlKgtKrywn)
		tjsl.GET("/getChartSummary", kgtKrywnController.GetChartSummary)
		tjsl.POST("/getLeaderBoard", kgtKrywnController.GetLeaderBoardKgtKrywn)

		// Koordinator
		tjsl.POST("/storeKoordinator", koorkgtController.StoreKoordinator)
		tjsl.GET("/showKoordinator/:slug", koorkgtController.ShowDetailKoordinator)
		tjsl.DELETE("/deleteKoordinator/:slug", koorkgtController.DeleteKoordinator)
		tjsl.GET("/listKoordinator", koorkgtController.ListKoordinator)
	}

	event := r.Group(os.Getenv("EVENT_API_URL"))
	{
		event.POST("/store_new", eventController.StoreEvent)
		event.POST("/updateStatusEvent", eventController.UpdateStatusEvent)
		event.GET("/getDataApproval/:nik", eventController.GetDataApproval)
		event.POST("/konfirmasiKehadiran", eventController.KonfirmasiKehadiran)
		event.POST("/getDataInFeed/:nik", eventController.GetDataInFeed)

		// DISPOSE
		event.POST("/storeDispose", eventController.StoreDispose)
		event.POST("/getDataDispose", eventController.GetDataDispose)

		event.GET("/getDataEvent/:nik", eventController.GetDataEvent)
		event.POST("/getDataByNik", eventController.GetDataByNik)
		event.POST("/deleteEvent", eventController.DeleteEvent)
		event.GET("/showEvent/:id/:nik", eventController.ShowEvent)
		event.DELETE("/deleteEventBooking/:id_booking", eventController.DeleteEventBooking)

		// GCS
		event.POST("/storeFileGCS", eventController.StoreFileGCS)
		event.POST("/renameFileGCS", eventController.RenameFileGCS)
		event.POST("/deleteFileGCS", eventController.DeleteFileGCS)

		// NOTULEN
		event.POST("/storeNotulen", eventController.StoreNotulen)
		event.GET("/getDataNotulen/:id", eventController.GetDataNotulen)
		event.DELETE("/deleteFileNotulen/:id", eventController.DeleteFileNotulen)

		// ROOM MASTER
		event.GET("/getCategoryRoom", eventController.GetCategoryRoom)
		event.POST("/getDataRoom", eventController.GetRoomEvent)
		event.POST("/getBookingRoom", eventController.GetBookingRoom)

		event.POST("/storeBookingRoom", eventController.StoreBookingRoomEvent)

		// PRESENCE
		event.POST("/storeEventPresence", eventController.StoreEventPresence)
		event.GET("/printDaftarHadir/:id", eventController.PrintDaftarHadir)

	}

	personalInformation := r.Group(os.Getenv("PERSONALINFORMATION_API_URL"))
	{
		personalInformation.POST("/storeData", userProfileController.StoreData)
		personalInformation.POST("/getData", userProfileController.GetData)
		personalInformation.GET("/getCategory", userProfileController.GetCategory)
	}

	profile := r.Group(os.Getenv("MOBILE_API_URL"))
	{
		profile.POST("/storeProfile", userProfileController.StoreProfile)
		profile.POST("/storeAboutUs", userProfileController.StoreAboutUs)
		profile.GET("/showAboutUs/:nik", userProfileController.GetShowAboutUs)
		profile.GET("/getSosialMediaInformation/:nik", userProfileController.GetSocialMediaInformation)

		profile.POST("/storeInformationContact", userProfileController.StoreInformationContact)
		profile.POST("/storeSkill", userProfileController.StoreSkill)
		profile.POST("/updateSkill", userProfileController.UpdateSkill)
		profile.POST("/deleteSkill", userProfileController.DeleteSkill)
		profile.POST("/updatePhotoProfile", userProfileController.UpdatePhotoProfile)
		profile.GET("/getSkill/:nik", userProfileController.GetSkill)
		profile.GET("/getPengalamanKerja/:nik", userProfileController.GetPengalamanKerja)
		profile.GET("/getContactInformation/:nik", userProfileController.GetContactInformation)
		profile.GET("/showProfile/:nik", userProfileController.ShowProfile)
	}

	r.Run(":9096")
}
