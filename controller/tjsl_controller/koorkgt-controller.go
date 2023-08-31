package tjsl_controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl"
	users "github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type KoorKgtController struct {
	KegiatanKoordinatorRepo *tjsl.KegiatanKoordinatorRepo
	KegiatanPhotosRepo      *tjsl.KegiatanPhotosRepo
	KegiatanMasterRepo      *tjsl.KegiatanMasterRepo
	KoordinatorPersonRepo   *tjsl.KoordinatorPersonRepo
	PihcMasterKaryRepo      *pihc.PihcMasterKaryRepo
	PihcMasterKaryRtRepo    *pihc.PihcMasterKaryRtRepo
}

func NewKoorKgtController(db *gorm.DB) *KoorKgtController {
	return &KoorKgtController{KegiatanKoordinatorRepo: tjsl.NewKegiatanKoordinatorRepo(db),
		KegiatanPhotosRepo:    tjsl.NewKegiatanPhotosRepo(db),
		KoordinatorPersonRepo: tjsl.NewKoordinatorPersonRepo(db),
		PihcMasterKaryRepo:    pihc.NewPihcMasterKaryRepo(db),
		PihcMasterKaryRtRepo:  pihc.NewPihcMasterKaryRtRepo(db),
		KegiatanMasterRepo:    tjsl.NewKegiatanMasterRepo(db)}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return (fe.Field() + " wajib di isi")
	case "min":
		return ("Peserta yang diundang minimal " + fe.Param() + " orang")
	case "validyear":
		return ("Field has an invalid value: " + fe.Field() + fe.Tag())
	}
	return "Unknown error"
}

