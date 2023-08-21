package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl/mstrKgt"
	users "github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type MstrKgtController struct {
	KegiatanMasterRepo   *mstrKgt.KegiatanMasterRepo
	PihcMasterKaryRtRepo *mstrKgt.PihcMasterKaryRtRepo
}

func NewMstrKgtController(db *gorm.DB) *MstrKgtController {
	return &MstrKgtController{
		KegiatanMasterRepo:   mstrKgt.NewKegiatanMasterRepo(db),
		PihcMasterKaryRtRepo: mstrKgt.NewPihcMasterKaryRtRepo(db)}
}

func (c *MstrKgtController) ListMasterKegiatan(ctx *gin.Context) {
	var inputan Authentication.ValidationLMK

	// nik := ctx.PostForm("nik")
	// tahun := ctx.PostForm("tahun")

	if err := ctx.ShouldBindJSON(&inputan); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "NIK / Tahun Tidak Boleh Kosong"})
		return
	}

	PIHC_MSTR_KRY_RT, err := c.PihcMasterKaryRtRepo.FindUserByNIK(inputan.NIK)

	comp_code := PIHC_MSTR_KRY_RT.Company

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Karyawan Tidak Ada",
			"Data":   nil,
		})
	}

	KegiatanMaster, err := c.KegiatanMasterRepo.FindUserByCompCodeYear(comp_code, inputan.Tahun)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   KegiatanMaster})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Tidak Ada",
			"Data":   nil,
		})
	}
}

func (c *MstrKgtController) StoreMasterKegiatan(ctx *gin.Context) {
	var km mstrKgt.KegiatanMaster
	var req Authentication.ValidationSMK

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

	if req.IdKegiatan != 0 {
		kgt_mstr, err_kgtmstr := c.KegiatanMasterRepo.FindData(req.IdKegiatan)

		if req.NamaKegiatan != "" {
			kgt_mstr.NamaKegiatan = req.NamaKegiatan
		}

		if req.DeskripsiKegiatan != "" {
			kgt_mstr.DeskripsiKegiatan = req.DeskripsiKegiatan
		}

		if err_kgtmstr == nil {
			kgt_mstr, err_update := c.KegiatanMasterRepo.Update(kgt_mstr)
			if err_update == nil {
				ctx.JSON(http.StatusOK, gin.H{
					"status":  http.StatusOK,
					"success": "Success",
					"data":    "Data berhasil diUpdate",
					"result":  &kgt_mstr,
				})
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"success": "Gagal mengupdate data",
				})
			}
		} else {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"info":   "Data Tidak Ada",
				"Data":   nil,
			})
		}
	} else {
		t := time.Now()

		if req.NamaKegiatan != "" {
			km.NamaKegiatan = req.NamaKegiatan
		}

		if req.DeskripsiKegiatan != "" {
			km.DeskripsiKegiatan = req.DeskripsiKegiatan
		}
		km.Periode = strconv.Itoa(t.Year())
		km.CompCode = comp_code
		km.CreatedBy = req.NIK
		km.Slug = users.String(12)
		km, err = c.KegiatanMasterRepo.Create(km)
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
				"data":    "Data berhasil ditambahkan",
				"result":  &km,
			})
		} else {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"success": "Gagal menambahkan data",
			})
		}
	}

}

func (c *MstrKgtController) DeleteMasterKegiatan(ctx *gin.Context) {
	slug := ctx.Param("slug")
	data, err := c.KegiatanMasterRepo.DelMasterKegiatanID(slug)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   data})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Tidak Ada",
			"Data":   nil,
		})
	}
}
