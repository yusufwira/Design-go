package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl/kgtKrywn"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl/mstrKgt"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl/photosKgt"
	users "github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type KgtKrywnController struct {
	KegiatanKaryawanRepo *kgtKrywn.KegiatanKaryawanRepo
	KegiatanMasterRepo   *mstrKgt.KegiatanMasterRepo
	KegiatanPhotosRepo   *photosKgt.KegiatanPhotosRepo
	PihcMasterKaryRtRepo *mstrKgt.PihcMasterKaryRtRepo
}

func NewKgtKrywnController(db *gorm.DB) *KgtKrywnController {
	return &KgtKrywnController{KegiatanKaryawanRepo: kgtKrywn.NewKegiatanKaryawanRepo(db),
		KegiatanMasterRepo:   mstrKgt.NewKegiatanMasterRepo(db),
		KegiatanPhotosRepo:   photosKgt.NewKegiatanPhotosRepo(db),
		PihcMasterKaryRtRepo: mstrKgt.NewPihcMasterKaryRtRepo(db)}
}

func (c *KgtKrywnController) StoreKgtKrywn(ctx *gin.Context) {
	var kk kgtKrywn.KegiatanKaryawan
	var kp photosKgt.KegiatanPhotos
	var req Authentication.AuthenticationKK

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "NIK Tidak Boleh Kosong"})
		return
	}

	PIHC_MSTR_KRY_RT, err := c.PihcMasterKaryRtRepo.FindUserByNIK(req.NIK)

	comp_code := PIHC_MSTR_KRY_RT.Company

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Karyawan Tidak Ada",
			"Data":   nil})
		return
	}

	if req.KegiatanParentId == 0 {
		kk.KegiatanParentId = 0
	} else {
		kk.KegiatanParentId = req.KegiatanParentId
	}

	t := time.Now()
	kk.NIK = req.NIK
	kk.Id = req.Id
	kk.NamaKegiatan = req.NamaKegiatan
	kk.TanggalKegiatan = req.TanggalKegiatan
	kk.LokasiKegiatan = req.LokasiKegiatan
	kk.DeskripsiKegiatan = req.DeskripsiKegiatan
	kk.Status = "WaitApv"
	kk.Manager = ""
	kk.CompCode = comp_code
	kk.Periode = strconv.Itoa(t.Year())

	if kk.Id != 0 {
		kgt_krywn := c.KegiatanKaryawanRepo.FindData(kk.Id)
		kk.Slug = kgt_krywn.Slug
		kk.CreatedAt = kgt_krywn.CreatedAt
		kk, err = c.KegiatanKaryawanRepo.Update(kk)
		data_kp := c.KegiatanPhotosRepo.FindData(kk.Id)
		fmt.Println(len(data_kp), len(req.Photos))
		if len(data_kp) == len(req.Photos) {
			for i := range data_kp {
				data_kp[i].OriginalName = req.Photos[i].OriginalName
				data_kp[i].Url = req.Photos[i].URL
				url, _ := c.KegiatanPhotosRepo.GetFileExtensionFromUrl(data_kp[i].Url)
				data_kp[i].Extendtion = url
				c.KegiatanPhotosRepo.Update(data_kp[i])
				fmt.Println("_______________X_________________")
			}
			fmt.Println("_______________Y_________________")
		}
		fmt.Println("_______________Z_________________")
		// for data1, data2 := range req.Photos,data_kp {
		// 	for _, data2 := range data_kp {
		// 		data2.OriginalName = data.OriginalName
		// 		data2.Url = data.URL
		// 		url, _ := c.KegiatanPhotosRepo.GetFileExtensionFromUrl(data2.Url)
		// 		data2.Extendtion = url
		// 		c.KegiatanPhotosRepo.Update(data2)
		// 	}
		// }
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
				"data":    "Data berhasil diUpdate",
			})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"success": "Gagal mengupdate data",
			})
		}
	} else {
		fmt.Println(req.CreatedAt)
		kk.Slug = users.String(12)
		kk, err = c.KegiatanKaryawanRepo.Create(kk)
		for _, data := range req.Photos {
			kp.KegiatanId = kk.Id
			kp.OriginalName = data.OriginalName
			kp.Url = data.URL
			url, _ := c.KegiatanPhotosRepo.GetFileExtensionFromUrl(kp.Url)
			kp.Extendtion = url
			// s := c.KegiatanPhotosRepo.LastString(strings.Split(data.OriginalName, "."))
			// kp.Extendtion = s
			c.KegiatanPhotosRepo.Create(kp)
		}
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
				"data":    "Data berhasil ditambahkan",
			})
		} else {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"success": "Gagal menambahkan data",
			})
		}
	}

}

func (c *KgtKrywnController) ListKgtKrywn(ctx *gin.Context) {

}

func (c *KgtKrywnController) ShowKgtKrywn(ctx *gin.Context) {

}

func (c *KgtKrywnController) DeleteKgtKrywn(ctx *gin.Context) {

}
