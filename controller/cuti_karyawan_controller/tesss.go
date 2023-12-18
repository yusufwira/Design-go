package cuti_karyawan_controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/cuti"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"gorm.io/gorm"
)

type TesssController struct {
	PengajuanAbsenRepo        *cuti.PengajuanAbsenRepo
	HistoryPengajuanAbsenRepo *cuti.HistoryPengajuanAbsenRepo
	SaldoCutiRepo             *cuti.SaldoCutiRepo
	HistorySaldoCutiRepo      *cuti.HistorySaldoCutiRepo
	TipeAbsenRepo             *cuti.TipeAbsenRepo
	FileAbsenRepo             *cuti.FileAbsenRepo
	TransaksiCutiRepo         *cuti.TransaksiCutiRepo
	PihcMasterKaryDbRepo      *pihc.PihcMasterKaryDbRepo
	PihcMasterCompanyRepo     *pihc.PihcMasterCompanyRepo
}

func NewTesssController(Db *gorm.DB) *TesssController {
	return &TesssController{PengajuanAbsenRepo: cuti.NewPengajuanAbsenRepo(Db),
		HistoryPengajuanAbsenRepo: cuti.NewHistoryPengajuanAbsenRepo(Db),
		SaldoCutiRepo:             cuti.NewSaldoCutiRepo(Db),
		HistorySaldoCutiRepo:      cuti.NewHistorySaldoCutiRepo(Db),
		TipeAbsenRepo:             cuti.NewTipeAbsenRepo(Db),
		FileAbsenRepo:             cuti.NewFileAbsenRepo(Db),
		TransaksiCutiRepo:         cuti.NewTransaksiCutiRepo(Db),
		PihcMasterKaryDbRepo:      pihc.NewPihcMasterKaryDbRepo(Db),
		PihcMasterCompanyRepo:     pihc.NewPihcMasterCompanyRepo(Db)}
}

// History
func saldoHistory(source cuti.SaldoCuti) cuti.HistorySaldoCuti {
	return cuti.HistorySaldoCuti{
		IdHistorySaldoCuti: source.IdSaldoCuti,
		TipeAbsenId:        source.TipeAbsenId,
		Nik:                source.Nik,
		Saldo:              source.Saldo,
		ValidFrom:          source.ValidFrom,
		ValidTo:            source.ValidTo,
		CreatedBy:          source.CreatedBy,
		CreatedAt:          source.CreatedAt,
		UpdatedAt:          source.UpdatedAt,
		Periode:            source.Periode,
		MaxHutang:          source.MaxHutang,
		ValidFromHutang:    source.ValidFromHutang,
	}
}
func pengajuanHistory(source cuti.PengajuanAbsen) cuti.HistoryPengajuanAbsen {
	return cuti.HistoryPengajuanAbsen{
		IDHistoryPengajuanAbsen: source.IdPengajuanAbsen,
		Nik:                     source.Nik,
		CompCode:                source.CompCode,
		TipeAbsenId:             source.TipeAbsenId,
		Deskripsi:               source.Deskripsi,
		MulaiAbsen:              source.MulaiAbsen,
		AkhirAbsen:              source.AkhirAbsen,
		TglPengajuan:            source.TglPengajuan,
		Status:                  source.Status,
		CreatedBy:               source.CreatedBy,
		Keterangan:              source.Keterangan,
		Periode:                 source.Periode,
		ApprovedBy:              source.ApprovedBy,
		JmlHariKalendar:         source.JmlHariKalendar,
		JmlHariKerja:            source.JmlHariKerja,
	}
}
func perhitungan(mulai time.Time, akhir time.Time) (int, int) {
	jmlhHariKalender := 0
	JmlHariKerja := 0
	for currentDate := mulai; !currentDate.After(akhir); currentDate = currentDate.AddDate(0, 0, 1) {
		jmlhHariKalender++
		if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
			JmlHariKerja++
		}
	}
	return jmlhHariKalender, JmlHariKerja
}
func perhitunganx(mulai time.Time, saldo int) (int, int) {
	jmlhHariKalender := 0
	JmlHariKerja := 0
	for currentDate := mulai; JmlHariKerja < saldo; currentDate = currentDate.AddDate(0, 0, 1) {
		jmlhHariKalender++
		if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
			JmlHariKerja++
		}
	}
	return jmlhHariKalender, JmlHariKerja
}

