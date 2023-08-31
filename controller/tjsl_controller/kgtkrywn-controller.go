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

type KgtKrywnController struct {
	KegiatanKaryawanRepo   *tjsl.KegiatanKaryawanRepo
	KegiatanMasterRepo     *tjsl.KegiatanMasterRepo
	KegiatanPhotosRepo     *tjsl.KegiatanPhotosRepo
	PihcMasterKaryRtRepo   *pihc.PihcMasterKaryRtRepo
	PihcMasterPositionRepo *pihc.PihcMasterPositionRepo
}

func NewKgtKrywnController(db *gorm.DB) *KgtKrywnController {
	return &KgtKrywnController{KegiatanKaryawanRepo: tjsl.NewKegiatanKaryawanRepo(db),
		KegiatanMasterRepo:     tjsl.NewKegiatanMasterRepo(db),
		KegiatanPhotosRepo:     tjsl.NewKegiatanPhotosRepo(db),
		PihcMasterKaryRtRepo:   pihc.NewPihcMasterKaryRtRepo(db),
		PihcMasterPositionRepo: pihc.NewPihcMasterPositionRepo(db)}
}

func (c *KgtKrywnController) StorePengajuanKegiatan(ctx *gin.Context) {
	var kk tjsl.KegiatanKaryawan
	var kp tjsl.KegiatanPhotos
	var req Authentication.ValidationSKKgt

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

	if req.Id != 0 {
		kgt_krywn, err_kgtkrywn := c.KegiatanKaryawanRepo.FindDataID(req.Id)
		kgt_krywn.NamaKegiatan = req.NamaKegiatan

		// tgl_kegiatan, _ := time.Parse(time.DateOnly, req.TanggalKegiatan)
		// kgt_krywn.TanggalKegiatan = datatypes.Date(tgl_kegiatan)
		kgt_krywn.TanggalKegiatan = req.TanggalKegiatan
		kgt_krywn.LokasiKegiatan = req.LokasiKegiatan
		kgt_krywn.DeskripsiKegiatan = req.DeskripsiKegiatan

		if err_kgtkrywn == nil {
			kgt_krywn, err_updte_kgtkrywn := c.KegiatanKaryawanRepo.Update(kgt_krywn)
			if err_updte_kgtkrywn == nil {
				var list_id_foto []int
				var kegiatan_id int
				var is_koordinator int

				if kgt_krywn.KoordinatorId == 0 {
					kegiatan_id = kgt_krywn.Id
					is_koordinator = 0
				} else {
					kegiatan_id = kgt_krywn.KoordinatorId
					is_koordinator = 1
				}

				for _, data := range req.Photos {
					kp.KegiatanId = kegiatan_id
					kp.IsKoordinator = is_koordinator
					kp.OriginalName = data.OriginalName
					kp.Url = data.URL
					url, _ := c.KegiatanPhotosRepo.GetFileExtensionFromUrl(kp.Url)
					kp.Extendtion = url
					kgt_photos := c.KegiatanPhotosRepo.Create(kp)
					list_id_foto = append(list_id_foto, kgt_photos.Id)
				}

				c.KegiatanPhotosRepo.DelPhotosIDLama(kegiatan_id, list_id_foto)

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
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"success": "Data Tidak Ditemukan",
			})
		}
	} else {
		// t := time.Now()
		kk.NIK = req.NIK
		kk.KegiatanParentId = req.KegiatanParentId
		kk.KoordinatorId = req.KoordinatorId
		kk.NamaKegiatan = req.NamaKegiatan

		// Using datatypes.Date
		// tgl_kegiatan, _ := time.Parse(time.DateOnly, req.TanggalKegiatan)
		// kk.TanggalKegiatan = datatypes.Date(tgl_kegiatan)
		kk.TanggalKegiatan = req.TanggalKegiatan
		kk.LokasiKegiatan = req.LokasiKegiatan
		kk.DeskripsiKegiatan = req.DeskripsiKegiatan
		kk.Status = "WaitApv"
		kk.Manager = ""
		kk.CompCode = comp_code
		// kk.Periode = strconv.Itoa(t.Year())
		kk.Periode = req.Tahun
		kk.Slug = users.String(12)
		kk, err_kgtkrywn := c.KegiatanKaryawanRepo.Create(kk)

		var kegiatan_id int
		var is_koordinator int
		if kk.KoordinatorId == 0 {
			kegiatan_id = kk.Id
			is_koordinator = 0
		} else {
			kegiatan_id = kk.KoordinatorId
			is_koordinator = 1
		}

		if err_kgtkrywn == nil {
			for _, data := range req.Photos {
				kp.KegiatanId = kegiatan_id
				kp.IsKoordinator = is_koordinator
				kp.OriginalName = data.OriginalName
				kp.Url = data.URL
				url, _ := c.KegiatanPhotosRepo.GetFileExtensionFromUrl(kp.Url)
				kp.Extendtion = url
				// s := c.KegiatanPhotosRepo.LastString(strings.Split(data.OriginalName, "."))
				// kp.Extendtion = s
				c.KegiatanPhotosRepo.Create(kp)
			}

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

func (c *KgtKrywnController) ShowPengajuanKegiatan(ctx *gin.Context) {
	var req Authentication.ValidationMyTjsl

	if err := ctx.ShouldBindQuery(&req); err != nil {
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
	data, err := c.KegiatanKaryawanRepo.FindDataNIKPeriode(req.Nik, req.Tahun)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   data,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Tidak Ada",
			"Data":   nil,
		})
	}
}

