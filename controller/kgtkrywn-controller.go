package controller

import (
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
	var req Authentication.ValidationKK

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

	if req.KegiatanParentId != 0 {
		kk.KegiatanParentId = req.KegiatanParentId
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

				for _, data := range req.Photos {
					kp.KegiatanId = kgt_krywn.Id
					kp.OriginalName = data.OriginalName
					kp.Url = data.URL
					url, _ := c.KegiatanPhotosRepo.GetFileExtensionFromUrl(kp.Url)
					kp.Extendtion = url
					kgt_photos := c.KegiatanPhotosRepo.Create(kp)
					list_id_foto = append(list_id_foto, kgt_photos.Id)
				}

				c.KegiatanPhotosRepo.DelPhotosIDLama(kgt_krywn.Id, list_id_foto)

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
		t := time.Now()
		kk.NIK = req.NIK
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
		kk.Periode = strconv.Itoa(t.Year())
		kk.Slug = users.String(12)
		kk, err_kgtkrywn := c.KegiatanKaryawanRepo.Create(kk)

		if err_kgtkrywn == nil {
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

func (c *KgtKrywnController) ListApprvlKgtKrywn(ctx *gin.Context) {

}

func (c *KgtKrywnController) ShowKgtKrywn(ctx *gin.Context) {
	var data Authentication.KegiatanKaryawanPhotos
	slug := ctx.Param("slug")

	data_kk, err_kk := c.KegiatanKaryawanRepo.FindDataSlug(slug)
	data_kp := c.KegiatanPhotosRepo.FindDataPhotosID(data_kk.Id)
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
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Tidak Ada",
			"Data":   nil,
		})
	}
}

func (c *KgtKrywnController) DeleteKgtKrywn(ctx *gin.Context) {
	slug := ctx.Param("slug")

	data_kk, err_kk := c.KegiatanKaryawanRepo.FindDataSlug(slug)
	if err_kk == nil {
		status := "WaitApv"
		c.KegiatanKaryawanRepo.DelKegiatanKaryawanID(data_kk.Slug, status)
		photos := c.KegiatanPhotosRepo.FindDataPhotosID(data_kk.Id)

		for _, data := range photos {
			c.KegiatanPhotosRepo.DelPhotosID(data.KegiatanId)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Tidak Ada",
			"Data":   nil,
		})
	}
}