// Pengajuan Cuti
func (c *TesssController) StoreCutiKaryawan(ctx *gin.Context) {
	var req Authentication.ValidasiStoreCutiKaryawan
	var sck cuti.PengajuanAbsen
	var fsc []cuti.FileAbsen
	var trsc []Authentication.SaldoCutiTransaksiPengajuan
	var history_saldo []cuti.HistorySaldoCuti

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

	PIHC_MSTR_KRY, err_krywn := c.PihcMasterKaryDbRepo.FindUserByNIK(req.Nik)
	comp_code := PIHC_MSTR_KRY.Company

	if err_krywn == nil {
		if req.IdPengajuanAbsen == 0 {
			sck.Nik = req.Nik
			sck.TipeAbsenId = &req.TipeAbsenId
			sck.CompCode = comp_code
			sck.Deskripsi = &req.Deskripsi
			sck.MulaiAbsen, _ = time.Parse(time.DateOnly, req.MulaiAbsen)
			sck.AkhirAbsen, _ = time.Parse(time.DateOnly, req.AkhirAbsen)
			sck.TglPengajuan, _ = time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
			stats := "WaitApproved"
			sck.Status = &stats
			periode := strconv.Itoa(time.Now().Year())
			sck.Periode = &periode
			sck.CreatedBy = &req.CreatedBy
			dataKaryawan, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(sck.Nik)
			if dataKaryawan.PosTitle != "Wakil Direktur Utama" {
				for dataKaryawan.PosTitle != "Wakil Direktur Utama" {
					dataKaryawan, _ = c.PihcMasterKaryDbRepo.FindUserAtasanBySupPosID(dataKaryawan.SupPosID)
				}
			} else {
				for dataKaryawan.PosTitle != "Direktur Utama" {
					dataKaryawan, _ = c.PihcMasterKaryDbRepo.FindUserAtasanBySupPosID(dataKaryawan.SupPosID)
				}
			}
			approvedBy := dataKaryawan.EmpNo
			sck.ApprovedBy = &approvedBy

			// Tipe Absen
			tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*sck.TipeAbsenId)

			if tipeAbsen.MaxAbsen != nil {
				// Tanpa Menggunakan Saldo
				jmlhHariKalender, JmlHariKerja := perhitungan(sck.MulaiAbsen, sck.AkhirAbsen)
				sck.JmlHariKalendar = &jmlhHariKalender
				sck.JmlHariKerja = &JmlHariKerja

				transaksi_cuti := cuti.TransaksiCuti{}
				sckData, _ := c.PengajuanAbsenRepo.Create(sck)
				if *tipeAbsen.TipeMaxAbsen == "hari_kalender" {
					transaksi_cuti.TipeHari = "hari_kalender"
					transaksi_cuti.JumlahCuti = jmlhHariKalender
				} else if *tipeAbsen.TipeMaxAbsen == "hari_kerja" {
					transaksi_cuti.TipeHari = "hari_kerja"
					transaksi_cuti.JumlahCuti = JmlHariKerja
				}

				convert := convertSourceTargetMyPengajuanAbsen(sckData, tipeAbsen)
				for _, fa := range req.FileAbsen {
					files := cuti.FileAbsen{
						PengajuanAbsenId: sckData.IdPengajuanAbsen,
						Filename:         fa.Filename,
						Url:              fa.URL,
						Extension:        fa.Extension,
					}
					fsc = append(fsc, files)
				}
				// CREATE FileAbsen
				files, _ := c.FileAbsenRepo.CreateArr(fsc)

				// Transaksi Cuti
				transaksi_cuti.PengajuanAbsenId = sckData.IdPengajuanAbsen
				transaksi_cuti.Nik = sckData.Nik
				if sckData.Periode != nil {
					transaksi_cuti.Periode = *sckData.Periode
				}

				// CREATE Transaksi Cuti
				c.TransaksiCutiRepo.Create(transaksi_cuti)

				if files == nil {
					files = []cuti.FileAbsen{}
				}
				data := Authentication.PengajuanAbsens{
					MyPengajuanAbsen: convert,
					File:             files,
				}

				ctx.JSON(http.StatusOK, gin.H{
					"status": http.StatusOK,
					"data":   data,
				})
			} else {
				// Menggunakan Saldo
				saldoPeriode, _ := c.SaldoCutiRepo.SaldoPeriode(req.TipeAbsenId, comp_code, sck.Nik, periode)
				isSaldo := false
				var totalKerja, totalKalender int = 0, 0
				keterangan := "Berada Diluar Masa Berlaku"
				fmt.Println(keterangan)

				jmlhHariKalender, JmlHariKerja := perhitungan(sck.MulaiAbsen, sck.AkhirAbsen)
				for _, xx := range saldoPeriode {
					x, y := perhitunganx(sck.MulaiAbsen, xx.Saldo)
					fmt.Println(x, y)
					fmt.Println(sck.MulaiAbsen.AddDate(0, 0, x-1).Format(time.DateOnly))
				}
				fmt.Println("--------------------")
				fmt.Println(jmlhHariKalender, JmlHariKerja)
				fmt.Println(sck.MulaiAbsen.AddDate(0, 0, jmlhHariKalender-1).Format(time.DateOnly))

				sck.JmlHariKalendar = &jmlhHariKalender
				sck.JmlHariKerja = &JmlHariKerja

				if comp_code == "A000" {
					for _, saldo := range saldoPeriode {
						if (sck.MulaiAbsen.After(saldo.ValidFrom) || sck.MulaiAbsen.Equal(saldo.ValidFrom)) &&
							(sck.AkhirAbsen.Before(saldo.ValidTo) || sck.AkhirAbsen.Equal(saldo.ValidTo)) {
							// MulaiAbsen >= ValidFrom && AkhirAbsen <= ValidTo
							hariKalender, hariKerja := perhitungan(sck.MulaiAbsen, sck.AkhirAbsen)

							if saldo.Saldo != 0 {
								if hariKerja <= saldo.Saldo {
									saldo.Saldo = saldo.Saldo - hariKerja
									isSaldo = true
								} else if hariKerja <= (saldo.Saldo+saldo.MaxHutang) && (saldo.MaxHutang != 0) {
									hutang := hariKerja - saldo.Saldo
									saldo.Saldo = 0
									saldo.MaxHutang = saldo.MaxHutang - hutang
									tahun, _ := strconv.Atoi(saldo.Periode)
									saldoNextYear, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipe(saldo.Nik, saldo.TipeAbsenId, strconv.Itoa(tahun+1))
									if saldoNextYear.Saldo-hutang >= 0 {
										saldoNextYear.Saldo = saldoNextYear.Saldo - hutang
										c.SaldoCutiRepo.Update(saldoNextYear)
										dataHistorySaldoCuti := saldoHistory(saldoNextYear)
										history_saldo = append(history_saldo, dataHistorySaldoCuti)
										isSaldo = true
									} else {
										keterangan = "Saldo untuk Periode Berikutnya Tidak Cukup"
									}
								}
							} else {
								keterangan = "Saldo Tidak Cukup"
							}
							if isSaldo {
								totalKerja += hariKerja
								totalKalender += hariKalender
								source := Authentication.SaldoCutiTransaksiPengajuan{
									SaldoCuti: saldo,
									JmlhCuti:  hariKerja,
								}
								trsc = append(trsc, source)
								dataHistorySaldoCuti := saldoHistory(saldo)
								history_saldo = append(history_saldo, dataHistorySaldoCuti)
							}
						} else {
							isSaldo = false
						}
					}
				} else {
					var newPeriode time.Time
					for _, saldo := range saldoPeriode {
						if (sck.MulaiAbsen.Before(saldo.ValidTo) || sck.MulaiAbsen.Equal(saldo.ValidTo)) &&
							sck.MulaiAbsen.After(saldo.ValidFrom) && (sck.AkhirAbsen.After(saldo.ValidTo) || sck.AkhirAbsen.Equal(saldo.ValidTo)) {
							// MulaiAbsen <= ValidTo && MulaiAbsen > ValidFrom && AkhirAbsen >= ValidTo
							// hariKalender, hariKerja := perhitungan(sck.MulaiAbsen, saldo.ValidTo.AddDate(0, 0, -1))
							// for currentDate := sck.MulaiAbsen; !currentDate.After(saldo.ValidTo.AddDate(0, 0, -1)); currentDate = currentDate.AddDate(0, 0, 1) {
							// 	hariKalender++
							// 	if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
							// 		hariKerja++
							// 	}
							// }

							isSaldo = true
							newPeriode = saldo.ValidTo
						} else if (sck.MulaiAbsen.After(saldo.ValidFrom) || sck.MulaiAbsen.Equal(saldo.ValidFrom)) &&
							(sck.AkhirAbsen.Before(saldo.ValidTo) || sck.AkhirAbsen.Equal(saldo.ValidTo)) {
							// MulaiAbsen >= ValidFrom && AkhirAbsen <= ValidTo
							isSaldo = true
						} else if (newPeriode.After(saldo.ValidFrom) || newPeriode.Equal(saldo.ValidFrom)) &&
							(sck.AkhirAbsen.Before(saldo.ValidTo) || sck.AkhirAbsen.Equal(saldo.ValidTo)) {
							// newPeriode>=ValidFrom && AkhirAbsen <= ValidTo(periode ke-2)
							isSaldo = true
						} else {
							isSaldo = false
						}
					}
				}
				// CREATE
				// if isSaldo {
				// 	if JmlHariKerja == totalKerja {
				// 		sckData, _ := c.PengajuanAbsenRepo.Create(sck)
				// 		dataHistoryPengajuanAbsen := pengajuanHistory(sckData)
				// 		c.HistoryPengajuanAbsenRepo.Create(dataHistoryPengajuanAbsen)

				// 		convert := convertSourceTargetMyPengajuanAbsen(sckData, tipeAbsen)
				// 		for _, fa := range req.FileAbsen {
				// 			files := cuti.FileAbsen{
				// 				PengajuanAbsenId: sckData.IdPengajuanAbsen,
				// 				Filename:         fa.Filename,
				// 				Url:              fa.URL,
				// 				Extension:        fa.Extension,
				// 			}
				// 			fsc = append(fsc, files)
				// 		}
				// 		files, _ := c.FileAbsenRepo.CreateArr(fsc)
				// 		for _, transaction := range trsc {
				// 			c.SaldoCutiRepo.Update(transaction.SaldoCuti)
				// 			transaksi_cuti := cuti.TransaksiCuti{
				// 				PengajuanAbsenId: sckData.IdPengajuanAbsen,
				// 				Nik:              sckData.Nik,
				// 				Periode:          transaction.Periode,
				// 				JumlahCuti:       transaction.JmlhCuti,
				// 				TipeHari:         "hari_kerja",
				// 			}
				// 			c.TransaksiCutiRepo.Create(transaksi_cuti)
				// 		}
				// 		for _, historySaldo := range history_saldo {
				// 			c.HistorySaldoCutiRepo.Create(historySaldo)
				// 		}
				// 		if files == nil {
				// 			files = []cuti.FileAbsen{}
				// 		}
				// 		data := Authentication.PengajuanAbsens{
				// 			MyPengajuanAbsen: convert,
				// 			File:             files,
				// 		}

				// 		ctx.JSON(http.StatusOK, gin.H{
				// 			"status": http.StatusOK,
				// 			"data":   data,
				// 		})
				// 	} else {
				// 		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				// 			"status":     http.StatusServiceUnavailable,
				// 			"keterangan": keterangan,
				// 		})
				// 	}
				// } else {
				// 	ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				// 		"status":     http.StatusServiceUnavailable,
				// 		"keterangan": keterangan,
				// 	})
				// }
			}
		} else {
			// Update
		}
	} else {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}

}

