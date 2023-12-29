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

type CutiKrywnController struct {
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

func NewCutiKrywnController(Db *gorm.DB) *CutiKrywnController {
	return &CutiKrywnController{PengajuanAbsenRepo: cuti.NewPengajuanAbsenRepo(Db),
		HistoryPengajuanAbsenRepo: cuti.NewHistoryPengajuanAbsenRepo(Db),
		SaldoCutiRepo:             cuti.NewSaldoCutiRepo(Db),
		HistorySaldoCutiRepo:      cuti.NewHistorySaldoCutiRepo(Db),
		TipeAbsenRepo:             cuti.NewTipeAbsenRepo(Db),
		FileAbsenRepo:             cuti.NewFileAbsenRepo(Db),
		TransaksiCutiRepo:         cuti.NewTransaksiCutiRepo(Db),
		PihcMasterKaryDbRepo:      pihc.NewPihcMasterKaryDbRepo(Db),
		PihcMasterCompanyRepo:     pihc.NewPihcMasterCompanyRepo(Db)}
}

// Pengajuan Cuti (DONE)
func (c *CutiKrywnController) StoreCutiKaryawan(ctx *gin.Context) {
	var req Authentication.ValidasiStoreCutiKaryawan
	var sck cuti.PengajuanAbsen
	var fsc []cuti.FileAbsen
	var trsc []Authentication.SaldoCutiTransaksiPengajuan

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

	if req.IdPengajuanAbsen != nil {
		req.IdPengajuanAbsen = ConvertInterfaceTypeDataToInt(req.IdPengajuanAbsen)
	}

	PIHC_MSTR_KRY, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(req.Nik)
	comp_code := PIHC_MSTR_KRY.Company

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

	// Mencari Atasan
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

	tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*sck.TipeAbsenId)

	saldoPeriode, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipeArr(sck.Nik, tipeAbsen.IdTipeAbsen, *sck.Periode)
	result := perhitungan(sck.MulaiAbsen, sck.AkhirAbsen)
	sck.JmlHariKalendar = &result[0]
	sck.JmlHariKerja = &result[1]
	saldo_digunakan := 0
	tipe_hari := ""
	if *tipeAbsen.TipeMaxAbsen == "hari_kalender" {
		saldo_digunakan = result[0]
		tipe_hari = "hari_kalender"
	} else if *tipeAbsen.TipeMaxAbsen == "hari_kerja" {
		saldo_digunakan = result[1]
		tipe_hari = "hari_kerja"
	}
	saldo_terpakai := 0
	var newPeriode time.Time
	keterangan := ""
	keterangan_x := ""
	checkSaldo := false
	create := false
	update := false
	isSaldo := false
	fmt.Println(req.IdPengajuanAbsen)

	if req.IdPengajuanAbsen == nil {
		if tipeAbsen.MaxAbsen != nil {
			JmlHariKerja := 0
			jmlhHariKalender := 0

			transaksi_cuti := cuti.TransaksiCuti{}
			if *tipeAbsen.TipeMaxAbsen == "hari_kalender" {
				for currentDate := sck.MulaiAbsen; jmlhHariKalender != *tipeAbsen.MaxAbsen; currentDate = currentDate.AddDate(0, 0, 1) {
					jmlhHariKalender++
					if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
						JmlHariKerja++
					}
				}
				sck.AkhirAbsen = sck.MulaiAbsen.AddDate(0, 0, jmlhHariKalender-1)
				// Transaksi Cuti
				transaksi_cuti.JumlahCuti = jmlhHariKalender
			} else if *tipeAbsen.TipeMaxAbsen == "hari_kerja" {
				for currentDate := sck.MulaiAbsen; JmlHariKerja != *tipeAbsen.MaxAbsen; currentDate = currentDate.AddDate(0, 0, 1) {
					jmlhHariKalender++
					if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
						JmlHariKerja++
					}
				}
				sck.AkhirAbsen = sck.MulaiAbsen.AddDate(0, 0, jmlhHariKalender-1)
				// Transaksi Cuti
				transaksi_cuti.JumlahCuti = JmlHariKerja
			}
			transaksi_cuti.TipeHari = tipe_hari

			sck.JmlHariKalendar = &jmlhHariKalender
			sck.JmlHariKerja = &JmlHariKerja

			sckData, _ := c.PengajuanAbsenRepo.Create(sck)
			// CREATE HistoryPengajuanAbsen
			history_pengajuan_absen := HistoryPengajuanCutiSet(sckData)
			c.HistoryPengajuanAbsenRepo.Create(history_pengajuan_absen)

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
			isHutang := false
			saldoHutang := 0
			for _, saldo := range saldoPeriode {
				var result [2]int
				jmlahCuti := 0

				if (sck.MulaiAbsen.Before(saldo.ValidTo) || sck.MulaiAbsen.Equal(saldo.ValidTo)) &&
					(sck.MulaiAbsen.After(saldo.ValidFrom) || sck.MulaiAbsen.Equal(saldo.ValidFrom)) &&
					(sck.AkhirAbsen.After(saldo.ValidTo) || sck.AkhirAbsen.Equal(saldo.ValidTo)) {
					// MulaiAbsen <= ValidTo && MulaiAbsen >= ValidFrom && AkhirAbsen>=ValidTo
					if saldo_digunakan <= saldo.Saldo {
						result = perhitungan(sck.MulaiAbsen, saldo.ValidTo)
						newPeriode = saldo.ValidTo
						isSaldo = true
						isHutang = false
					} else if sck.MulaiAbsen.After(saldo.ValidFromHutang) { // MulaiAbsen >= ValidFromHutang
						if saldo_digunakan <= saldo.Saldo+saldo.MaxHutang {
							result = perhitungan(sck.MulaiAbsen, saldo.ValidTo)
							newPeriode = saldo.ValidTo
							isSaldo = true
							isHutang = true
						} else {
							keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
								"</b> dan Maksimal Berhutang " + strconv.Itoa(saldo.MaxHutang)
							isHutang = false
							isSaldo = false
						}
					} else { // Diluar Masa Berlaku Hutang
						keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
							"</b> dan Anda Berada diluar Masa Berlaku Hutang, Masa Berlaku hutang dimulai pada " +
							saldo.ValidFromHutang.Format(time.DateOnly)
						isHutang = false
						isSaldo = false
					}
				} else if (newPeriode.After(saldo.ValidFrom) || newPeriode.Equal(saldo.ValidFrom)) &&
					(sck.AkhirAbsen.Before(saldo.ValidTo) || sck.AkhirAbsen.Equal(saldo.ValidTo)) {
					// newPeriode>=ValidFrom && AkhirAbsen<=ValidTo (periode ke-2)
					if saldo_digunakan <= saldo.Saldo-saldoHutang {
						result = perhitungan(newPeriode.AddDate(0, 0, 1), sck.AkhirAbsen)
						isHutang = false
						isSaldo = true
					} else if sck.MulaiAbsen.After(saldo.ValidFromHutang) {
						if saldo_digunakan <= saldo.Saldo-saldoHutang+saldo.MaxHutang {
							result = perhitungan(newPeriode.AddDate(0, 0, 1), sck.AkhirAbsen)
							isHutang = true
							isSaldo = true
						} else {
							keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
								"</b> dan Maksimal Berhutang " + strconv.Itoa(saldo.MaxHutang)
							isHutang = false
							isSaldo = false
						}
					} else {
						keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
							"</b> dan Anda Berada diluar Masa Berlaku Hutang, Masa Berlaku hutang dimulai pada " +
							saldo.ValidFromHutang.Format(time.DateOnly)
						isHutang = false
						isSaldo = false
					}
				} else if (sck.MulaiAbsen.After(saldo.ValidFrom) || sck.MulaiAbsen.Equal(saldo.ValidFrom)) &&
					(sck.AkhirAbsen.Before(saldo.ValidTo) || sck.AkhirAbsen.Equal(saldo.ValidTo)) {
					// MulaiAbsen >= ValidFrom && AkhirAbsen <= ValidTo
					if saldo_digunakan <= saldo.Saldo {
						result = perhitungan(sck.MulaiAbsen, sck.AkhirAbsen)
						isHutang = false
						isSaldo = true
					} else if sck.MulaiAbsen.After(saldo.ValidFromHutang) {
						if saldo_digunakan <= saldo.Saldo+saldo.MaxHutang {
							result = perhitungan(sck.MulaiAbsen, sck.AkhirAbsen)
							isHutang = true
							isSaldo = true
						} else {
							keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
								"</b> dan Maksimal Berhutang " + strconv.Itoa(saldo.MaxHutang)
							isHutang = false
							isSaldo = false
						}
					} else {
						keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
							"</b> dan Anda Berada diluar Masa Berlaku Hutang, Masa Berlaku hutang dimulai pada " +
							saldo.ValidFromHutang.Format(time.DateOnly)
						isHutang = false
						isSaldo = false
					}

				}

				if *tipeAbsen.TipeMaxAbsen == "hari_kerja" {
					jmlahCuti = result[1]
					saldo_terpakai += jmlahCuti
				} else if *tipeAbsen.TipeMaxAbsen == "hari_kalender" {
					jmlahCuti = result[0]
					saldo_terpakai += jmlahCuti
				}

				if isSaldo {
					source := Authentication.SaldoCutiTransaksiPengajuan{
						SaldoCuti: saldo,
						JmlhCuti:  jmlahCuti,
					}
					trsc = append(trsc, source)
				}

				if isHutang {
					saldoHutang = saldo_terpakai - saldo.Saldo
				}
			}

			if saldo_digunakan == saldo_terpakai {
				checkSaldo = true
				create = true
			} else {
				keterangan_x = "Maaf Saldo Anda Tidak Cukup" + keterangan
			}
		}
	} else {
		pengajuan_absen, _ := c.PengajuanAbsenRepo.FindDataIdPengajuan(req.IdPengajuanAbsen)
		if *pengajuan_absen.Status == "Rejected" || *pengajuan_absen.Status == "WaitApproved" {
			tipe_absen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(req.TipeAbsenId)
			file_absen, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuan(pengajuan_absen.IdPengajuanAbsen)
			transaksi_absen, _ := c.TransaksiCutiRepo.FindDataTransaksiIDPengajuan(pengajuan_absen.IdPengajuanAbsen)

			sck.IdPengajuanAbsen = pengajuan_absen.IdPengajuanAbsen
			sck.TipeAbsenId = &tipe_absen.IdTipeAbsen
			sck.CreatedAt = pengajuan_absen.CreatedAt

			for _, delete_file := range file_absen {
				c.FileAbsenRepo.Delete(delete_file)
			}
			for _, delete_transaksi := range transaksi_absen {
				c.TransaksiCutiRepo.Delete(delete_transaksi)
			}
			for _, fa := range req.FileAbsen {
				files := cuti.FileAbsen{
					PengajuanAbsenId: pengajuan_absen.IdPengajuanAbsen,
					Filename:         fa.Filename,
					Url:              fa.URL,
					Extension:        fa.Extension,
				}
				fsc = append(fsc, files)
			}

			if tipeAbsen.MaxAbsen != nil {
				JmlHariKerja := 0
				jmlhHariKalender := 0

				transaksi_cuti := cuti.TransaksiCuti{}
				if *tipeAbsen.TipeMaxAbsen == "hari_kalender" {
					for currentDate := sck.MulaiAbsen; jmlhHariKalender != *tipeAbsen.MaxAbsen; currentDate = currentDate.AddDate(0, 0, 1) {
						jmlhHariKalender++
						if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
							JmlHariKerja++
						}
					}
					sck.AkhirAbsen = sck.MulaiAbsen.AddDate(0, 0, jmlhHariKalender-1)
					// Transaksi Cuti
					transaksi_cuti.TipeHari = "hari_kalender"
					transaksi_cuti.JumlahCuti = jmlhHariKalender
				} else if *tipeAbsen.TipeMaxAbsen == "hari_kerja" {
					for currentDate := sck.MulaiAbsen; JmlHariKerja != *tipeAbsen.MaxAbsen; currentDate = currentDate.AddDate(0, 0, 1) {
						jmlhHariKalender++
						if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
							JmlHariKerja++
						}
					}
					sck.AkhirAbsen = sck.MulaiAbsen.AddDate(0, 0, jmlhHariKalender-1)
					// Transaksi Cuti
					transaksi_cuti.TipeHari = "hari_kerja"
					transaksi_cuti.JumlahCuti = JmlHariKerja
				}
				sck.JmlHariKalendar = &jmlhHariKalender
				sck.JmlHariKerja = &JmlHariKerja

				sckData, _ := c.PengajuanAbsenRepo.Update(sck)
				// CREATE HistoryPengajuanAbsen
				history_pengajuan_absen := HistoryPengajuanCutiSet(sckData)
				c.HistoryPengajuanAbsenRepo.Create(history_pengajuan_absen)

				convert := convertSourceTargetMyPengajuanAbsen(sckData, tipeAbsen)

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
				isHutang := false
				saldoHutang := 0
				for _, saldo := range saldoPeriode {
					var result [2]int
					jmlahCuti := 0

					if (sck.MulaiAbsen.Before(saldo.ValidTo) || sck.MulaiAbsen.Equal(saldo.ValidTo)) &&
						(sck.MulaiAbsen.After(saldo.ValidFrom) || sck.MulaiAbsen.Equal(saldo.ValidFrom)) &&
						(sck.AkhirAbsen.After(saldo.ValidTo) || sck.AkhirAbsen.Equal(saldo.ValidTo)) {
						// MulaiAbsen <= ValidTo && MulaiAbsen >= ValidFrom && AkhirAbsen>=ValidTo
						if saldo_digunakan <= saldo.Saldo {
							result = perhitungan(sck.MulaiAbsen, saldo.ValidTo)
							newPeriode = saldo.ValidTo
							isSaldo = true
							isHutang = false
						} else if sck.MulaiAbsen.After(saldo.ValidFromHutang) { // MulaiAbsen >= ValidFromHutang
							if saldo_digunakan <= saldo.Saldo+saldo.MaxHutang {
								result = perhitungan(sck.MulaiAbsen, saldo.ValidTo)
								newPeriode = saldo.ValidTo
								isSaldo = true
								isHutang = true
							} else {
								keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
									"</b> dan Maksimal Berhutang " + strconv.Itoa(saldo.MaxHutang)
								isHutang = false
								isSaldo = false
							}
						} else { // Diluar Masa Berlaku Hutang
							keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
								"</b> dan Anda Berada diluar Masa Berlaku Hutang, Masa Berlaku hutang dimulai pada " +
								saldo.ValidFromHutang.Format(time.DateOnly)
							isHutang = false
							isSaldo = false
						}
					} else if (newPeriode.After(saldo.ValidFrom) || newPeriode.Equal(saldo.ValidFrom)) &&
						(sck.AkhirAbsen.Before(saldo.ValidTo) || sck.AkhirAbsen.Equal(saldo.ValidTo)) {
						// newPeriode>=ValidFrom && AkhirAbsen<=ValidTo (periode ke-2)
						if saldo_digunakan <= saldo.Saldo-saldoHutang {
							result = perhitungan(newPeriode.AddDate(0, 0, 1), sck.AkhirAbsen)
							isHutang = false
							isSaldo = true
						} else if sck.MulaiAbsen.After(saldo.ValidFromHutang) {
							if saldo_digunakan <= saldo.Saldo-saldoHutang+saldo.MaxHutang {
								result = perhitungan(newPeriode.AddDate(0, 0, 1), sck.AkhirAbsen)
								isHutang = true
								isSaldo = true
							} else {
								keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
									"</b> dan Maksimal Berhutang " + strconv.Itoa(saldo.MaxHutang)
								isHutang = false
								isSaldo = false
							}
						} else {
							keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
								"</b> dan Anda Berada diluar Masa Berlaku Hutang, Masa Berlaku hutang dimulai pada " +
								saldo.ValidFromHutang.Format(time.DateOnly)
							isHutang = false
							isSaldo = false
						}
					} else if (sck.MulaiAbsen.After(saldo.ValidFrom) || sck.MulaiAbsen.Equal(saldo.ValidFrom)) &&
						(sck.AkhirAbsen.Before(saldo.ValidTo) || sck.AkhirAbsen.Equal(saldo.ValidTo)) {
						// MulaiAbsen >= ValidFrom && AkhirAbsen <= ValidTo
						if saldo_digunakan <= saldo.Saldo {
							result = perhitungan(sck.MulaiAbsen, sck.AkhirAbsen)
							isHutang = false
							isSaldo = true
						} else if sck.MulaiAbsen.After(saldo.ValidFromHutang) {
							if saldo_digunakan <= saldo.Saldo+saldo.MaxHutang {
								result = perhitungan(sck.MulaiAbsen, sck.AkhirAbsen)
								isHutang = true
								isSaldo = true
							} else {
								keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
									"</b> dan Maksimal Berhutang " + strconv.Itoa(saldo.MaxHutang)
								isHutang = false
								isSaldo = false
							}
						} else {
							keterangan = ", Anda memiliki <b>Saldo: " + strconv.Itoa(saldo.Saldo) +
								"</b> dan Anda Berada diluar Masa Berlaku Hutang, Masa Berlaku hutang dimulai pada " +
								saldo.ValidFromHutang.Format(time.DateOnly)
							isHutang = false
							isSaldo = false
						}

					}

					if *tipeAbsen.TipeMaxAbsen == "hari_kerja" {
						jmlahCuti = result[1]
						saldo_terpakai += jmlahCuti
					} else if *tipeAbsen.TipeMaxAbsen == "hari_kalender" {
						jmlahCuti = result[0]
						saldo_terpakai += jmlahCuti
					}

					if isSaldo {
						source := Authentication.SaldoCutiTransaksiPengajuan{
							SaldoCuti: saldo,
							JmlhCuti:  jmlahCuti,
						}
						trsc = append(trsc, source)
					}

					if isHutang {
						saldoHutang = saldo_terpakai - saldo.Saldo
					}
				}

				if saldo_digunakan == saldo_terpakai {
					checkSaldo = true
					update = true
				} else {
					keterangan_x = "Maaf Saldo Anda Tidak Cukup" + keterangan
				}
			}
		}
	}

	if checkSaldo {
		if create {
			sckData, _ := c.PengajuanAbsenRepo.Create(sck)
			history_pengajuan_absen := HistoryPengajuanCutiSet(sckData)
			c.HistoryPengajuanAbsenRepo.Create(history_pengajuan_absen)

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

			// // Transaksi Cuti
			transaksi_cuti := cuti.TransaksiCuti{}

			if sckData.Periode != nil {
				transaksi_cuti.Periode = *sckData.Periode
			}
			for _, transaction := range trsc {
				transaksi_cuti := cuti.TransaksiCuti{}
				transaksi_cuti.PengajuanAbsenId = sckData.IdPengajuanAbsen
				transaksi_cuti.Nik = sckData.Nik
				transaksi_cuti.Periode = transaction.Periode
				transaksi_cuti.JumlahCuti = transaction.JmlhCuti
				transaksi_cuti.TipeHari = tipe_hari
				c.TransaksiCutiRepo.Create(transaksi_cuti)
			}

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
		} else if !create {
			if update {
				sckData, _ := c.PengajuanAbsenRepo.Update(sck)
				history_pengajuan_absen := HistoryPengajuanCutiSet(sckData)
				c.HistoryPengajuanAbsenRepo.Create(history_pengajuan_absen)

				convert := convertSourceTargetMyPengajuanAbsen(sckData, tipeAbsen)

				// CREATE FileAbsen
				files, _ := c.FileAbsenRepo.CreateArr(fsc)

				// // Transaksi Cuti
				transaksi_cuti := cuti.TransaksiCuti{}

				if sckData.Periode != nil {
					transaksi_cuti.Periode = *sckData.Periode
				}
				for _, transaction := range trsc {
					transaksi_cuti := cuti.TransaksiCuti{}
					transaksi_cuti.PengajuanAbsenId = sckData.IdPengajuanAbsen
					transaksi_cuti.Nik = sckData.Nik
					transaksi_cuti.Periode = transaction.Periode
					transaksi_cuti.JumlahCuti = transaction.JmlhCuti
					transaksi_cuti.TipeHari = tipe_hari
					c.TransaksiCutiRepo.Create(transaksi_cuti)
				}

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
			} else if !update {
				ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"status":     http.StatusServiceUnavailable,
					"keterangan": "Gagal Mengubah Pengajuan Absen",
				})
			} else {
				ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"status":     http.StatusServiceUnavailable,
					"keterangan": "Gagal Membuat Pengajuan Absen",
				})
			}

		}
	} else {
		if keterangan_x != "" {
			fmt.Println("A")
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"status":     http.StatusServiceUnavailable,
				"keterangan": keterangan_x,
			})
		} else {
			fmt.Println("B")
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
func (c *CutiKrywnController) ShowDetailPengajuanCuti(ctx *gin.Context) {
	data := Authentication.PengajuanAbsens{}
	id := ctx.Param("id_pengajuan_absen")
	id_pengajuan, _ := strconv.Atoi(id)
	fmt.Println(id_pengajuan)

	data_pengajuan, err_pengajuan := c.PengajuanAbsenRepo.FindDataIdPengajuan(id_pengajuan)

	if err_pengajuan == nil {
		data_tipe_absen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*data_pengajuan.TipeAbsenId)
		data_file_absen, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuan(data_pengajuan.IdPengajuanAbsen)

		convert := convertSourceTargetMyPengajuanAbsen(data_pengajuan, data_tipe_absen)

		if data_file_absen == nil {
			data_file_absen = []cuti.FileAbsen{}
		}
		data.MyPengajuanAbsen = convert
		data.File = data_file_absen

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"data":   data,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
func (c *CutiKrywnController) GetMyPengajuanCuti(ctx *gin.Context) {
	var req Authentication.ValidationNIKTahun
	var data []Authentication.PengajuanAbsens

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
	dataDB, err := c.PengajuanAbsenRepo.FindDataNIKPeriode(req.NIK, req.Tahun)

	if err == nil {
		for _, myCuti := range dataDB {
			data_tipe_absen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*myCuti.TipeAbsenId)
			data_file_absen, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuan(myCuti.IdPengajuanAbsen)
			convert := convertSourceTargetMyPengajuanAbsen(myCuti, data_tipe_absen)

			if data_file_absen == nil {
				data_file_absen = []cuti.FileAbsen{}
			}
			list := Authentication.PengajuanAbsens{
				MyPengajuanAbsen: convert,
				File:             data_file_absen,
			}
			data = append(data, list)
		}

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
func (c *CutiKrywnController) DeletePengajuanCuti(ctx *gin.Context) {
	id := ctx.Param("id_pengajuan_absen")
	id_pengajuan_absen, _ := strconv.Atoi(id)

	pengajuanAbsen, err := c.PengajuanAbsenRepo.FindDataIdPengajuan(id_pengajuan_absen)

	if err == nil {
		if *pengajuanAbsen.Status == "WaitApproved" || *pengajuanAbsen.Status == "Rejected" {
			c.PengajuanAbsenRepo.DelPengajuanCuti(pengajuanAbsen.IdPengajuanAbsen)
			transaksi_cuti, _ := c.TransaksiCutiRepo.FindDataTransaksiIDPengajuan(pengajuanAbsen.IdPengajuanAbsen)

			for _, tr := range transaksi_cuti {
				c.TransaksiCutiRepo.Delete(tr)
			}
			ctx.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"info":   "success",
			})
		} else {
			ctx.AbortWithStatus(http.StatusServiceUnavailable)
		}
	} else {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// Approval Cuti (DONE)
func (c *CutiKrywnController) StoreApprovePengajuanAbsen(ctx *gin.Context) {
	var req Authentication.ValidationApprovalAtasanPengajuanAbsen
	var eksekusi_saldo []cuti.SaldoCuti
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

	pengajuan_absen, err := c.PengajuanAbsenRepo.FindDataIdPengajuan(req.IdPengajuanAbsen)
	if err == nil {
		if *pengajuan_absen.Status == "WaitApproved" {
			if req.Status == "Rejected" {
				pengajuan_absen.Status = &req.Status
				if req.Keterangan != "" {
					pengajuan_absen.Keterangan = &req.Keterangan
				}
				updated_pengajuan, _ := c.PengajuanAbsenRepo.Update(pengajuan_absen)
				history_pengajuan := HistoryPengajuanCutiSet(updated_pengajuan)
				c.HistoryPengajuanAbsenRepo.Create(history_pengajuan)

				ctx.JSON(http.StatusOK, gin.H{
					"status":     http.StatusOK,
					"keterangan": "Success",
				})
			} else if req.Status == "Approved" {
				pengajuan_absen.Status = &req.Status
				if req.Keterangan != "" {
					pengajuan_absen.Keterangan = &req.Keterangan
				}

				tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*pengajuan_absen.TipeAbsenId)

				approve := false
				if tipeAbsen.MaxAbsen != nil {
					approve = true
				} else {
					trc, _ := c.TransaksiCutiRepo.FindDataTransaksiIDPengajuan(pengajuan_absen.IdPengajuanAbsen)

					saldo_hutang := 0
					hutang := false
					for _, transaction := range trc {
						my_saldo, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipe(transaction.Nik, tipeAbsen.IdTipeAbsen, transaction.Periode)

						if transaction.JumlahCuti <= my_saldo.Saldo {
							my_saldo.Saldo = my_saldo.Saldo - transaction.JumlahCuti
							approve = true
						} else if transaction.JumlahCuti <= my_saldo.Saldo+my_saldo.MaxHutang {
							saldo_hutang = transaction.JumlahCuti - my_saldo.Saldo
							my_saldo.MaxHutang = my_saldo.MaxHutang - saldo_hutang
							my_saldo.Saldo = 0
							approve = true
							hutang = true
						} else {
							approve = false
							hutang = false
						}

						if hutang {
							periode, _ := strconv.Atoi(my_saldo.Periode)
							next_periode := strconv.Itoa(periode + (my_saldo.ValidTo.Year() - my_saldo.ValidFrom.Year()))

							next_saldo_periode, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipe(my_saldo.Nik, my_saldo.TipeAbsenId, next_periode)
							if saldo_hutang <= next_saldo_periode.Saldo {
								next_saldo_periode.Saldo = next_saldo_periode.Saldo - saldo_hutang
								approve = true
								eksekusi_saldo = append(eksekusi_saldo, next_saldo_periode)
							} else {
								approve = false
								break
							}
						}

						if approve {
							eksekusi_saldo = append(eksekusi_saldo, my_saldo)
						}
					}
				}

				if approve {
					updated_pengajuan, _ := c.PengajuanAbsenRepo.Update(pengajuan_absen)
					history_pengajuan := HistoryPengajuanCutiSet(updated_pengajuan)
					c.HistoryPengajuanAbsenRepo.Create(history_pengajuan)

					for _, saldo := range eksekusi_saldo {
						updated_saldo, _ := c.SaldoCutiRepo.Update(saldo)
						history_saldo := HistorySaldoCutiSet(updated_saldo)
						c.HistorySaldoCutiRepo.Create(history_saldo)
					}
					ctx.JSON(http.StatusOK, gin.H{
						"status":     http.StatusOK,
						"keterangan": "Success",
					})
				} else {
					*pengajuan_absen.Keterangan = "Di Tolak, Saldo Anda Tidak Cukup"
					*pengajuan_absen.Status = "Rejected"
					updated_pengajuan, _ := c.PengajuanAbsenRepo.Update(pengajuan_absen)
					history_pengajuan := HistoryPengajuanCutiSet(updated_pengajuan)
					c.HistoryPengajuanAbsenRepo.Create(history_pengajuan)

					ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
						"status":     http.StatusServiceUnavailable,
						"keterangan": "Maaf Pengajuan tidak dapat di Approve, Saldo Tidak Cukup, Silahkan Mengajukan kembali",
					})
				}
			}
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
func (c *CutiKrywnController) ListApprvlCuti(ctx *gin.Context) {
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
func (c *CutiKrywnController) ShowDetailApprovalPengajuanCuti(ctx *gin.Context) {
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

// Tipe Saldo Cuti (DONE)
func (c *CutiKrywnController) GetTipeAbsenSaldoPengajuan(ctx *gin.Context) {
	var req Authentication.ValidationNIKTahun

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
	data := []Authentication.GetTipeAbsenSaldoIndiv{}
	data2 := []Authentication.GetTipeAbsenSaldoIndiv{}

	pihc_mstr_krywn, err_pihc_mstr_krywn := c.PihcMasterKaryDbRepo.FindUserByNIK(req.NIK)

	if err_pihc_mstr_krywn == nil {
		fmt.Println(pihc_mstr_krywn.Company)
		TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenPengajuan(pihc_mstr_krywn.Company)

		for _, dataCuti := range TipeAbsen {
			tipeSaldoCuti := Authentication.GetTipeAbsenSaldoIndiv{
				IdTipeAbsen:   dataCuti.IdTipeAbsen,
				NamaTipeAbsen: dataCuti.NamaTipeAbsen,
				TipeMaxAbsen:  *dataCuti.TipeMaxAbsen,
			}
			if dataCuti.MaxAbsen == nil {
				saldoCutiPerTipe, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipeArr(pihc_mstr_krywn.EmpNo, dataCuti.IdTipeAbsen, req.Tahun)

				arr_saldo := []Authentication.SaldoIndiv{}
				for _, my_saldo := range saldoCutiPerTipe {
					saldo := Authentication.SaldoIndiv{
						Saldo:           my_saldo.Saldo,
						ValidFrom:       my_saldo.ValidFrom.Format(time.DateOnly),
						ValidTo:         my_saldo.ValidTo.Format(time.DateOnly),
						Periode:         my_saldo.Periode,
						MaxHutang:       my_saldo.MaxHutang,
						ValidFromHutang: my_saldo.ValidFromHutang.Format(time.DateOnly),
					}
					arr_saldo = append(arr_saldo, saldo)

				}
				tipeSaldoCuti.MySaldo = arr_saldo
				if dataCuti.NamaTipeAbsen == "Cuti Tahunan" {
					data = append(data, tipeSaldoCuti)
				} else {
					data2 = append(data2, tipeSaldoCuti)
				}
			} else {
				tipeSaldoCuti.MaxAbsen = dataCuti.MaxAbsen
				data2 = append(data2, tipeSaldoCuti)
			}
		}

		data = append(data, data2...)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data,
	})
}
func (c *CutiKrywnController) GetAdminTipeAbsen(ctx *gin.Context) {
	nik := ctx.Query("nik")
	data := []Authentication.GetTipeAbsenKaryawanSaldo{}
	data2 := []Authentication.GetTipeAbsenKaryawanSaldo{}

	if nik == "" {
		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": "Nik wajib di isi"})

		return
	}

	pihc_mstr_krywn, err_pihc_mstr_krywn := c.PihcMasterKaryDbRepo.FindUserByNIK(nik)

	if err_pihc_mstr_krywn == nil {
		TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenSaldo(pihc_mstr_krywn.Company)

		for _, dataCuti := range TipeAbsen {
			TipeAbsenKaryawan := Authentication.GetTipeAbsenKaryawanSaldo{
				IdTipeAbsen:   dataCuti.IdTipeAbsen,
				NamaTipeAbsen: dataCuti.NamaTipeAbsen,
				CompCode:      dataCuti.CompCode,
				MaxAbsen:      dataCuti.MaxAbsen,
				TipeMaxAbsen:  dataCuti.TipeMaxAbsen,
				CreatedAt:     dataCuti.CreatedAt,
				UpdatedAt:     dataCuti.UpdatedAt,
			}
			if dataCuti.NamaTipeAbsen == "Cuti Tahunan" {
				data = append(data, TipeAbsenKaryawan)
			} else {
				data2 = append(data2, TipeAbsenKaryawan)
			}
		}

		data = append(data, data2...)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data,
	})
}

// Saldo Cuti (DONE)
func (c *CutiKrywnController) StoreAdminSaldoCutiKaryawan(ctx *gin.Context) {
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

	if req.IDSaldo != nil {
		req.IDSaldo = ConvertInterfaceTypeDataToInt(req.IDSaldo)
	}

	kebenaran := false
	var keterangan string
	if req.IDSaldo == nil {
		sc.TipeAbsenId = req.TipeAbsenId
		sc.Nik = req.Nik
		sc.Saldo = req.Saldo
		sc.ValidFrom, _ = time.Parse(time.DateOnly, req.ValidFrom)
		sc.ValidTo, _ = time.Parse(time.DateOnly, req.ValidTo)
		sc.CreatedBy = req.CreatedBy

		// periode := strconv.Itoa(time.Now().Year())
		periode := strconv.Itoa(sc.ValidFrom.Year())
		sc.Periode = periode
		sc.MaxHutang = req.MaxHutang
		sc.ValidFromHutang, _ = time.Parse(time.DateOnly, req.ValidFromHutang)

		saldoCuti, err := c.SaldoCutiRepo.Create(sc)
		if err == nil {
			data = SaldoCutiKaryawanSet(saldoCuti)
			dataHistorySaldoCuti := HistorySaldoCutiSet(saldoCuti)
			c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)

			kebenaran = true
			keterangan = "Success"
		} else {
			data = Authentication.SaldoCutiKaryawan{}

			kebenaran = false
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
			data = SaldoCutiKaryawanSet(saldoCuti)
			dataHistorySaldoCuti := HistorySaldoCutiSet(saldoCuti)
			c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)

			kebenaran = true
			keterangan = "Success"
		} else {
			data = Authentication.SaldoCutiKaryawan{}

			kebenaran = false
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
func (c *CutiKrywnController) ListAdminSaldoCutiKaryawan(ctx *gin.Context) {
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
			dataSaldoCuti := ListSaldoCutiKaryawanSet(dataSaldoo, company, TipeAbsen)

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
func (c *CutiKrywnController) GetAdminSaldoCuti(ctx *gin.Context) {
	id := ctx.Param("id_saldo_cuti")
	var data Authentication.ListSaldoCutiKaryawan

	id_saldo, _ := strconv.Atoi(id)
	saldoCuti, err := c.SaldoCutiRepo.GetSaldoCutiByID(id_saldo)
	if err == nil {
		karyawan, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(saldoCuti.Nik)
		company, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(karyawan.Company)
		TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(saldoCuti.TipeAbsenId)
		data = ListSaldoCutiKaryawanSet(saldoCuti, company, TipeAbsen)
		if karyawan.Nama != nil && *karyawan.Nama != "" {
			data.Nama = *karyawan.Nama
		}

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
func (c *CutiKrywnController) DeleteAdminSaldoCuti(ctx *gin.Context) {
	id := ctx.Param("id_saldo_cuti")
	var data Authentication.SaldoCutiKaryawan

	id_saldo, _ := strconv.Atoi(id)
	saldoCuti, err := c.SaldoCutiRepo.DelAdminSaldoCuti(id_saldo)

	if err == nil {
		data = SaldoCutiKaryawanSet(saldoCuti)

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
