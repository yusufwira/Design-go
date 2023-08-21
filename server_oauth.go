package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"

	"github.com/yusufwira/lern-golang-gin/connection"
	"github.com/yusufwira/lern-golang-gin/controller"
)

func main() {
	db := connection.Database()
	mstrKgtController := controller.NewMstrKgtController(db)
	kgtKrywnController := controller.NewKgtKrywnController(db)
	UserController := controller.NewUserController(db)

	r := gin.Default()

	connection.Middleware()

	auth := r.Group("/oauth2")
	{
		auth.GET("/token", ginserver.HandleTokenRequest)
	}

	api := r.Group("/api")
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
	}
	tjsl := r.Group("/api/tjsl")
	{
		tjsl.POST("/listKegiatan", mstrKgtController.ListMasterKegiatan)
		tjsl.POST("/storeMasterKegiatan", mstrKgtController.StoreMasterKegiatan)
		tjsl.DELETE("/deleteMasterKegiatan/:slug", mstrKgtController.DeleteMasterKegiatan)

		tjsl.POST("/storePengajuan", kgtKrywnController.StorePengajuanKegiatan)
		// tjsl.POST("listApprovalTjsl", kgtKrywnController.ListApprvlKgtKrywn)
		tjsl.GET("/showPengajuan/:slug", kgtKrywnController.ShowDetailPengajuanKegiatan)
		tjsl.GET("/myTjsl", kgtKrywnController.ShowPengajuanKegiatan)
		tjsl.DELETE("/deletePengajuan/:slug", kgtKrywnController.DeletePengajuanKegiatan)
	}
	r.Run(":9096")
}