// APPROVAL Cuti
func (c *TesssController) ListApprvlCuti(ctx *gin.Context) {
	var req Authentication.ValidationNIKTahunStatus
	list_aprvl := []Authentication.ListApprovalCuti{}

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

	var dataDB []cuti.PengajuanAbsen
	var err error

	if req.Status == "WaitApproved" {
		dB, errs := c.PengajuanAbsenRepo.FindDataNIKPeriodeApprovalWaiting(req.NIK, req.Tahun, req.Status, req.IsManager)
		err = errs
		dataDB = dB
	} else {
		dB, errs := c.PengajuanAbsenRepo.FindDataNIKPeriodeApproval(req.NIK, req.Tahun, req.IsManager)
		err = errs
		dataDB = dB
	}

	var arrNIK []string
	var arrTipeAbsenID []string
	var arrCompany []string
	var arrIdPengajuanAbsen []int

	if err == nil {
		for _, myCuti := range dataDB {
			arrIdPengajuanAbsen = append(arrIdPengajuanAbsen, myCuti.IdPengajuanAbsen)
			arrNIK = append(arrNIK, myCuti.Nik)
			arrTipeAbsenID = append(arrTipeAbsenID, *myCuti.TipeAbsenId)
		}
		karyawan, _ := c.PihcMasterKaryDbRepo.FindUserByNIKArray(arrNIK)
		tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByIDArray(arrTipeAbsenID)
		files, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuanArray(arrIdPengajuanAbsen)
		for _, myKrywn := range karyawan {
			arrCompany = append(arrCompany, myKrywn.Company)
		}
		companys, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompanyArray(arrCompany)

		for _, myCuti := range dataDB {
			myFiles := []cuti.FileAbsen{}
			list_pengajuan := Authentication.ListApprovalCuti{}
			// Karyawan
			for _, myKaryawan := range karyawan {
				if myCuti.Nik == myKaryawan.EmpNo {
					for _, myCompany := range companys {
						if myKaryawan.Company == myCompany.Code {
							list_pengajuan.PihcMasterKary = convertSourceTargetDataKaryawan(myKaryawan)
							list_pengajuan.PihcMasterCompany = myCompany
							foto := "https://storage.googleapis.com/lumen-oauth-storage/DataKaryawan/Foto/" + myCompany.Code + "/" + myKaryawan.EmpNo + ".jpg"
							respons, err := http.Get(foto)
							if err != nil || respons.StatusCode != http.StatusOK {
								foto = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
							}
							list_pengajuan.Foto = foto
							list_pengajuan.FotoDefault = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
						}
					}
				}
			}
			// Tipe Absen
			for _, myTipeAbsen := range tipeAbsen {
				if *myCuti.TipeAbsenId == myTipeAbsen.IdTipeAbsen {
					list_pengajuan.TipeAbsen = myTipeAbsen
					result := convertSourceTargetMyPengajuanAbsen(myCuti, myTipeAbsen)
					list_pengajuan.IdPengajuanAbsen = result.IdPengajuanAbsen
					list_pengajuan.MulaiAbsen = result.MulaiAbsen
					list_pengajuan.AkhirAbsen = result.AkhirAbsen
					list_pengajuan.TglPengajuan = result.TglPengajuan
					list_pengajuan.Periode = result.Periode
					if result.Status != nil && *result.Status != "" {
						list_pengajuan.Status = *result.Status
					}
					if result.Deskripsi != nil && *result.Deskripsi != "" {
						list_pengajuan.Deskripsi = *result.Deskripsi
					}
				}
			}
			for _, list_file := range files {
				if myCuti.IdPengajuanAbsen == list_file.PengajuanAbsenId {
					myFiles = append(myFiles, list_file)
				}
			}
			list_pengajuan.FileAbsen = myFiles
			list_aprvl = append(list_aprvl, list_pengajuan)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   list_aprvl,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *TesssController) ShowDetailApprovalPengajuanCuti(ctx *gin.Context) {
	id := ctx.Param("id_pengajuan_absen")
	id_pengajuan, _ := strconv.Atoi(id)

	list_aprvl := Authentication.ListApprovalCuti{}

	dataDB, err := c.PengajuanAbsenRepo.FindDataIdPengajuan(id_pengajuan)

	if err == nil {
		tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*dataDB.TipeAbsenId)
		karyawan, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(dataDB.Nik)
		companys, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(karyawan.Company)
		files, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuan(dataDB.IdPengajuanAbsen)
		if files == nil {
			files = []cuti.FileAbsen{}
		}

		data_karyawan_convert := convertSourceTargetDataKaryawan(karyawan)
		result := convertSourceTargetMyPengajuanAbsen(dataDB, tipeAbsen)

		list_aprvl.IdPengajuanAbsen = result.IdPengajuanAbsen
		list_aprvl.PihcMasterKary = data_karyawan_convert
		list_aprvl.PihcMasterCompany = companys
		list_aprvl.TipeAbsen = tipeAbsen
		list_aprvl.MulaiAbsen = result.MulaiAbsen
		list_aprvl.AkhirAbsen = result.AkhirAbsen
		list_aprvl.TglPengajuan = result.TglPengajuan
		list_aprvl.FileAbsen = files
		list_aprvl.Periode = result.Periode
		foto := "https://storage.googleapis.com/lumen-oauth-storage/DataKaryawan/Foto/" + companys.Code + "/" + result.Nik + ".jpg"
		respons, err := http.Get(foto)
		if err != nil || respons.StatusCode != http.StatusOK {
			foto = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
		}
		list_aprvl.Foto = foto
		list_aprvl.FotoDefault = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"

		if result.Status != nil && *result.Status != "" {
			list_aprvl.Status = *result.Status
		}
		if result.Deskripsi != nil && *result.Deskripsi != "" {
			list_aprvl.Deskripsi = *result.Deskripsi
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   list_aprvl,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

// REPAIRED
// func (c *TesssController) StoreApprovePengajuanAbsen(ctx *gin.Context) {
// 	var req Authentication.ValidationApprovalAtasanPengajuanAbsen

// 	if err := ctx.ShouldBind(&req); err != nil {
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

// 	pengajuan_absen, err := c.PengajuanAbsenRepo.FindDataIdPengajuan(req.IdPengajuanAbsen)

// 	if err == nil {
// 		if req.Status == "Approved" {
// 			if req.Status != "" {
// 				pengajuan_absen.Status = &req.Status
// 			}
// 			if req.Keterangan != "" {
// 				pengajuan_absen.Keterangan = &req.Keterangan
// 			}

// 			c.PengajuanAbsenRepo.Update(pengajuan_absen)

// 			ctx.JSON(http.StatusOK, gin.H{
// 				"status": http.StatusOK,
// 				"info":   "Success",
// 			})
// 		} else {
// 			transaksi, _ := c.TransaksiCutiRepo.FindDataTransaksiIDPengajuan(pengajuan_absen.IdPengajuanAbsen)
// 			tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*pengajuan_absen.TipeAbsenId)
// 			if (tipeAbsen.MaxAbsen == nil) || (*tipeAbsen.MaxAbsen == 0) {
// 				for _, data_transaksi := range transaksi {
// 					saldo_cuti, err_saldo := c.SaldoCutiRepo.GetSaldoCutiPerTipe(pengajuan_absen.Nik, *pengajuan_absen.TipeAbsenId, data_transaksi.Periode)
// 					saldo_cuti_history, _ := c.HistorySaldoCutiRepo.GetSaldoHistory
// 					if err_saldo == nil {
// 						saldo_cuti.Saldo = saldo_cuti.Saldo + data_transaksi.JumlahCuti
// 						c.SaldoCutiRepo.Update(saldo_cuti)
// 					}
// 				}
// 			}

// 			if req.Status != "" {
// 				pengajuan_absen.Status = &req.Status
// 			}
// 			if req.Keterangan != "" {
// 				pengajuan_absen.Keterangan = &req.Keterangan
// 			}

// 			c.PengajuanAbsenRepo.Update(pengajuan_absen)

// 			// c.TransaksiCutiRepo.Delete(pengajuan_absen.IdPengajuanAbsen)

// 			ctx.JSON(http.StatusOK, gin.H{
// 				"status": http.StatusOK,
// 				"info":   "Success",
// 			})
// 		}
// 	} else {
// 		ctx.AbortWithStatus(http.StatusInternalServerError)
// 	}
// }

// SALDO CUTI
func (c *TesssController) StoreAdminSaldoCutiKaryawan(ctx *gin.Context) {
	var req Authentication.ValidasiStoreSaldoCuti
	var data Authentication.SaldoCutiKaryawan
	var sc cuti.SaldoCuti

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

	kebenaran := false
	var keterangan string
	if req.IDSaldo == 0 {
		sc.TipeAbsenId = req.TipeAbsenId
		sc.Nik = req.Nik
		sc.Saldo = req.Saldo
		sc.ValidFrom, _ = time.Parse(time.DateOnly, req.ValidFrom)
		sc.ValidTo, _ = time.Parse(time.DateOnly, req.ValidTo)
		sc.CreatedBy = req.CreatedBy

		periode := strconv.Itoa(time.Now().Year())
		sc.Periode = periode
		sc.MaxHutang = req.MaxHutang
		sc.ValidFromHutang, _ = time.Parse(time.DateOnly, req.ValidFrom)

		saldoCuti, err := c.SaldoCutiRepo.Create(sc)
		if err == nil {
			dataSaldoCuti := Authentication.SaldoCutiKaryawan{
				IdSaldoCuti:     saldoCuti.IdSaldoCuti,
				TipeAbsenId:     saldoCuti.TipeAbsenId,
				Nik:             saldoCuti.Nik,
				Saldo:           saldoCuti.Saldo,
				ValidFrom:       saldoCuti.ValidFrom.Format(time.DateOnly),
				ValidTo:         saldoCuti.ValidTo.Format(time.DateOnly),
				CreatedBy:       saldoCuti.CreatedBy,
				CreatedAt:       saldoCuti.CreatedAt,
				UpdatedAt:       saldoCuti.UpdatedAt,
				Periode:         saldoCuti.Periode,
				MaxHutang:       saldoCuti.MaxHutang,
				ValidFromHutang: saldoCuti.ValidFromHutang.Format(time.DateOnly),
			}
			dataHistorySaldoCuti := cuti.HistorySaldoCuti{
				IdHistorySaldoCuti: saldoCuti.IdSaldoCuti,
				TipeAbsenId:        saldoCuti.TipeAbsenId,
				Nik:                saldoCuti.Nik,
				Saldo:              saldoCuti.Saldo,
				ValidFrom:          saldoCuti.ValidFrom,
				ValidTo:            saldoCuti.ValidTo,
				CreatedBy:          saldoCuti.CreatedBy,
				CreatedAt:          saldoCuti.UpdatedAt,
				Periode:            saldoCuti.Periode,
				MaxHutang:          saldoCuti.MaxHutang,
				ValidFromHutang:    saldoCuti.ValidFromHutang,
			}
			c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)

			data = dataSaldoCuti

			kebenaran = true
			keterangan = "Success"
		} else {
			data = Authentication.SaldoCutiKaryawan{}

			keterangan = "Gagal Store Saldo Cuti"
		}
	} else {
		sc, _ := c.SaldoCutiRepo.GetSaldoCutiByID(req.IDSaldo)
		sc.Saldo = req.Saldo
		sc.ValidFrom, _ = time.Parse(time.DateOnly, req.ValidFrom)
		sc.ValidTo, _ = time.Parse(time.DateOnly, req.ValidTo)
		sc.MaxHutang = req.MaxHutang
		sc.ValidFromHutang, _ = time.Parse(time.DateOnly, req.ValidFrom)

		saldoCuti, err := c.SaldoCutiRepo.Update(sc)
		if err == nil {
			dataSaldoCuti := Authentication.SaldoCutiKaryawan{
				IdSaldoCuti:     saldoCuti.IdSaldoCuti,
				TipeAbsenId:     saldoCuti.TipeAbsenId,
				Nik:             saldoCuti.Nik,
				Saldo:           saldoCuti.Saldo,
				ValidFrom:       saldoCuti.ValidFrom.Format(time.DateOnly),
				ValidTo:         saldoCuti.ValidTo.Format(time.DateOnly),
				CreatedBy:       saldoCuti.CreatedBy,
				CreatedAt:       saldoCuti.CreatedAt,
				UpdatedAt:       saldoCuti.UpdatedAt,
				Periode:         saldoCuti.Periode,
				MaxHutang:       saldoCuti.MaxHutang,
				ValidFromHutang: saldoCuti.ValidFromHutang.Format(time.DateOnly),
			}
			dataHistorySaldoCuti := cuti.HistorySaldoCuti{
				IdHistorySaldoCuti: saldoCuti.IdSaldoCuti,
				TipeAbsenId:        saldoCuti.TipeAbsenId,
				Nik:                saldoCuti.Nik,
				Saldo:              saldoCuti.Saldo,
				ValidFrom:          saldoCuti.ValidFrom,
				ValidTo:            saldoCuti.ValidTo,
				CreatedBy:          saldoCuti.CreatedBy,
				CreatedAt:          saldoCuti.UpdatedAt,
				Periode:            saldoCuti.Periode,
				MaxHutang:          saldoCuti.MaxHutang,
				ValidFromHutang:    saldoCuti.ValidFromHutang,
			}
			c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)

			data = dataSaldoCuti

			kebenaran = true
			keterangan = "Success"
		} else {
			data = Authentication.SaldoCutiKaryawan{}

			keterangan = "Gagal Update Saldo Cuti"
		}
	}

	if kebenaran {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"Success": keterangan,
			"data":    data,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"status":  http.StatusServiceUnavailable,
			"Success": keterangan,
			"data":    data,
		})
	}
}

