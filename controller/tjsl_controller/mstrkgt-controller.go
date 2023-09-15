package tjsl_controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl"
	users "github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type MstrKgtController struct {
	KegiatanMasterRepo   *tjsl.KegiatanMasterRepo
	PihcMasterKaryRtRepo *pihc.PihcMasterKaryRtRepo
}

func NewMstrKgtController(db *gorm.DB) *MstrKgtController {
	return &MstrKgtController{
		KegiatanMasterRepo:   tjsl.NewKegiatanMasterRepo(db),
		PihcMasterKaryRtRepo: pihc.NewPihcMasterKaryRtRepo(db)}
}

func (c *MstrKgtController) ListMasterKegiatan(ctx *gin.Context) {
	var req Authentication.ValidationLMK

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

	PIHC_MSTR_KRY_RT, err := c.PihcMasterKaryRtRepo.FindUserByNIK(req.NIK)

	comp_code := PIHC_MSTR_KRY_RT.Company

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Karyawan Tidak Ada",
			"Data":   nil,
		})
	}

	KegiatanMaster, err := c.KegiatanMasterRepo.FindUserByCompCodeYear(comp_code, req.Tahun)

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

func (c *MstrKgtController) GetMasterKegiatan(ctx *gin.Context) {
	slug := ctx.Param("slug")

	KegiatanMaster, err := c.KegiatanMasterRepo.FindDataBySlug(slug)

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
	var km tjsl.KegiatanMaster
	var req Authentication.ValidationSMK

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
		kgt_mstr, err_kgtmstr := c.KegiatanMasterRepo.FindDataById(req.IdKegiatan)

		if req.NamaKegiatan != "" {
			kgt_mstr.NamaKegiatan = req.NamaKegiatan
		}

		if req.DeskripsiKegiatan != nil {
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

		if req.DeskripsiKegiatan != nil {
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