func (c *KgtKrywnController) ShowDetailPengajuanKegiatan(ctx *gin.Context) {
	var data Authentication.KegiatanKaryawanPhotos
	slug := ctx.Param("slug")

	data_kk, err_kk := c.KegiatanKaryawanRepo.FindDataSlug(slug)

	var kegiatan_id int
	var is_koordinator int
	if data_kk.KoordinatorId == 0 {
		kegiatan_id = data_kk.Id
		is_koordinator = 0
	} else {
		kegiatan_id = data_kk.KoordinatorId
		is_koordinator = 1
	}
	data_kp := c.KegiatanPhotosRepo.FindDataPhotosID(kegiatan_id, is_koordinator)

	data_pihc, err_pihc := c.PihcMasterKaryRtRepo.FindUserByNIK(data_kk.NIK)

	data.IDKegiatan = data_kk.Id
	data.SlugKegiatan = data_kk.Slug
	data.Nik = data_kk.NIK
	data.Nama = data_pihc.Nama
	data.PhotoProfile = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
	data.DeptTitle = data_pihc.DeptTitle

	if data_kk.KegiatanParentId == 0 {
		data.Jenis = "Kegiatan sosial kemasyarakatan diluar perusahaan"
	} else {
		data.Jenis = "Kegiatan Tanggung Jawab Sosial dan Lingkungan (TJSL) perusahaan"
	}

	data.KoordinatorID = data_kk.KoordinatorId
	data.SlugKoordinator = 0
	data.SlugKegiatanParent = 0
	data.KegiatanParentID = data_kk.KegiatanParentId
	data.NamaKegiatan = data_kk.NamaKegiatan

	rfc339, _ := time.Parse(time.RFC3339, data_kk.TanggalKegiatan)
	tgl_kegiatan_nonformat := rfc339.Format(time.DateOnly)
	year, month, day := rfc339.Date()
	tanggal := strconv.Itoa(day)
	bulan := month.String()
	tahun := strconv.Itoa(year)
	tgl_kegiatan := tanggal + " " + bulan + " " + tahun

	data.TanggalKegiatan = tgl_kegiatan
	data.TanggalKegiatanNonFormat = tgl_kegiatan_nonformat
	data.LokasiKegiatan = data_kk.LokasiKegiatan
	data.Deskripsi = data_kk.DeskripsiKegiatan
	data.Status = data_kk.Status
	data.AlasanPenolakan = data_kk.DescDecline
	data.PhotoKegiatan = data_kp
	data.Tahun = data_kk.Periode

	if (err_kk == nil) || (err_pihc == nil) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   data,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *KgtKrywnController) DeletePengajuanKegiatan(ctx *gin.Context) {
	slug := ctx.Param("slug")

	data_kk, err_kk := c.KegiatanKaryawanRepo.FindDataSlug(slug)
	if err_kk == nil {
		c.KegiatanKaryawanRepo.DelKegiatanKaryawanID(data_kk.Slug)

		var kegiatan_id int
		var is_koordinator int
		if data_kk.KoordinatorId == 0 {
			kegiatan_id = data_kk.Id
			is_koordinator = 0
		} else {
			kegiatan_id = data_kk.KoordinatorId
			is_koordinator = 1
		}
		photos := c.KegiatanPhotosRepo.FindDataPhotosID(kegiatan_id, is_koordinator)

		for _, data := range photos {
			c.KegiatanPhotosRepo.DelPhotosID(data.KegiatanId)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *KgtKrywnController) StoreApprovePengajuanKegiatan(ctx *gin.Context) {
	var req Authentication.ValidationApprovalAtasan

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

	kegiatan_karyawan, err := c.KegiatanKaryawanRepo.FindDataSlug(req.SlugKegiatan)
	kegiatan_karyawan.Status = req.Status

	if err == nil {
		c.KegiatanKaryawanRepo.Update(kegiatan_karyawan)
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *KgtKrywnController) ListApprvlKgtKrywn(ctx *gin.Context) {
	var req Authentication.ValidationListApproval
	var list_aprvl []Authentication.ListApprovalTJSL

	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	status := "WaitApv"
	kgt_krwyn, _ := c.KegiatanKaryawanRepo.FindDataNIKCompCodePeriode(req.NIK, req.Tahun, req.CompCode, status)

	for _, data := range kgt_krwyn {
		pihc_mster_krywn, _ := c.PihcMasterKaryRtRepo.FindUserByNIK(data.NIK)

		var kegiatan_id int
		var is_koordinator int
		if data.KoordinatorId == 0 {
			kegiatan_id = data.Id
			is_koordinator = 0
		} else {
			kegiatan_id = data.KoordinatorId
			is_koordinator = 1
		}

		data_kp := c.KegiatanPhotosRepo.FindDataPhotosID(kegiatan_id, is_koordinator)

		pihc_mster_position, _ := c.PihcMasterPositionRepo.FindUserByPosID(pihc_mster_krywn.PosID)

		var jenis_kegiatan string
		if data.KegiatanParentId == 0 {
			jenis_kegiatan = "Kegiatan sosial kemasyarakatan diluar perusahaan"
		} else {
			jenis_kegiatan = "Kegiatan Tanggung Jawab Sosial dan Lingkungan (TJSL) perusahaan"
		}

		data_list := Authentication.ListApprovalTJSL{
			SlugKegiatan:    data.Slug,
			Nik:             data.NIK,
			Nama:            pihc_mster_krywn.Nama,
			PhotoProfile:    "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg",
			Email:           pihc_mster_krywn.Email,
			PosID:           pihc_mster_krywn.PosID,
			PosTitle:        pihc_mster_krywn.PosTitle,
			DeptTitle:       pihc_mster_krywn.DeptTitle,
			Jenis:           jenis_kegiatan,
			NamaKegiatan:    data.NamaKegiatan,
			TanggalKegiatan: data.TanggalKegiatan,
			LokasiKegiatan:  data.LokasiKegiatan,
			Deskripsi:       data.DeskripsiKegiatan,
			Status:          data.Status,
			PhotoKegiatan:   data_kp,
			Short:           pihc_mster_position.Short,
			LogoCompany:     "https://storage.googleapis.com/lumen-oauth-storage/company/logo-pi-full.png",
		}
		list_aprvl = append(list_aprvl, data_list)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"info":   "Success",
		"data":   list_aprvl,
	})
}
