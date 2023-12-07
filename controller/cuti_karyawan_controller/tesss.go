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
			approvedBy := "82105096"
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
