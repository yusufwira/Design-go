package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/mstrKgt"
	users "github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type MstrKgtController struct {
	TJSLRepo *mstrKgt.TJSLRepo
	DboRepo  *mstrKgt.DboRepo
}

func NewMstrKgtController(db *gorm.DB) *MstrKgtController {
	return &MstrKgtController{TJSLRepo: mstrKgt.NewTJSLRepo(db)}
}

func (c *MstrKgtController) ListMasterKegiatan(ctx *gin.Context) {
	var inputan Authentication.AuthenticationLMK

	if err := ctx.ShouldBindJSON(&inputan); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "NIK / Tahun Tidak Boleh Kosong"})
		return
	}

	PIHC_MSTR_KRY_RT, err := c.DboRepo.FindUserByNIK(inputan.NIK)

	comp_code := PIHC_MSTR_KRY_RT.Company

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Karyawan Tidak Ada",
			"Data":   nil,
		})
	}

	KegiatanMaster, err := c.TJSLRepo.FindUserByCompCodeYear(comp_code, inputan.Tahun)

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
	var inputan Authentication.AuthenticationSMK

	if err := ctx.ShouldBindJSON(&inputan); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "NIK Tidak Boleh Kosong"})
		return
	}

	PIHC_MSTR_KRY_RT, err := c.DboRepo.FindUserByNIK(inputan.NIK)

	comp_code := PIHC_MSTR_KRY_RT.Company

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Karyawan Tidak Ada",
			"Data":   nil})
		return
	}

	km.NamaKegiatan = inputan.NamaKegiatan
	km.DeskripsiKegiatan = inputan.DeskripsiKegiatan
	km.Periode = inputan.Periode
	km.CompCode = comp_code
	km.CreatedBy = inputan.NIK
	km.Slug = users.String(12)

	km, err = c.TJSLRepo.Create(km)
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

func (c *MstrKgtController) UpdateMasterKegiatan(ctx *gin.Context) {
	slug := ctx.Param("slug")
	var req Authentication.AuthenticationSMK

	km, err := c.TJSLRepo.FindNIKbySlug(slug)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Karyawan Tidak Ada",
			"Data":   nil,
		})
		return
	}

	req.NIK = km.CreatedBy

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if req.NamaKegiatan != "" {
		km.NamaKegiatan = req.NamaKegiatan
	}
	if req.DeskripsiKegiatan != "" {
		km.DeskripsiKegiatan = req.DeskripsiKegiatan
	}

	if req.Periode != "" {
		km.Periode = req.Periode
	}

	km, err = c.TJSLRepo.Update(km)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    "Data berhasil diUpdate",
			"result":  &km,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"success": "Gagal mengupdate data",
		})
	}
}

func (c *MstrKgtController) DeleteMasterKegiatan(ctx *gin.Context) {
	slug := ctx.Param("slug")
	data, err := c.TJSLRepo.DelMasterKegiatanID(slug)

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