func (c *TesssController) ListAdminSaldoCutiKaryawan(ctx *gin.Context) {
	var req Authentication.ValidasiListSaldoCuti
	data := []Authentication.ListSaldoCutiKaryawan{}

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

	saldoCuti, err := c.SaldoCutiRepo.FindSaldoCutiKaryawanAdmin(req.Nik, req.Tahun)

	if err == nil {
		for _, dataSaldoo := range saldoCuti {
			karyawan, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(dataSaldoo.Nik)
			company, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(karyawan.Company)
			TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(dataSaldoo.TipeAbsenId)
			TipeAbsenKaryawan := Authentication.GetTipeAbsenKaryawanSaldo{
				IdTipeAbsen:   TipeAbsen.IdTipeAbsen,
				NamaTipeAbsen: TipeAbsen.NamaTipeAbsen,
				TipeMaxAbsen:  TipeAbsen.TipeMaxAbsen,
				CreatedAt:     TipeAbsen.CreatedAt,
				UpdatedAt:     TipeAbsen.UpdatedAt,
			}
			if TipeAbsen.CompCode != nil && *TipeAbsen.CompCode != "" {
				TipeAbsenKaryawan.CompCode = *TipeAbsen.CompCode
			}
			if TipeAbsen.MaxAbsen != nil && *TipeAbsen.MaxAbsen != 0 {
				TipeAbsenKaryawan.MaxAbsen = *TipeAbsen.MaxAbsen
			}
			CompanyKaryawans := Authentication.CompanyKaryawan{
				Code: company.Code,
				Name: company.Name,
			}
			dataSaldoCuti := Authentication.ListSaldoCutiKaryawan{
				IdSaldoCuti:               dataSaldoo.IdSaldoCuti,
				GetTipeAbsenKaryawanSaldo: TipeAbsenKaryawan,
				Nik:                       dataSaldoo.Nik,
				CompanyKaryawan:           CompanyKaryawans,
				Saldo:                     dataSaldoo.Saldo,
				ValidFrom:                 dataSaldoo.ValidFrom.Format(time.DateOnly),
				ValidTo:                   dataSaldoo.ValidTo.Format(time.DateOnly),
				CreatedBy:                 dataSaldoo.CreatedBy,
				CreatedAt:                 dataSaldoo.CreatedAt,
				UpdatedAt:                 dataSaldoo.UpdatedAt,
				Periode:                   dataSaldoo.Periode,
				MaxHutang:                 dataSaldoo.MaxHutang,
				ValidFromHutang:           dataSaldoo.ValidFromHutang.Format(time.DateOnly),
			}
			if karyawan.Nama != nil && *karyawan.Nama != "" {
				dataSaldoCuti.Nama = *karyawan.Nama
			}
			data = append(data, dataSaldoCuti)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
}

func (c *TesssController) GetAdminSaldoCuti(ctx *gin.Context) {
	id := ctx.Param("id_saldo_cuti")
	var data Authentication.ListSaldoCutiKaryawan

	id_saldo, _ := strconv.Atoi(id)
	saldoCuti, err := c.SaldoCutiRepo.GetSaldoCutiByID(id_saldo)
	if err == nil {
		karyawan, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(saldoCuti.Nik)
		company, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(karyawan.Company)
		TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(saldoCuti.TipeAbsenId)
		TipeAbsenKaryawan := Authentication.GetTipeAbsenKaryawanSaldo{
			IdTipeAbsen:   TipeAbsen.IdTipeAbsen,
			NamaTipeAbsen: TipeAbsen.NamaTipeAbsen,
			TipeMaxAbsen:  TipeAbsen.TipeMaxAbsen,
			CreatedAt:     TipeAbsen.CreatedAt,
			UpdatedAt:     TipeAbsen.UpdatedAt,
		}
		if TipeAbsen.CompCode != nil && *TipeAbsen.CompCode != "" {
			TipeAbsenKaryawan.CompCode = *TipeAbsen.CompCode
		}
		if TipeAbsen.MaxAbsen != nil && *TipeAbsen.MaxAbsen != 0 {
			TipeAbsenKaryawan.MaxAbsen = *TipeAbsen.MaxAbsen
		}
		CompanyKaryawans := Authentication.CompanyKaryawan{
			Code: company.Code,
			Name: company.Name,
		}
		dataSaldoCuti := Authentication.ListSaldoCutiKaryawan{
			IdSaldoCuti:               saldoCuti.IdSaldoCuti,
			GetTipeAbsenKaryawanSaldo: TipeAbsenKaryawan,
			Nik:                       saldoCuti.Nik,
			CompanyKaryawan:           CompanyKaryawans,
			Saldo:                     saldoCuti.Saldo,
			ValidFrom:                 saldoCuti.ValidFrom.Format(time.DateOnly),
			ValidTo:                   saldoCuti.ValidTo.Format(time.DateOnly),
			CreatedBy:                 saldoCuti.CreatedBy,
			CreatedAt:                 saldoCuti.CreatedAt,
			UpdatedAt:                 saldoCuti.UpdatedAt,
			Periode:                   saldoCuti.Periode,
			MaxHutang:                 saldoCuti.MaxHutang,
			ValidFromHutang:           saldoCuti.ValidFromHutang.Format(time.DateOnly),
		}
		if karyawan.Nama != nil && *karyawan.Nama != "" {
			dataSaldoCuti.Nama = *karyawan.Nama
		}
		data = dataSaldoCuti

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   data})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Data Tidak Ada",
			"data":   nil,
		})
	}
}