func (c *KoorKgtController) StoreKoordinator(ctx *gin.Context) {
	var kp tjsl.KegiatanPhotos
	var koorKgt tjsl.KegiatanKoordinator
	var koorPerson tjsl.KoordinatorPerson
	var req Authentication.ValidationKKoor

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

	PIHC_MSTR_KRY_RT, _ := c.PihcMasterKaryRtRepo.FindUserByNIK(req.Nik)

	comp_code := PIHC_MSTR_KRY_RT.Company

	if req.Id != 0 {
		kgt_koor, err_kgt_koor := c.KegiatanKoordinatorRepo.FindDataID(req.Id)
		kgt_koor.Nama = req.Nama

		if err_kgt_koor == nil {
			kgt_koor, err_updte_kgt_koor := c.KegiatanKoordinatorRepo.Update(kgt_koor)
			if err_updte_kgt_koor == nil {
				var list_id_foto []int

				for _, dataPhotos := range req.Photos {
					kp.KegiatanId = kgt_koor.IdKoordinator
					kp.IsKoordinator = 1
					kp.OriginalName = dataPhotos.OriginalName
					kp.Url = dataPhotos.URL
					// url, _ := c.KegiatanPhotosRepo.GetFileExtensionFromUrl(kp.Url)
					kp.Extendtion = dataPhotos.Extension
					kgt_photos := c.KegiatanPhotosRepo.Create(kp)
					list_id_foto = append(list_id_foto, kgt_photos.Id)
				}

				c.KegiatanPhotosRepo.DelPhotosIDLama(kgt_koor.IdKoordinator, list_id_foto)

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
		koorKgt.KegiatanParentId = 0
		koorKgt.Nama = req.Nama
		koorKgt.CreatedBy = req.Nik
		koorKgt.CompCode = comp_code
		koorKgt.Slug = users.String(12)
		// koorKgt.Periode = strconv.Itoa(t.Year())
		koorKgt.Periode = req.Tahun

		koorKgt, err_koorKgt := c.KegiatanKoordinatorRepo.Create(koorKgt)
		if err_koorKgt == nil {
			for _, dataPhotos := range req.Photos {
				kp.KegiatanId = koorKgt.IdKoordinator
				kp.IsKoordinator = 1
				kp.OriginalName = dataPhotos.OriginalName
				kp.Url = dataPhotos.URL
				// url, _ := c.KegiatanPhotosRepo.GetFileExtensionFromUrl(kp.Url)
				kp.Extendtion = dataPhotos.Extension
				// s := c.KegiatanPhotosRepo.LastString(strings.Split(data.OriginalName, "."))
				// kp.Extendtion = s
				c.KegiatanPhotosRepo.Create(kp)
			}

			req.Person = append(req.Person, koorKgt.CreatedBy)
			for _, dataPerson := range req.Person {
				koorPerson.KoordinatorId = koorKgt.IdKoordinator
				koorPerson.NIK = dataPerson
				c.KoordinatorPersonRepo.Create(koorPerson)
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

func (c *KoorKgtController) ShowDetailKoordinator(ctx *gin.Context) {
	var data Authentication.KegiatanDetailKoordinatorPhotos
	slug := ctx.Param("slug")

	data_koor, err_kgt_koor := c.KegiatanKoordinatorRepo.FindDataSlug(slug)

	data_kp := c.KegiatanPhotosRepo.FindDataPhotosID(data_koor.IdKoordinator, 1)
	data_person := c.KoordinatorPersonRepo.FindDataKoorPersonID(data_koor.IdKoordinator)

	data.IDKoordinator = data_koor.IdKoordinator
	data.KegiatanParentID = 0
	data.Nama = data_koor.Nama
	data.CreatedBy = data_koor.CreatedBy
	data.CreatedAt = data_koor.CreatedAt
	data.UpdatedAt = data_koor.UpdatedAt
	data.CompCode = data_koor.CompCode
	data.Slug = data_koor.Slug
	data.Periode = data_koor.Periode
	data.Photos = data_kp

	// var list_person []Authentication.Personal
	data_list := make([]Authentication.Personal, len(data_person))
	for i, list_data_person := range data_person {
		data_karyawan, _ := c.PihcMasterKaryRepo.FindUserByNIK(list_data_person.NIK)

		data_list[i] = Authentication.Personal{
			ID:            list_data_person.Id,
			KoordinatorID: list_data_person.KoordinatorId,
			Nik:           list_data_person.NIK,
			CreatedAt:     list_data_person.CreatedAt,
			UpdatedAt:     list_data_person.UpdatedAt,
			Employee:      data_karyawan,
			URLPhoto:      "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg",
		}
		// data_list := Authentication.Personal{
		// 	ID:            list_data_person.Id,
		// 	KoordinatorID: list_data_person.KoordinatorId,
		// 	Nik:           list_data_person.NIK,
		// 	CreatedAt:     list_data_person.CreatedAt,
		// 	UpdatedAt:     list_data_person.UpdatedAt,
		// 	Employee:      data_karyawan,
		// 	URLPhoto:      "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg",
		// }
		// list_person = append(list_person, data_list)
	}
	data.Person = data_list

	if err_kgt_koor == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   data,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *KoorKgtController) DeleteKoordinator(ctx *gin.Context) {
	slug := ctx.Param("slug")

	data_koor_kgt, err_koor_kgt := c.KegiatanKoordinatorRepo.FindDataSlug(slug)

	if err_koor_kgt == nil {
		c.KegiatanKoordinatorRepo.DelKegiatanKoordinatorID(data_koor_kgt.Slug)
		photos := c.KegiatanPhotosRepo.FindDataPhotosID(data_koor_kgt.IdKoordinator, 1)
		person := c.KoordinatorPersonRepo.FindDataKoorPersonID(data_koor_kgt.IdKoordinator)

		for _, dataPhotos := range photos {
			c.KegiatanPhotosRepo.DelPhotosID(dataPhotos.KegiatanId)
		}

		for _, dataPerson := range person {
			c.KoordinatorPersonRepo.DelPersonID(dataPerson.KoordinatorId)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *KoorKgtController) ListKoordinator(ctx *gin.Context) {
	var req Authentication.ValidationListKoordinator
	var dataKoor []Authentication.KegiatanListKoordinatorPhotos
	// var listKoorKegiatan []tjsl.KegiatanKoordinator

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

	if req.Slug != "" {
		kegiatanMster, errKgtmster := c.KegiatanMasterRepo.FindDataBySlug(req.Slug)

		if errKgtmster == nil {
			koorKegiatan, errkoorKgt := c.KegiatanKoordinatorRepo.FindDataParentID(kegiatanMster.IdKegiatan, req.NIK)

			if errkoorKgt == nil {
				for _, dataKoorKgt := range koorKegiatan {
					data_karyawan, _ := c.PihcMasterKaryRtRepo.FindUserByNIK(dataKoorKgt.CreatedBy)
					data_list := Authentication.KegiatanListKoordinatorPhotos{
						IDKoordinator:    dataKoorKgt.IdKoordinator,
						KegiatanParentID: dataKoorKgt.KegiatanParentId,
						Nama:             dataKoorKgt.Nama,
						CreatedBy:        dataKoorKgt.CreatedBy,
						CreatedAt:        dataKoorKgt.CreatedAt,
						UpdatedAt:        dataKoorKgt.UpdatedAt,
						CompCode:         dataKoorKgt.CompCode,
						Slug:             dataKoorKgt.Slug,
						Periode:          dataKoorKgt.Periode,
						Employee:         data_karyawan}
					dataKoor = append(dataKoor, data_list)
				}
			}
		}
	} else {
		dataList, errDataList := c.KegiatanKoordinatorRepo.FindDataKoorIDLuarKegiatan(req.NIK)

		if errDataList == nil {
			for _, dataKoorPerson := range dataList {
				data_karyawan, _ := c.PihcMasterKaryRtRepo.FindUserByNIK(dataKoorPerson.CreatedBy)
				listKoor := Authentication.KegiatanListKoordinatorPhotos{
					IDKoordinator:    dataKoorPerson.IdKoordinator,
					KegiatanParentID: dataKoorPerson.KegiatanParentId,
					Nama:             dataKoorPerson.Nama,
					CreatedBy:        dataKoorPerson.CreatedBy,
					CreatedAt:        dataKoorPerson.CreatedAt,
					UpdatedAt:        dataKoorPerson.UpdatedAt,
					CompCode:         dataKoorPerson.CompCode,
					Slug:             dataKoorPerson.Slug,
					Periode:          dataKoorPerson.Periode,
					Employee:         data_karyawan}
				dataKoor = append(dataKoor, listKoor)
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"info":   "Success",
		"data":   dataKoor,
	})
}
