package tjsl_controller

import (
	"errors"
	"fmt"
	"net/http"
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
		parsedTime, _ := time.Parse(time.DateTime, req.TanggalKegiatan)
		kgt_krywn.TanggalKegiatan = parsedTime
		kgt_krywn.LokasiKegiatan = req.LokasiKegiatan
		if req.DeskripsiKegiatan != nil {
			kgt_krywn.DeskripsiKegiatan = req.DeskripsiKegiatan
		}

		if err_kgtkrywn == nil {
			kgt_krywn, err_updte_kgtkrywn := c.KegiatanKaryawanRepo.Update(kgt_krywn)
			if err_updte_kgtkrywn == nil {
				var list_id_foto []int
				var kegiatan_id int
				var is_koordinator int

				if kgt_krywn.KoordinatorId == nil {
					kegiatan_id = kgt_krywn.Id
					is_koordinator = 0
				} else {
					kegiatan_id = *kgt_krywn.KoordinatorId
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
		parsedTime, _ := time.Parse(time.DateTime, req.TanggalKegiatan)
		kk.TanggalKegiatan = parsedTime
		kk.LokasiKegiatan = req.LokasiKegiatan
		kk.DeskripsiKegiatan = req.DeskripsiKegiatan
		kk.Status = "WaitApv"
		kk.Manager = nil
		kk.CompCode = comp_code
		// kk.Periode = strconv.Itoa(t.Year())
		kk.Periode = req.Tahun
		kk.Slug = users.String(12)
		kk, err_kgtkrywn := c.KegiatanKaryawanRepo.Create(kk)

		var kegiatan_id *int
		var is_koordinator int
		if kk.KoordinatorId == nil {
			kegiatan_id = &kk.Id
			is_koordinator = 0
		} else {
			kegiatan_id = kk.KoordinatorId
			is_koordinator = 1
		}

		if err_kgtkrywn == nil {
			for _, data := range req.Photos {
				if kegiatan_id != nil {
					kp.KegiatanId = *kegiatan_id
				}
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

// func (c *KgtKrywnController) GetChartSummary(ctx *gin.Context) {
// 	var req Authentication.ValidationMyTjsl

// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		var ve validator.ValidationErrors
// 		if errors.As(err, &ve) {
// 			out := make([]Authentication.ErrorMsg, len(ve))
// 			for i, fe := range ve {
// 				out[i] = Authentication.ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
// 			}
// 			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": out})
// 		}
// 		return
// 	}

// 	status := "Approved"
// 	bulan := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
// 	var total int
// 	list_bulan := &Authentication.Month{}
// 	var data Authentication.ListChartSummary
// 	for _, intBulan := range bulan {
// 		jumlahPerbulan, _ := c.KegiatanKaryawanRepo.RekapPerbulan(req.Nik, req.Tahun, status, intBulan)
// 		switch intBulan {
// 		case 1:
// 			list_bulan.Num1 = jumlahPerbulan
// 		case 2:
// 			list_bulan.Num2 = jumlahPerbulan
// 		case 3:
// 			list_bulan.Num3 = jumlahPerbulan
// 		case 4:
// 			list_bulan.Num4 = jumlahPerbulan
// 		case 5:
// 			list_bulan.Num5 = jumlahPerbulan
// 		case 6:
// 			list_bulan.Num6 = jumlahPerbulan
// 		case 7:
// 			list_bulan.Num7 = jumlahPerbulan
// 		case 8:
// 			list_bulan.Num8 = jumlahPerbulan
// 		case 9:
// 			list_bulan.Num9 = jumlahPerbulan
// 		case 10:
// 			list_bulan.Num10 = jumlahPerbulan
// 		case 11:
// 			list_bulan.Num11 = jumlahPerbulan
// 		case 12:
// 			list_bulan.Num12 = jumlahPerbulan
// 		default:
// 			fmt.Println("Tidak MASUK")
// 		}
// 		total += jumlahPerbulan
// 		// data, err := c.KegiatanKaryawanRepo.FindDataNIKPeriode(req.Nik, req.Tahun)
// 	}
// 	data.Month = *list_bulan
// 	data.TotalIndividu = total

//		ctx.JSON(http.StatusOK, gin.H{
//			"status": http.StatusOK,
//			"info":   "Success",
//			"data":   data,
//		})
//	}
func (c *KgtKrywnController) GetChartSummary(ctx *gin.Context) {
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

	status := "Approved"
	list_bulan := &Authentication.RekapPerbulan{}
	var data Authentication.ListChartSummary
	karyawan := &Authentication.Employee{}

	dataKaryawan, _ := c.PihcMasterKaryRtRepo.FindUserRekapByNIK(req.Nik)

	if dataKaryawan == nil {
		var data2 Authentication.ListChartNotFoundDataSummary
		data2.RekapPerbulan.Month = list_bulan.Month
		data2.RekapPerbulan.TotalIndividu = list_bulan.TotalIndividu

		fmt.Println("NIL")
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   data2,
		})
	} else {
		karyawan.EmpNama = dataKaryawan.EmpNama
		karyawan.Nik = dataKaryawan.Nik
		karyawan.PosID = dataKaryawan.PosID
		karyawan.PosTitle = dataKaryawan.PosTitle
		karyawan.DeptID = dataKaryawan.DeptID
		karyawan.DeptTitle = dataKaryawan.DeptTitle
		karyawan.KompID = dataKaryawan.KompID
		karyawan.KompTitle = dataKaryawan.KompTitle
		karyawan.DirID = dataKaryawan.DirID
		karyawan.DirTitle = dataKaryawan.DirTitle
		karyawan.Photo = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"

		jumlahPerbulan, _ := c.KegiatanKaryawanRepo.RekapPerbulan(req.Nik, req.Tahun, status)
		for _, dataPerbulan := range jumlahPerbulan {
			if dataPerbulan.Bulan == 1 {
				list_bulan.Num1 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 2 {
				list_bulan.Num2 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 3 {
				list_bulan.Num3 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 4 {
				list_bulan.Num4 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 5 {
				list_bulan.Num5 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 6 {
				list_bulan.Num6 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 7 {
				list_bulan.Num7 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 8 {
				list_bulan.Num8 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 9 {
				list_bulan.Num9 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 10 {
				list_bulan.Num10 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 11 {
				list_bulan.Num11 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.Bulan == 12 {
				list_bulan.Num12 = dataPerbulan.JumlahPerbulan
			}
			if dataPerbulan.TotalPertahun != 0 {
				fmt.Println("masuk")
				list_bulan.TotalIndividu = dataPerbulan.TotalPertahun
			}
		}
		data.RekapPerbulan.Month = list_bulan.Month
		data.RekapPerbulan.TotalIndividu = list_bulan.TotalIndividu
		data.Employee = *karyawan

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   data,
		})
	}
}

func (c *KgtKrywnController) ShowDetailPengajuanKegiatan(ctx *gin.Context) {
	var data Authentication.KegiatanKaryawanPhotos
	slug := ctx.Param("slug")

	data_kk, err_kk := c.KegiatanKaryawanRepo.FindDataSlug(slug)

	var kegiatan_id int
	var is_koordinator int
	if data_kk.KoordinatorId == nil {
		kegiatan_id = data_kk.Id
		is_koordinator = 0
	} else {
		kegiatan_id = *data_kk.KoordinatorId
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

	if data_kk.KegiatanParentId == nil {
		data.Jenis = "Kegiatan sosial kemasyarakatan diluar perusahaan"
	} else {
		data.Jenis = "Kegiatan Tanggung Jawab Sosial dan Lingkungan (TJSL) perusahaan"
	}

	if data_kk.KoordinatorId != nil {
		data.KoordinatorID = data_kk.KoordinatorId
	}

	data.SlugKoordinator = nil
	data.SlugKegiatanParent = nil
	if data_kk.KegiatanParentId != nil {
		data.KegiatanParentID = data_kk.KegiatanParentId
	}
	data.NamaKegiatan = data_kk.NamaKegiatan

	// rfc339, _ := time.Parse(time.RFC3339, data_kk.TanggalKegiatan)
	// tgl_kegiatan_nonformat := rfc339.Format(time.DateOnly)
	// year, month, day := rfc339.Date()
	// tanggal := strconv.Itoa(day)
	// bulan := month.String()
	// tahun := strconv.Itoa(year)
	// tgl_kegiatan := tanggal + " " + bulan + " " + tahun

	data.TanggalKegiatan = data_kk.TanggalKegiatan.Format("02 January 2006")
	data.TanggalKegiatanNonFormat = data_kk.TanggalKegiatan.Format("2006-01-02")
	data.LokasiKegiatan = data_kk.LokasiKegiatan
	if data_kk.DeskripsiKegiatan != nil {
		data.Deskripsi = data_kk.DeskripsiKegiatan
	}
	data.Status = data_kk.Status
	if data_kk.DescDecline != nil {
		data.AlasanPenolakan = data_kk.DescDecline
	}
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
		if data_kk.KoordinatorId == nil {
			kegiatan_id = data_kk.Id
			is_koordinator = 0
		} else {
			kegiatan_id = *data_kk.KoordinatorId
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

// func (c *KgtKrywnController) ListApprvlKgtKrywn(ctx *gin.Context) {
// 	var req Authentication.ValidationListApproval
// 	list_aprvl := []Authentication.ListApprovalTJSL{}

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		var ve validator.ValidationErrors
// 		if errors.As(err, &ve) {
// 			out := make([]Authentication.ErrorMsg, len(ve))
// 			for i, fe := range ve {
// 				out[i] = Authentication.ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
// 			}
// 			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": out})
// 		}
// 		return
// 	}

// 	status := "WaitApv"
// 	kgt_krwyn, _ := c.KegiatanKaryawanRepo.FindDataNIKCompCodePeriode(req.NIK, req.Tahun, req.CompCode, status)

// 	for _, data := range kgt_krwyn {
// 		pihc_mster_krywn, _ := c.PihcMasterKaryRtRepo.FindUserByNIK(data.NIK)

// 		var kegiatan_id int
// 		var is_koordinator int
// 		if data.KoordinatorId == nil || *data.KoordinatorId == 0 {
// 			kegiatan_id = data.Id
// 			is_koordinator = 0
// 		} else {
// 			kegiatan_id = *data.KoordinatorId
// 			is_koordinator = 1
// 		}

// 		data_kp := c.KegiatanPhotosRepo.FindDataPhotosID(kegiatan_id, is_koordinator)

// 		// pihc_mster_position, _ := c.PihcMasterPositionRepo.FindUserByPosID(pihc_mster_krywn.PosID)

// 		var jenis_kegiatan string
// 		if data.KegiatanParentId == nil {
// 			jenis_kegiatan = "Kegiatan sosial kemasyarakatan diluar perusahaan"
// 		} else {
// 			jenis_kegiatan = "Kegiatan Tanggung Jawab Sosial dan Lingkungan (TJSL) perusahaan"
// 		}

// 		data_list := Authentication.ListApprovalTJSL{
// 			SlugKegiatan:    data.Slug,
// 			Nik:             data.NIK,
// 			Nama:            pihc_mster_krywn.Nama,
// 			PhotoProfile:    "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg",
// 			Email:           pihc_mster_krywn.Email,
// 			PosID:           pihc_mster_krywn.PosID,
// 			PosTitle:        pihc_mster_krywn.PosTitle,
// 			DeptTitle:       pihc_mster_krywn.DeptTitle,
// 			Jenis:           jenis_kegiatan,
// 			NamaKegiatan:    data.NamaKegiatan,
// 			TanggalKegiatan: data.TanggalKegiatan.Format("02 January 2006"),
// 			LokasiKegiatan:  data.LokasiKegiatan,
// 			Deskripsi:       data.DeskripsiKegiatan,
// 			Status:          data.Status,
// 			PhotoKegiatan:   data_kp,
// 			Short:           nil,
// 			LogoCompany:     "https://storage.googleapis.com/lumen-oauth-storage/company/logo-pi-full.png",
// 		}
// 		// Short:           pihc_mster_position.Short,
// 		list_aprvl = append(list_aprvl, data_list)
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"status": http.StatusOK,
// 		"info":   "Success",
// 		"data":   list_aprvl,
// 	})
// }

func (c *KgtKrywnController) ListApprvlKgtKrywn(ctx *gin.Context) {
	var req Authentication.ValidationListApproval
	list_aprvl := []Authentication.ListApprovalTJSL{}

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
	kgt_krywn, _ := c.KegiatanKaryawanRepo.ListKegiatanKaryawanApprvalWait(req.NIK, req.Tahun, req.CompCode, status)
	// kgt_krwyn, _ := c.KegiatanKaryawanRepo.FindDataNIKCompCodePeriode(req.NIK, req.Tahun, req.CompCode, status)

	for _, data := range kgt_krywn {
		// 	pihc_mster_krywn, _ := c.PihcMasterKaryRtRepo.FindUserByNIK(data.NIK)

		var kegiatan_id int
		var is_koordinator int
		if data.KoordinatorId == nil || *data.KoordinatorId == 0 {
			kegiatan_id = data.Id
			is_koordinator = 0
		} else {
			kegiatan_id = *data.KoordinatorId
			is_koordinator = 1
		}

		data_kp := c.KegiatanPhotosRepo.FindDataPhotosID(kegiatan_id, is_koordinator)

		// 	// pihc_mster_position, _ := c.PihcMasterPositionRepo.FindUserByPosID(pihc_mster_krywn.PosID)

		var jenis_kegiatan string
		if data.KegiatanParentId == nil {
			jenis_kegiatan = "Kegiatan sosial kemasyarakatan diluar perusahaan"
		} else {
			jenis_kegiatan = "Kegiatan Tanggung Jawab Sosial dan Lingkungan (TJSL) perusahaan"
		}

		data_list := Authentication.ListApprovalTJSL{
			SlugKegiatan:    data.Slug,
			Nik:             data.NIK,
			Nama:            data.Nama,
			PhotoProfile:    "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg",
			Email:           data.Email,
			PosID:           data.PosID,
			PosTitle:        data.PosTitle,
			DeptTitle:       data.DeptTitle,
			Jenis:           jenis_kegiatan,
			NamaKegiatan:    data.NamaKegiatan,
			TanggalKegiatan: data.TanggalKegiatan.Format("02 January 2006"),
			LokasiKegiatan:  data.LokasiKegiatan,
			Deskripsi:       data.DeskripsiKegiatan,
			Status:          data.Status,
			PhotoKegiatan:   data_kp,
			Short:           nil,
			LogoCompany:     "https://storage.googleapis.com/lumen-oauth-storage/company/logo-pi-full.png",
		}
		// Short:           pihc_mster_position.Short,
		list_aprvl = append(list_aprvl, data_list)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"info":   "Success",
		"data":   list_aprvl,
	})
}