func (c *TesssController) DeleteAdminSaldoCuti(ctx *gin.Context) {
	id := ctx.Param("id_saldo_cuti")
	var data Authentication.SaldoCutiKaryawan

	id_saldo, _ := strconv.Atoi(id)
	saldoCuti, err := c.SaldoCutiRepo.DelAdminSaldoCuti(id_saldo)

	if err == nil {
		dataSaldoCuti := Authentication.SaldoCutiKaryawan{
			IdSaldoCuti:     saldoCuti.IdSaldoCuti,
			TipeAbsenId:     saldoCuti.TipeAbsenId,
			Nik:             saldoCuti.Nik,
			Saldo:           saldoCuti.Saldo,
			ValidFrom:       saldoCuti.ValidFrom.Format(time.DateOnly),
			ValidTo:         saldoCuti.ValidTo.Format(time.DateOnly),
			CreatedBy:       saldoCuti.CreatedBy,
			CreatedAt:       saldoCuti.CreatedAt,
			UpdatedAt:       saldoCuti.UpdatedAt,
			Periode:         saldoCuti.Periode,
			MaxHutang:       saldoCuti.MaxHutang,
			ValidFromHutang: saldoCuti.ValidFromHutang.Format(time.DateOnly),
		}
		dataHistorySaldoCuti := cuti.HistorySaldoCuti{
			IdHistorySaldoCuti: saldoCuti.IdSaldoCuti,
			TipeAbsenId:        saldoCuti.TipeAbsenId,
			Nik:                saldoCuti.Nik,
			Saldo:              saldoCuti.Saldo,
			ValidFrom:          saldoCuti.ValidFrom,
			ValidTo:            saldoCuti.ValidTo,
			CreatedBy:          saldoCuti.CreatedBy,
			CreatedAt:          saldoCuti.UpdatedAt,
			Periode:            saldoCuti.Periode,
			MaxHutang:          saldoCuti.MaxHutang,
			ValidFromHutang:    saldoCuti.ValidFromHutang,
		}
		c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)
		data = dataSaldoCuti

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

// Tipe Saldo Cuti
// func (c *TesssController) GetTipeAbsenSaldoPengajuan(ctx *gin.Context) {
// 	var req Authentication.ValidationNIKTahun

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
// 	data := []Authentication.GetTipeAbsenSaldoIndiv{}
// 	data2 := []Authentication.GetTipeAbsenSaldoIndiv{}

// 	pihc_mstr_krywn, err_pihc_mstr_krywn := c.PihcMasterKaryDbRepo.FindUserByNIK(req.NIK)

// 	if err_pihc_mstr_krywn == nil {
// 		fmt.Println(pihc_mstr_krywn.Company)
// 		TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenPengajuan(pihc_mstr_krywn.Company)

// 		for _, dataCuti := range TipeAbsen {
// 			saldoCutiPerTipe, err := c.SaldoCutiRepo.GetSaldoCutiPerTipe(pihc_mstr_krywn.EmpNo, dataCuti.IdTipeAbsen, req.Tahun)
// 			if err == nil {
// 				if dataCuti.NamaTipeAbsen == "Cuti Tahunan" {
// 					max_absen := &Authentication.MaxAbsenIndiv{
// 						TipeMaxAbsen: dataCuti.TipeMaxAbsen,
// 					}
// 					saldo := &Authentication.SaldoIndiv{
// 						Saldo:           saldoCutiPerTipe.Saldo,
// 						ValidFrom:       saldoCutiPerTipe.ValidFrom.Format(time.DateOnly),
// 						ValidTo:         saldoCutiPerTipe.ValidTo.Format(time.DateOnly),
// 						Periode:         saldoCutiPerTipe.Periode,
// 						MaxHutang:       saldoCutiPerTipe.MaxHutang,
// 						ValidFromHutang: saldoCutiPerTipe.ValidFromHutang.Format(time.DateOnly),
// 					}
// 					tipeSaldoCuti := Authentication.GetTipeAbsenSaldoIndiv{
// 						IdTipeAbsen:   dataCuti.IdTipeAbsen,
// 						NamaTipeAbsen: dataCuti.NamaTipeAbsen,
// 					}
// 					if dataCuti.MaxAbsen != nil && *dataCuti.MaxAbsen != 0 {
// 						max_absen.MaxAbsen = *dataCuti.MaxAbsen
// 						tipeSaldoCuti.MaxAbsenIndiv = max_absen
// 					}
// 					if saldo.Periode != "" {
// 						tipeSaldoCuti.SaldoIndiv = saldo
// 					}

// 					data = append(data, tipeSaldoCuti)
// 				} else {
// 					max_absen := &Authentication.MaxAbsenIndiv{
// 						TipeMaxAbsen: dataCuti.TipeMaxAbsen,
// 					}
// 					saldo := &Authentication.SaldoIndiv{
// 						Saldo:           saldoCutiPerTipe.Saldo,
// 						ValidFrom:       saldoCutiPerTipe.ValidFrom.Format(time.DateOnly),
// 						ValidTo:         saldoCutiPerTipe.ValidTo.Format(time.DateOnly),
// 						Periode:         saldoCutiPerTipe.Periode,
// 						MaxHutang:       saldoCutiPerTipe.MaxHutang,
// 						ValidFromHutang: saldoCutiPerTipe.ValidFromHutang.Format(time.DateOnly),
// 					}
// 					tipeSaldoCuti := Authentication.GetTipeAbsenSaldoIndiv{
// 						IdTipeAbsen:   dataCuti.IdTipeAbsen,
// 						NamaTipeAbsen: dataCuti.NamaTipeAbsen,
// 					}
// 					if saldo.Periode != "" {
// 						tipeSaldoCuti.SaldoIndiv = saldo
// 					}
// 					if dataCuti.MaxAbsen != nil && *dataCuti.MaxAbsen != 0 {
// 						max_absen.MaxAbsen = *dataCuti.MaxAbsen
// 						tipeSaldoCuti.MaxAbsenIndiv = max_absen
// 					}
// 					data2 = append(data2, tipeSaldoCuti)
// 				}
// 			} else {
// 				max_absen := &Authentication.MaxAbsenIndiv{
// 					TipeMaxAbsen: dataCuti.TipeMaxAbsen,
// 				}
// 				saldo := &Authentication.SaldoIndiv{
// 					Saldo:           saldoCutiPerTipe.Saldo,
// 					ValidFrom:       saldoCutiPerTipe.ValidFrom.Format(time.DateOnly),
// 					ValidTo:         saldoCutiPerTipe.ValidTo.Format(time.DateOnly),
// 					Periode:         saldoCutiPerTipe.Periode,
// 					MaxHutang:       saldoCutiPerTipe.MaxHutang,
// 					ValidFromHutang: saldoCutiPerTipe.ValidFromHutang.Format(time.DateOnly),
// 				}
// 				tipeSaldoCuti := Authentication.GetTipeAbsenSaldoIndiv{
// 					IdTipeAbsen:   dataCuti.IdTipeAbsen,
// 					NamaTipeAbsen: dataCuti.NamaTipeAbsen,
// 				}
// 				if saldo.Periode != "" {
// 					tipeSaldoCuti.SaldoIndiv = saldo
// 				}
// 				if dataCuti.MaxAbsen != nil && *dataCuti.MaxAbsen != 0 {
// 					max_absen.MaxAbsen = *dataCuti.MaxAbsen
// 					tipeSaldoCuti.MaxAbsenIndiv = max_absen
// 				}

// 				data2 = append(data2, tipeSaldoCuti)
// 			}
// 		}

// 		data = append(data, data2...)
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"status":  http.StatusOK,
// 		"success": "Success",
// 		"data":    data,
// 	})
// }