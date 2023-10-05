package tjsl_controller

import (
	"errors"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl"
	users "github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type KoorKgtController struct {
	KegiatanKaryawanRepo    *tjsl.KegiatanKaryawanRepo
	KegiatanKoordinatorRepo *tjsl.KegiatanKoordinatorRepo
	KegiatanPhotosRepo      *tjsl.KegiatanPhotosRepo
	KegiatanMasterRepo      *tjsl.KegiatanMasterRepo
	KoordinatorPersonRepo   *tjsl.KoordinatorPersonRepo
	PihcMasterKaryDbRepo    *pihc.PihcMasterKaryDbRepo
	PihcMasterKaryRtDbRepo  *pihc.PihcMasterKaryRtDbRepo
	PihcMasterKaryRtRepo    *pihc.PihcMasterKaryRtRepo
}

func NewKoorKgtController(Db *gorm.DB, StorageClient *storage.Client) *KoorKgtController {
	return &KoorKgtController{KegiatanKoordinatorRepo: tjsl.NewKegiatanKoordinatorRepo(Db),
		KegiatanKaryawanRepo:   tjsl.NewKegiatanKaryawanRepo(Db, StorageClient),
		KegiatanPhotosRepo:     tjsl.NewKegiatanPhotosRepo(Db),
		KoordinatorPersonRepo:  tjsl.NewKoordinatorPersonRepo(Db),
		PihcMasterKaryDbRepo:   pihc.NewPihcMasterKaryDbRepo(Db),
		PihcMasterKaryRtDbRepo: pihc.NewPihcMasterKaryRtDbRepo(Db),
		PihcMasterKaryRtRepo:   pihc.NewPihcMasterKaryRtRepo(Db),
		KegiatanMasterRepo:     tjsl.NewKegiatanMasterRepo(Db)}
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

	PIHC_MSTR_KRY, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(req.Nik)

	comp_code := PIHC_MSTR_KRY.Company

	if req.Id != 0 {
		kgt_koor, err_kgt_koor := c.KegiatanKoordinatorRepo.FindDataID(req.Id)
		kgt_koor.Nama = req.Nama

		if req.KegiatanParentId != nil {
			kgt_koor.KegiatanParentId = req.KegiatanParentId
		}

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
		if req.KegiatanParentId != nil {
			koorKgt.KegiatanParentId = req.KegiatanParentId
		}
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
	if data_koor.KegiatanParentId != nil {
		data.KegiatanParentID = *data_koor.KegiatanParentId
	}
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
		data_karyawan, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(list_data_person.NIK)

		// data_list[i].Employee.TglLahir = data_karyawan.TglLahir.Format("2006-01-02")
		data_karyawan_convert := convertSourceTargetDataKaryawan(data_karyawan)
		data_list[i] = Authentication.Personal{
			ID:            list_data_person.Id,
			KoordinatorID: list_data_person.KoordinatorId,
			Nik:           list_data_person.NIK,
			CreatedAt:     list_data_person.CreatedAt,
			UpdatedAt:     list_data_person.UpdatedAt,
			Employee:      data_karyawan_convert,
		}

		files, err := c.KegiatanKaryawanRepo.FindPhotosKaryawan(data_karyawan.EmpNo, data_karyawan.Company)
		if err != nil {
			data_list[i].URLPhoto = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
		} else {
			data_list[i].URLPhoto = "https://storage.googleapis.com/" + files
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

func convertSourceTargetDataKaryawan(source pihc.PihcMasterKaryDb) pihc.PihcMasterKary {
	return pihc.PihcMasterKary{
		EmpNo:          source.EmpNo,
		Nama:           source.Nama,
		Gender:         source.Gender,
		Agama:          source.Agama,
		StatusKawin:    source.StatusKawin,
		Anak:           source.Anak,
		Mdg:            "0",
		EmpGrade:       source.EmpGrade,
		EmpGradeTitle:  source.EmpGradeTitle,
		Area:           source.Area,
		AreaTitle:      source.AreaTitle,
		SubArea:        source.SubArea,
		SubAreaTtitle:  source.SubAreaTtitle,
		Contract:       source.Contract,
		Pendidikan:     source.Pendidikan,
		Company:        source.Company,
		Lokasi:         source.Lokasi,
		EmployeeStatus: source.EmployeeStatus,
		Email:          source.Email,
		HP:             source.HP,
		TglLahir:       source.TglLahir.Format("2006-01-02"),
		PosID:          source.PosID,
		PosTitle:       source.PosTitle,
		SubPosID:       source.SubPosID,
		PosGrade:       source.PosGrade,
		PosKategori:    source.PosKategori,
		OrgID:          source.OrgID,
		OrgTitle:       source.OrgTitle,
		DeptID:         source.DeptID,
		DeptTitle:      source.DeptTitle,
		KompID:         source.KompID,
		KompTitle:      source.KompTitle,
		DirID:          source.DirID,
		DirTitle:       source.DirTitle,
		PosLevel:       source.PosLevel,
		SupEmpNo:       source.SupEmpNo,
		BagID:          source.BagID,
		BagTitle:       source.BagTitle,
		SeksiID:        source.SeksiID,
		SeksiTitle:     source.SeksiTitle,
		PreNameTitle:   source.PreNameTitle,
		PostNameTitle:  source.PostNameTitle,
		NoNPWP:         source.NoNPWP,
		BankAccount:    source.BankAccount,
		BankName:       source.BankName,
		MdgDate:        source.MdgDate,
		PayScale:       source.PayScale,
		CCCode:         source.CCCode,
		Nickname:       source.Nickname,
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
	dataKoor := []Authentication.KegiatanListKoordinatorPhotos{}

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
		dataList, _ := c.KegiatanKoordinatorRepo.ListKoordinatorDalamKegiatan(req.Slug, req.NIK)

		for _, item := range dataList {
			// Create a new instance of authentication.KegiatanListKoordinatorPhotos
			koorItem := Authentication.KegiatanListKoordinatorPhotos{
				KegiatanKoordinator: item.KegiatanKoordinator,
				Employee:            item.Employee,
			}

			// Append koorItem to dataKoor
			dataKoor = append(dataKoor, koorItem)
		}
	} else {
		dataList, _ := c.KegiatanKoordinatorRepo.ListKoordinatorLuarKegiatan(req.NIK)

		for _, item := range dataList {
			// Create a new instance of authentication.KegiatanListKoordinatorPhotos
			koorItem := Authentication.KegiatanListKoordinatorPhotos{
				KegiatanKoordinator: item.KegiatanKoordinator,
				Employee:            item.Employee,
			}

			// Append koorItem to dataKoor
			dataKoor = append(dataKoor, koorItem)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"info":   "Success",
		"data":   dataKoor,
	})
}

// func (c *KoorKgtController) ListKoordinator(ctx *gin.Context) {
// 	var req Authentication.ValidationListKoordinator

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

// 	dataList, _ := c.KegiatanKoordinatorRepo.ListKoordinatorLuarKegiatan(req.NIK)

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"status": http.StatusOK,
// 		"info":   "Success",
// 		"data":   dataList,
// 	})
// }
