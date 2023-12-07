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

// Convert Db to struct
func convertSourceTargetDataKaryawan(source pihc.PihcMasterKaryDb) pihc.PihcMasterKary {
	var location *string
	var seksiID *string
	var seksiTitle *string
	var preName *string
	var postName *string
	var NoNPWP *string
	var bankAccount *string
	var bankName *string
	var payScale *string
	if source.Lokasi != "" {
		location = &source.Lokasi
	}
	if source.SeksiID != "" {
		seksiID = &source.SeksiID
	}
	if source.SeksiTitle != "" {
		seksiTitle = &source.SeksiTitle
	}
	if source.PreNameTitle != "" {
		preName = &source.PreNameTitle
	}
	if source.PostNameTitle != "" {
		postName = &source.PostNameTitle
	}
	if source.NoNPWP != "" {
		NoNPWP = &source.NoNPWP
	}
	if source.BankAccount != "" {
		bankAccount = &source.BankAccount
	}
	if source.BankName != "" {
		bankName = &source.BankName
	}
	if source.PayScale != "" {
		payScale = &source.PayScale
	}
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
		SubAreaTitle:   source.SubAreaTitle,
		Contract:       source.Contract,
		Pendidikan:     source.Pendidikan,
		Company:        source.Company,
		Lokasi:         location,
		EmployeeStatus: source.EmployeeStatus,
		Email:          source.Email,
		HP:             source.HP,
		TglLahir:       source.TglLahir.Format("2006-01-02"),
		PosID:          source.PosID,
		PosTitle:       source.PosTitle,
		SupPosID:       source.SupPosID,
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
		SeksiID:        seksiID,
		SeksiTitle:     seksiTitle,
		PreNameTitle:   preName,
		PostNameTitle:  postName,
		NoNPWP:         NoNPWP,
		BankAccount:    bankAccount,
		BankName:       bankName,
		MdgDate:        source.MdgDate,
		PayScale:       payScale,
		CCCode:         source.CCCode,
		Nickname:       source.Nickname,
	}
}
func convertSourceTargetMyPengajuanAbsen(source cuti.PengajuanAbsen, source2 cuti.TipeAbsen) cuti.MyPengajuanAbsen {
	return cuti.MyPengajuanAbsen{
		IdPengajuanAbsen: source.IdPengajuanAbsen,
		Nik:              source.Nik,
		CompCode:         source.CompCode,
		TipeAbsen:        source2,
		Deskripsi:        source.Deskripsi,
		MulaiAbsen:       source.MulaiAbsen.Format(time.DateOnly),
		AkhirAbsen:       source.AkhirAbsen.Format(time.DateOnly),
		TglPengajuan:     source.TglPengajuan.Format(time.DateOnly),
		Status:           source.Status,
		CreatedBy:        source.CreatedBy,
		CreatedAt:        source.CreatedAt,
		UpdatedAt:        source.UpdatedAt,
		Keterangan:       source.Keterangan,
		Periode:          source.Periode,
		ApprovedBy:       source.ApprovedBy,
		JmlHariKalendar:  source.JmlHariKalendar,
		JmlHariKerja:     source.JmlHariKerja,
	}
}

// Pengajuan Cuti
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

	PIHC_MSTR_KRY, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(req.Nik)

	comp_code := PIHC_MSTR_KRY.Company

	if req.IdPengajuanAbsen == 0 {
		// ID PengajuanAbsen == 0
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
		for dataKaryawan.PosTitle != "Wakil Direktur Utama" {
			dataKaryawan, _ = c.PihcMasterKaryDbRepo.FindUserAtasanBySupPosID(dataKaryawan.SupPosID)
		}
		approvedBy := dataKaryawan.EmpNo
		sck.ApprovedBy = &approvedBy

		tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*sck.TipeAbsenId)

		if tipeAbsen.MaxAbsen != nil {
			// MaxAbsen != nil
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

			sckData, _ := c.PengajuanAbsenRepo.Create(sck)

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
			// MaxAbsen == nil (Menggunakan Saldo)
			existSaldo, saldoCuti, _ := c.SaldoCutiRepo.FindExistSaldo2Periode(req.TipeAbsenId, sck.Nik, req.MulaiAbsen, req.AkhirAbsen)

			if existSaldo {
				// Ada Saldo
				isMax := false
				// Menghitung jumlah hari kerja dan hari kalender
				JmlHariKerja := 0
				jmlhHariKalender := 0

				for currentDate := sck.MulaiAbsen; !currentDate.After(sck.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
					jmlhHariKalender++
					if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
						JmlHariKerja++
					}
				}
				fmt.Println("Hari Kalender: ", jmlhHariKalender, ", Hari Kerja:", JmlHariKerja)
				sck.JmlHariKalendar = &jmlhHariKalender
				sck.JmlHariKerja = &JmlHariKerja

				isHutang := false
				var keterangan string
				indexHutang := 0
				var newPeriode time.Time
				totalKerja := 0
				nextyear := false

				// Loop Saldo Cuti
				for _, dataSaldo := range saldoCuti {
					hariKerja := 0
					hariKalender := 0

					if (sck.MulaiAbsen.Before(dataSaldo.ValidTo) || sck.MulaiAbsen.Equal(dataSaldo.ValidTo)) &&
						sck.MulaiAbsen.After(dataSaldo.ValidFrom) && (sck.AkhirAbsen.After(dataSaldo.ValidTo) || sck.AkhirAbsen.Equal(dataSaldo.ValidTo)) {
						// MulaiAbsen <= ValidTo && MulaiAbsen > ValidFrom && AkhirAbsen>=ValidTo
						fmt.Println("A")
						for currentDate := sck.MulaiAbsen; !currentDate.After(dataSaldo.ValidTo.AddDate(0, 0, -1)); currentDate = currentDate.AddDate(0, 0, 1) {
							hariKalender++
							if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
								hariKerja++
								fmt.Print(hariKerja, " ")
							}
						}
						fmt.Println()
						newPeriode = dataSaldo.ValidTo
					} else if (sck.MulaiAbsen.After(dataSaldo.ValidFrom) || sck.MulaiAbsen.Equal(dataSaldo.ValidFrom)) &&
						(sck.AkhirAbsen.Before(dataSaldo.ValidTo) || sck.AkhirAbsen.Equal(dataSaldo.ValidTo)) {
						// MulaiAbsen >= ValidFrom && AkhirAbsen <= ValidTo
						fmt.Println("B")
						if sck.MulaiAbsen.After(dataSaldo.ValidFrom) && sck.AkhirAbsen.Before(dataSaldo.ValidTo) {
							for currentDate := sck.MulaiAbsen; !currentDate.After(sck.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
								hariKalender++
								if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
									hariKerja++
									fmt.Print(hariKerja, " ")
								}
							}
							nextyear = true
						} else {
							for currentDate := dataSaldo.ValidFrom; !currentDate.After(sck.MulaiAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
								hariKalender++
								if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
									hariKerja++
									fmt.Print(hariKerja, " ")
								}
							}
							for currentDate := sck.MulaiAbsen.AddDate(0, 0, 1); !currentDate.After(sck.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
								hariKalender++
								if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
									hariKerja++
									fmt.Print(hariKerja, " ")
								}
							}
						}
						fmt.Println()
					} else if (newPeriode.After(dataSaldo.ValidFrom) || newPeriode.Equal(dataSaldo.ValidFrom)) &&
						(sck.AkhirAbsen.Before(dataSaldo.ValidTo) || sck.AkhirAbsen.Equal(dataSaldo.ValidTo)) {
						// newPeriode>=ValidFrom && AkhirAbsen <= ValidTo(periode ke-2)
						fmt.Println("C")
						for currentDate := dataSaldo.ValidFrom; !currentDate.After(newPeriode); currentDate = currentDate.AddDate(0, 0, 1) {
							hariKalender++
							if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
								hariKerja++
								fmt.Print(hariKerja, " ")
							}
						}
						for currentDate := newPeriode.AddDate(0, 0, 1); !currentDate.After(sck.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
							hariKalender++
							if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
								hariKerja++
								fmt.Print(hariKerja, " ")
							}
						}
						fmt.Println()
					}

					if isHutang {
						// periode ke-2
						if dataSaldo.Saldo-indexHutang >= 0 {
							dataSaldo.Saldo = dataSaldo.Saldo - indexHutang
						}
					}

					fmt.Println(hariKalender, hariKerja, dataSaldo.Saldo, dataSaldo.MaxHutang, dataSaldo.Periode, totalKerja)

					if dataSaldo.Saldo != 0 && hariKerja != 0 {
						if hariKerja <= dataSaldo.Saldo {
							fmt.Println("hariKerja <= dataSaldo.Saldo")
							isMax = true
							dataSaldo.Saldo = dataSaldo.Saldo - hariKerja
						} else if hariKerja <= (dataSaldo.Saldo+dataSaldo.MaxHutang) && dataSaldo.MaxHutang != 0 {
							fmt.Println("hariKerja <= (dataSaldo.Saldo+dataSaldo.MaxHutang) && dataSaldo.MaxHutang != 0")
							if sck.MulaiAbsen.After(dataSaldo.ValidFromHutang) || sck.MulaiAbsen.Equal(dataSaldo.ValidFromHutang) {
								isMax = true
								isHutang = isMax
								hutang := hariKerja - dataSaldo.Saldo
								indexHutang = hutang
								dataSaldo.Saldo = 0
								dataSaldo.MaxHutang = dataSaldo.MaxHutang - hutang
								fmt.Println(hariKerja, dataSaldo.Saldo, dataSaldo.MaxHutang, hutang)
								if nextyear {
									fmt.Println("nextyear")
									tahun, _ := strconv.Atoi(dataSaldo.Periode)
									scNextYear, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipe(dataSaldo.Nik, dataSaldo.TipeAbsenId, strconv.Itoa(tahun+1))
									if scNextYear.Saldo-indexHutang >= 0 {
										fmt.Println(scNextYear.Saldo)
										scNextYear.Saldo = scNextYear.Saldo - indexHutang
										fmt.Println(scNextYear.Saldo)
										c.SaldoCutiRepo.Update(scNextYear)
										dataHistorySaldoCuti := cuti.HistorySaldoCuti{
											IdHistorySaldoCuti: scNextYear.IdSaldoCuti,
											TipeAbsenId:        scNextYear.TipeAbsenId,
											Nik:                scNextYear.Nik,
											Saldo:              scNextYear.Saldo,
											ValidFrom:          scNextYear.ValidFrom,
											ValidTo:            scNextYear.ValidTo,
											CreatedBy:          scNextYear.CreatedBy,
											CreatedAt:          scNextYear.CreatedAt,
											UpdatedAt:          scNextYear.UpdatedAt,
											Periode:            scNextYear.Periode,
											MaxHutang:          scNextYear.MaxHutang,
											ValidFromHutang:    scNextYear.ValidFromHutang,
										}
										c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)
									}
								}
							} else {
								fmt.Println("X")
								isMax = false
								keterangan = "Berada di luar Masa Berlaku Hutang"
							}
						} else {
							fmt.Println("Y")
							isMax = false
							keterangan = "Saldo Tidak Cukup"
						}
					} else if hariKerja != 0 && dataSaldo.MaxHutang != 0 {
						if hariKerja <= dataSaldo.MaxHutang {
							fmt.Println("hariKerja <= dataSaldo.MaxHutang")
							isMax = true
							isHutang = isMax
							hutang := hariKerja - dataSaldo.Saldo
							indexHutang = hutang
							dataSaldo.Saldo = 0
							dataSaldo.MaxHutang = dataSaldo.MaxHutang - hutang
							if nextyear {
								fmt.Println("nextyear")
								tahun, _ := strconv.Atoi(dataSaldo.Periode)
								scNextYear, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipe(dataSaldo.Nik, dataSaldo.TipeAbsenId, strconv.Itoa(tahun+1))
								if scNextYear.Saldo-indexHutang >= 0 {
									fmt.Println(scNextYear.Saldo)
									scNextYear.Saldo = scNextYear.Saldo - indexHutang
									fmt.Println(scNextYear.Saldo)
									c.SaldoCutiRepo.Update(scNextYear)
									dataHistorySaldoCuti := cuti.HistorySaldoCuti{
										IdHistorySaldoCuti: scNextYear.IdSaldoCuti,
										TipeAbsenId:        scNextYear.TipeAbsenId,
										Nik:                scNextYear.Nik,
										Saldo:              scNextYear.Saldo,
										ValidFrom:          scNextYear.ValidFrom,
										ValidTo:            scNextYear.ValidTo,
										CreatedBy:          scNextYear.CreatedBy,
										CreatedAt:          scNextYear.CreatedAt,
										UpdatedAt:          scNextYear.UpdatedAt,
										Periode:            scNextYear.Periode,
										MaxHutang:          scNextYear.MaxHutang,
										ValidFromHutang:    scNextYear.ValidFromHutang,
									}
									c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)
								}
							}
							fmt.Println(hariKerja, dataSaldo.Saldo, dataSaldo.MaxHutang, hutang)
						} else {
							fmt.Println("ZZZ")
							isMax = false
							keterangan = "Saldo Hutang Tidak Cukup"
						}

					} else {
						fmt.Println("Z")
						isMax = false
						keterangan = "Saldo Tidak Cukup"
					}

					if isMax {
						totalKerja = totalKerja + hariKerja
						source := Authentication.SaldoCutiTransaksiPengajuan{
							SaldoCuti: dataSaldo,
							JmlhCuti:  hariKerja,
						}
						trsc = append(trsc, source)
					}
				}

				// CREATE
				if isMax {
					if totalKerja == JmlHariKerja {
						sckData, _ := c.PengajuanAbsenRepo.Create(sck)

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
						files, _ := c.FileAbsenRepo.CreateArr(fsc)
						for _, transaction := range trsc {
							c.SaldoCutiRepo.Update(transaction.SaldoCuti)
							fmt.Println("XX")
							if !nextyear {
								fmt.Println("XXX")
								dataHistorySaldoCuti := cuti.HistorySaldoCuti{
									IdHistorySaldoCuti: transaction.SaldoCuti.IdSaldoCuti,
									TipeAbsenId:        transaction.SaldoCuti.TipeAbsenId,
									Nik:                transaction.SaldoCuti.Nik,
									Saldo:              transaction.SaldoCuti.Saldo,
									ValidFrom:          transaction.SaldoCuti.ValidFrom,
									ValidTo:            transaction.SaldoCuti.ValidTo,
									CreatedBy:          transaction.SaldoCuti.CreatedBy,
									CreatedAt:          transaction.SaldoCuti.CreatedAt,
									UpdatedAt:          transaction.SaldoCuti.UpdatedAt,
									Periode:            transaction.SaldoCuti.Periode,
									MaxHutang:          transaction.SaldoCuti.MaxHutang,
									ValidFromHutang:    transaction.SaldoCuti.ValidFromHutang,
								}
								c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)
							}

							if transaction.JmlhCuti != 0 {
								transaksi_cuti := cuti.TransaksiCuti{
									PengajuanAbsenId: sckData.IdPengajuanAbsen,
									Nik:              sckData.Nik,
									Periode:          transaction.Periode,
									JumlahCuti:       transaction.JmlhCuti,
									TipeHari:         "hari_kerja",
								}
								c.TransaksiCutiRepo.Create(transaksi_cuti)
							}
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
					} else {
						ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
							"status":     http.StatusServiceUnavailable,
							"keterangan": keterangan,
						})
					}
				} else {
					ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
						"status":     http.StatusServiceUnavailable,
						"keterangan": keterangan,
					})
				}
			} else {
				ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"status":     http.StatusServiceUnavailable,
					"keterangan": "Berada diluar Masa Berlaku",
				})
			}
		}
	} else {
		// ID PengajuanAbsen != 0
		pengajuan_absen, _ := c.PengajuanAbsenRepo.FindDataIdPengajuan(req.IdPengajuanAbsen)
		tipe_absen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*pengajuan_absen.TipeAbsenId)
		file_absen, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuan(pengajuan_absen.IdPengajuanAbsen)

		if *pengajuan_absen.Status == "WaitApproved" {
			fmt.Println("A")
			if tipe_absen.MaxAbsen == nil {
				fmt.Println("B")
				// Transaksi
				transaksi, _ := c.TransaksiCutiRepo.FindDataTransaksiIDPengajuan(pengajuan_absen.IdPengajuanAbsen)
				for _, data_transaksi := range transaksi {
					saldo_cuti, err_saldo := c.SaldoCutiRepo.GetSaldoCutiPerTipe(pengajuan_absen.Nik, *pengajuan_absen.TipeAbsenId, data_transaksi.Periode)
					if err_saldo == nil {
						saldo_cuti.Saldo = saldo_cuti.Saldo + data_transaksi.JumlahCuti
						c.SaldoCutiRepo.Update(saldo_cuti)
						dataHistorySaldoCuti := cuti.HistorySaldoCuti{
							IdHistorySaldoCuti: saldo_cuti.IdSaldoCuti,
							TipeAbsenId:        saldo_cuti.TipeAbsenId,
							Nik:                saldo_cuti.Nik,
							Saldo:              saldo_cuti.Saldo,
							ValidFrom:          saldo_cuti.ValidFrom,
							ValidTo:            saldo_cuti.ValidTo,
							CreatedBy:          saldo_cuti.CreatedBy,
							CreatedAt:          saldo_cuti.CreatedAt,
							UpdatedAt:          saldo_cuti.UpdatedAt,
							Periode:            saldo_cuti.Periode,
							MaxHutang:          saldo_cuti.MaxHutang,
							ValidFromHutang:    saldo_cuti.ValidFromHutang,
						}
						c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)
						c.TransaksiCutiRepo.Delete(data_transaksi)
					}
				}
			}
			for _, delete_file := range file_absen {
				c.FileAbsenRepo.Delete(delete_file)
			}
		} else if *pengajuan_absen.Status == "Approved" {
			*pengajuan_absen.Status = "WaitApproved"
			if tipe_absen.MaxAbsen == nil {
				// Transaksi
				transaksi, _ := c.TransaksiCutiRepo.FindDataTransaksiIDPengajuan(pengajuan_absen.IdPengajuanAbsen)
				for _, data_transaksi := range transaksi {
					saldo_cuti, err_saldo := c.SaldoCutiRepo.GetSaldoCutiPerTipe(pengajuan_absen.Nik, *pengajuan_absen.TipeAbsenId, data_transaksi.Periode)
					if err_saldo == nil {
						saldo_cuti.Saldo = saldo_cuti.Saldo + data_transaksi.JumlahCuti
						c.SaldoCutiRepo.Update(saldo_cuti)
						dataHistorySaldoCuti := cuti.HistorySaldoCuti{
							IdHistorySaldoCuti: saldo_cuti.IdSaldoCuti,
							TipeAbsenId:        saldo_cuti.TipeAbsenId,
							Nik:                saldo_cuti.Nik,
							Saldo:              saldo_cuti.Saldo,
							ValidFrom:          saldo_cuti.ValidFrom,
							ValidTo:            saldo_cuti.ValidTo,
							CreatedBy:          saldo_cuti.CreatedBy,
							CreatedAt:          saldo_cuti.CreatedAt,
							UpdatedAt:          saldo_cuti.UpdatedAt,
							Periode:            saldo_cuti.Periode,
							MaxHutang:          saldo_cuti.MaxHutang,
							ValidFromHutang:    saldo_cuti.ValidFromHutang,
						}
						c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)
						c.TransaksiCutiRepo.Delete(data_transaksi)
					}
				}
			}
			for _, delete_file := range file_absen {
				c.FileAbsenRepo.Delete(delete_file)
			}
			pengajuan_absen.Keterangan = nil
		} else {
			*pengajuan_absen.Status = "WaitApproved"
			pengajuan_absen.Keterangan = nil
			for _, delete_file := range file_absen {
				c.FileAbsenRepo.Delete(delete_file)
			}
		}
		pengajuan_absen.TipeAbsenId = &req.TipeAbsenId
		pengajuan_absen.CompCode = comp_code
		pengajuan_absen.Deskripsi = &req.Deskripsi
		pengajuan_absen.MulaiAbsen, _ = time.Parse(time.DateOnly, req.MulaiAbsen)
		pengajuan_absen.AkhirAbsen, _ = time.Parse(time.DateOnly, req.AkhirAbsen)
		pengajuan_absen.TglPengajuan, _ = time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
		periode := strconv.Itoa(time.Now().Year())
		pengajuan_absen.Periode = &periode
		pengajuan_absen.CreatedBy = &req.CreatedBy
		approvedBy := "82105096"
		pengajuan_absen.ApprovedBy = &approvedBy

		tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*pengajuan_absen.TipeAbsenId)

		if tipeAbsen.MaxAbsen != nil {
			// MaxAbsen != nil
			JmlHariKerja := 0
			jmlhHariKalender := 0

			if *tipeAbsen.TipeMaxAbsen == "hari_kalender" {
				for currentDate := pengajuan_absen.MulaiAbsen; jmlhHariKalender != *tipeAbsen.MaxAbsen; currentDate = currentDate.AddDate(0, 0, 1) {
					jmlhHariKalender++
					if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
						JmlHariKerja++
					}
				}
				pengajuan_absen.AkhirAbsen = pengajuan_absen.MulaiAbsen.AddDate(0, 0, jmlhHariKalender-1)
			} else if *tipeAbsen.TipeMaxAbsen == "hari_kerja" {
				for currentDate := pengajuan_absen.MulaiAbsen; JmlHariKerja != *tipeAbsen.MaxAbsen; currentDate = currentDate.AddDate(0, 0, 1) {
					jmlhHariKalender++
					if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
						JmlHariKerja++
					}
				}
				pengajuan_absen.AkhirAbsen = pengajuan_absen.MulaiAbsen.AddDate(0, 0, jmlhHariKalender-1)
			}
			pengajuan_absen.JmlHariKalendar = &jmlhHariKalender
			pengajuan_absen.JmlHariKerja = &JmlHariKerja
			sckData, _ := c.PengajuanAbsenRepo.Update(pengajuan_absen)

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
			files, _ := c.FileAbsenRepo.CreateArr(fsc)

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
			// MaxAbsen == nil (Menggunakan Saldo)
			existSaldo, saldoCuti, _ := c.SaldoCutiRepo.FindExistSaldo2Periode(req.TipeAbsenId, pengajuan_absen.Nik, req.MulaiAbsen, req.AkhirAbsen)

			if existSaldo {
				// Ada Saldo
				isMax := false
				// Menghitung jumlah hari kerja dan hari kalender
				JmlHariKerja := 0
				jmlhHariKalender := 0

				for currentDate := pengajuan_absen.MulaiAbsen; !currentDate.After(pengajuan_absen.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
					jmlhHariKalender++
					if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
						JmlHariKerja++
					}
				}
				fmt.Println("Hari Kalender: ", jmlhHariKalender, ", Hari Kerja:", JmlHariKerja)
				pengajuan_absen.JmlHariKalendar = &jmlhHariKalender
				pengajuan_absen.JmlHariKerja = &JmlHariKerja

				isHutang := false
				var keterangan string
				indexHutang := 0
				var newPeriode time.Time
				totalKerja := 0
				nextyear := false

				// Loop Saldo Cuti
				for _, dataSaldo := range saldoCuti {
					hariKerja := 0
					hariKalender := 0

					if (pengajuan_absen.MulaiAbsen.Before(dataSaldo.ValidTo) || pengajuan_absen.MulaiAbsen.Equal(dataSaldo.ValidTo)) &&
						pengajuan_absen.MulaiAbsen.After(dataSaldo.ValidFrom) && (pengajuan_absen.AkhirAbsen.After(dataSaldo.ValidTo) || pengajuan_absen.AkhirAbsen.Equal(dataSaldo.ValidTo)) {
						// MulaiAbsen <= ValidTo && MulaiAbsen > ValidFrom && AkhirAbsen>=ValidTo
						fmt.Println("A")
						for currentDate := pengajuan_absen.MulaiAbsen; !currentDate.After(dataSaldo.ValidTo.AddDate(0, 0, -1)); currentDate = currentDate.AddDate(0, 0, 1) {
							hariKalender++
							if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
								hariKerja++
								fmt.Print(hariKerja, " ")
							}
						}
						fmt.Println()
						newPeriode = dataSaldo.ValidTo
					} else if (pengajuan_absen.MulaiAbsen.After(dataSaldo.ValidFrom) || pengajuan_absen.MulaiAbsen.Equal(dataSaldo.ValidFrom)) &&
						(pengajuan_absen.AkhirAbsen.Before(dataSaldo.ValidTo) || pengajuan_absen.AkhirAbsen.Equal(dataSaldo.ValidTo)) {
						// MulaiAbsen >= ValidFrom && AkhirAbsen <= ValidTo
						fmt.Println("B")
						if pengajuan_absen.MulaiAbsen.After(dataSaldo.ValidFrom) && pengajuan_absen.AkhirAbsen.Before(dataSaldo.ValidTo) {
							for currentDate := pengajuan_absen.MulaiAbsen; !currentDate.After(pengajuan_absen.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
								hariKalender++
								if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
									hariKerja++
									fmt.Print(hariKerja, " ")
								}
							}
							nextyear = true
						} else {
							for currentDate := dataSaldo.ValidFrom; !currentDate.After(pengajuan_absen.MulaiAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
								hariKalender++
								if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
									hariKerja++
									fmt.Print(hariKerja, " ")
								}
							}
							for currentDate := pengajuan_absen.MulaiAbsen.AddDate(0, 0, 1); !currentDate.After(pengajuan_absen.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
								hariKalender++
								if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
									hariKerja++
									fmt.Print(hariKerja, " ")
								}
							}
						}
						fmt.Println()
					} else if (newPeriode.After(dataSaldo.ValidFrom) || newPeriode.Equal(dataSaldo.ValidFrom)) &&
						(pengajuan_absen.AkhirAbsen.Before(dataSaldo.ValidTo) || pengajuan_absen.AkhirAbsen.Equal(dataSaldo.ValidTo)) {
						// newPeriode>=ValidFrom && AkhirAbsen <= ValidTo(periode ke-2)
						fmt.Println("C")
						for currentDate := dataSaldo.ValidFrom; !currentDate.After(newPeriode); currentDate = currentDate.AddDate(0, 0, 1) {
							hariKalender++
							if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
								hariKerja++
								fmt.Print(hariKerja, " ")
							}
						}
						for currentDate := newPeriode.AddDate(0, 0, 1); !currentDate.After(pengajuan_absen.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
							hariKalender++
							if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
								hariKerja++
								fmt.Print(hariKerja, " ")
							}
						}
						fmt.Println()
					}

					if isHutang {
						// periode ke-2
						if dataSaldo.Saldo-indexHutang >= 0 {
							dataSaldo.Saldo = dataSaldo.Saldo - indexHutang
						}
					}

					fmt.Println(hariKalender, hariKerja, dataSaldo.Saldo, dataSaldo.MaxHutang, dataSaldo.Periode, totalKerja)

					if dataSaldo.Saldo != 0 && hariKerja != 0 {
						if hariKerja <= dataSaldo.Saldo {
							fmt.Println("hariKerja <= dataSaldo.Saldo")
							isMax = true
							dataSaldo.Saldo = dataSaldo.Saldo - hariKerja
						} else if hariKerja <= (dataSaldo.Saldo+dataSaldo.MaxHutang) && dataSaldo.MaxHutang != 0 {
							fmt.Println("hariKerja <= (dataSaldo.Saldo+dataSaldo.MaxHutang) && dataSaldo.MaxHutang != 0")
							if pengajuan_absen.MulaiAbsen.After(dataSaldo.ValidFromHutang) || pengajuan_absen.MulaiAbsen.Equal(dataSaldo.ValidFromHutang) {
								isMax = true
								isHutang = isMax
								hutang := hariKerja - dataSaldo.Saldo
								indexHutang = hutang
								dataSaldo.Saldo = 0
								dataSaldo.MaxHutang = dataSaldo.MaxHutang - hutang
								fmt.Println(hariKerja, dataSaldo.Saldo, dataSaldo.MaxHutang, hutang)
								if nextyear {
									fmt.Println("nextyear")
									tahun, _ := strconv.Atoi(dataSaldo.Periode)
									scNextYear, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipe(dataSaldo.Nik, dataSaldo.TipeAbsenId, strconv.Itoa(tahun+1))
									if scNextYear.Saldo-indexHutang >= 0 {
										fmt.Println(scNextYear.Saldo)
										scNextYear.Saldo = scNextYear.Saldo - indexHutang
										fmt.Println(scNextYear.Saldo)
										c.SaldoCutiRepo.Update(scNextYear)
										dataHistorySaldoCuti := cuti.HistorySaldoCuti{
											IdHistorySaldoCuti: scNextYear.IdSaldoCuti,
											TipeAbsenId:        scNextYear.TipeAbsenId,
											Nik:                scNextYear.Nik,
											Saldo:              scNextYear.Saldo,
											ValidFrom:          scNextYear.ValidFrom,
											ValidTo:            scNextYear.ValidTo,
											CreatedBy:          scNextYear.CreatedBy,
											CreatedAt:          scNextYear.CreatedAt,
											UpdatedAt:          scNextYear.UpdatedAt,
											Periode:            scNextYear.Periode,
											MaxHutang:          scNextYear.MaxHutang,
											ValidFromHutang:    scNextYear.ValidFromHutang,
										}
										c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)
									}
								}
							} else {
								fmt.Println("X")
								isMax = false
								keterangan = "Berada di luar Masa Berlaku Hutang"
							}
						} else {
							fmt.Println("Y")
							isMax = false
							keterangan = "Saldo Tidak Cukup"
						}
					} else if hariKerja != 0 && dataSaldo.MaxHutang != 0 {
						if hariKerja <= dataSaldo.MaxHutang {
							fmt.Println("hariKerja <= dataSaldo.MaxHutang")
							isMax = true
							isHutang = isMax
							hutang := hariKerja - dataSaldo.Saldo
							indexHutang = hutang
							dataSaldo.Saldo = 0
							dataSaldo.MaxHutang = dataSaldo.MaxHutang - hutang
							if nextyear {
								fmt.Println("nextyear")
								tahun, _ := strconv.Atoi(dataSaldo.Periode)
								scNextYear, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipe(dataSaldo.Nik, dataSaldo.TipeAbsenId, strconv.Itoa(tahun+1))
								if scNextYear.Saldo-indexHutang >= 0 {
									fmt.Println(scNextYear.Saldo)
									scNextYear.Saldo = scNextYear.Saldo - indexHutang
									fmt.Println(scNextYear.Saldo)
									c.SaldoCutiRepo.Update(scNextYear)
									dataHistorySaldoCuti := cuti.HistorySaldoCuti{
										IdHistorySaldoCuti: scNextYear.IdSaldoCuti,
										TipeAbsenId:        scNextYear.TipeAbsenId,
										Nik:                scNextYear.Nik,
										Saldo:              scNextYear.Saldo,
										ValidFrom:          scNextYear.ValidFrom,
										ValidTo:            scNextYear.ValidTo,
										CreatedBy:          scNextYear.CreatedBy,
										CreatedAt:          scNextYear.CreatedAt,
										UpdatedAt:          scNextYear.UpdatedAt,
										Periode:            scNextYear.Periode,
										MaxHutang:          scNextYear.MaxHutang,
										ValidFromHutang:    scNextYear.ValidFromHutang,
									}
									c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)
								}
							}
							fmt.Println(hariKerja, dataSaldo.Saldo, dataSaldo.MaxHutang, hutang)
						} else {
							fmt.Println("ZZZ")
							isMax = false
							keterangan = "Saldo Hutang Tidak Cukup"
						}
					} else {
						fmt.Println("Z")
						isMax = false
						keterangan = "Saldo Tidak Cukup"
					}
					if isMax {
						totalKerja = totalKerja + hariKerja
						source := Authentication.SaldoCutiTransaksiPengajuan{
							SaldoCuti: dataSaldo,
							JmlhCuti:  hariKerja,
						}
						trsc = append(trsc, source)
					}
				}

				// Update
				if isMax {
					if totalKerja == JmlHariKerja {
						sckData, _ := c.PengajuanAbsenRepo.Update(pengajuan_absen)

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
						files, _ := c.FileAbsenRepo.CreateArr(fsc)
						for _, transaction := range trsc {
							c.SaldoCutiRepo.Update(transaction.SaldoCuti)
							if !nextyear {
								dataHistorySaldoCuti := cuti.HistorySaldoCuti{
									IdHistorySaldoCuti: transaction.SaldoCuti.IdSaldoCuti,
									TipeAbsenId:        transaction.SaldoCuti.TipeAbsenId,
									Nik:                transaction.SaldoCuti.Nik,
									Saldo:              transaction.SaldoCuti.Saldo,
									ValidFrom:          transaction.SaldoCuti.ValidFrom,
									ValidTo:            transaction.SaldoCuti.ValidTo,
									CreatedBy:          transaction.SaldoCuti.CreatedBy,
									CreatedAt:          transaction.SaldoCuti.CreatedAt,
									UpdatedAt:          transaction.SaldoCuti.UpdatedAt,
									Periode:            transaction.SaldoCuti.Periode,
									MaxHutang:          transaction.SaldoCuti.MaxHutang,
									ValidFromHutang:    transaction.SaldoCuti.ValidFromHutang,
								}
								c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)
							}

							if transaction.JmlhCuti != 0 {
								transaksi_cuti := cuti.TransaksiCuti{
									PengajuanAbsenId: sckData.IdPengajuanAbsen,
									Nik:              sckData.Nik,
									Periode:          transaction.Periode,
									JumlahCuti:       transaction.JmlhCuti,
								}
								c.TransaksiCutiRepo.Create(transaksi_cuti)
							}
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
					} else {
						ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
							"status":     http.StatusServiceUnavailable,
							"keterangan": keterangan,
						})
					}
				} else {
					ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
						"status":     http.StatusServiceUnavailable,
						"keterangan": keterangan,
					})
				}
			} else {
				ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"status":     http.StatusServiceUnavailable,
					"keterangan": "Berada diluar Masa Berlaku",
				})
			}
		}
	}
}
func (c *CutiKrywnController) ShowDetailPengajuanCuti(ctx *gin.Context) {
	data := Authentication.PengajuanAbsens{}
	id := ctx.Param("id_pengajuan_absen")
	id_pengajuan, _ := strconv.Atoi(id)

	data_pengajuan, err_pengajuan := c.PengajuanAbsenRepo.FindDataIdPengajuan(id_pengajuan)
	data_tipe_absen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*data_pengajuan.TipeAbsenId)
	data_file_absen, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuan(data_pengajuan.IdPengajuanAbsen)

	convert := convertSourceTargetMyPengajuanAbsen(data_pengajuan, data_tipe_absen)

	if data_file_absen == nil {
		data_file_absen = []cuti.FileAbsen{}
	}
	data.MyPengajuanAbsen = convert
	data.File = data_file_absen

	if err_pengajuan == nil {
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
	var data []cuti.MyPengajuanAbsen

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

	for _, myCuti := range dataDB {
		tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*myCuti.TipeAbsenId)

		result := convertSourceTargetMyPengajuanAbsen(myCuti, tipeAbsen)
		data = append(data, result)
	}

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
func (c *CutiKrywnController) DeletePengajuanCuti(ctx *gin.Context) {
	id := ctx.Param("id_pengajuan_absen")

	id_pengajuan_absen, _ := strconv.Atoi(id)
	pengajuanAbsen, err := c.PengajuanAbsenRepo.DelPengajuanCuti(id_pengajuan_absen)
	if *pengajuanAbsen.Status == "WaitApproved" {
		tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*pengajuanAbsen.TipeAbsenId)
		transaksi_cuti, _ := c.TransaksiCutiRepo.FindDataTransaksiIDPengajuan(pengajuanAbsen.IdPengajuanAbsen)
		if tipeAbsen.MaxAbsen == nil {
			for _, tr := range transaksi_cuti {
				saldo, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipe(tr.Nik, tipeAbsen.IdTipeAbsen, tr.Periode)
				saldo.Saldo = saldo.Saldo + tr.JumlahCuti
				c.SaldoCutiRepo.Update(saldo)
				c.TransaksiCutiRepo.Delete(tr)
			}
		} else {
			for _, tr := range transaksi_cuti {
				c.TransaksiCutiRepo.Delete(tr)
			}
		}
	}

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"info":   "success",
		})
	} else {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// Approval Cuti
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
		fmt.Println("A")
		return
	}

	var dataDB []cuti.PengajuanAbsen
	var err error

	if req.Status == "WaitApproved" {
		dB, errs := c.PengajuanAbsenRepo.FindDataNIKPeriodeApprovalWaiting(req.NIK, req.Tahun, req.Status)
		err = errs
		dataDB = dB
	} else {
		dB, errs := c.PengajuanAbsenRepo.FindDataNIKPeriodeApproval(req.NIK, req.Tahun)
		err = errs
		dataDB = dB
	}
	// status := "WaitApproved"

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

		// for _, myCuti := range dataDB {
		// 	tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*myCuti.TipeAbsenId)
		// 	karyawan, _ := c.PihcMasterKaryDbRepo.FindUserByNIK(myCuti.Nik)
		// companys, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(karyawan.Company)
		// files, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuan(myCuti.IdPengajuanAbsen)
		// 	if files == nil {
		// 		files = []cuti.FileAbsen{}
		// 	}

		// 	data_karyawan_convert := convertSourceTargetDataKaryawan(karyawan)

		// 	result := convertSourceTargetMyPengajuanAbsen(myCuti, tipeAbsen)
		// 	foto := "https://storage.googleapis.com/lumen-oauth-storage/DataKaryawan/Foto/" + companys.Code + "/" + result.Nik + ".jpg"
		// 	respons, err := http.Get(foto)
		// 	if err != nil || respons.StatusCode != http.StatusOK {
		// 		foto = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
		// 	}
		// 	list_pengajuan := Authentication.ListApprovalCuti{
		// 		IdPengajuanAbsen:  result.IdPengajuanAbsen,
		// 		PihcMasterKary:    data_karyawan_convert,
		// 		PihcMasterCompany: companys,
		// 		TipeAbsen:         tipeAbsen,
		// 		MulaiAbsen:        result.MulaiAbsen,
		// 		AkhirAbsen:        result.AkhirAbsen,
		// 		TglPengajuan:      result.TglPengajuan,
		// 		FileAbsen:         files,
		// 		Periode:           result.Periode,
		// 		Foto:              foto,
		// 		FotoDefault:       "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg",
		// 	}
		// 	if result.Status != nil && *result.Status != "" {
		// 		list_pengajuan.Status = *result.Status
		// 	}
		// 	if result.Deskripsi != nil && *result.Deskripsi != "" {
		// 		list_pengajuan.Deskripsi = *result.Deskripsi
		// 	}
		// 	list_aprvl = append(list_aprvl, list_pengajuan)

		// }
		fmt.Println("Success")
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

func (c *CutiKrywnController) StoreApprovePengajuanAbsen(ctx *gin.Context) {
	var req Authentication.ValidationApprovalAtasanPengajuanAbsen

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
		fmt.Println("Masuk")
		if req.Status == "Approved" {
			fmt.Println("Approved")
			if req.Status != "" {
				pengajuan_absen.Status = &req.Status
			}
			if req.Keterangan != "" {
				pengajuan_absen.Keterangan = &req.Keterangan
			}

			c.PengajuanAbsenRepo.Update(pengajuan_absen)

			ctx.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"info":   "Success",
			})
		} else {
			fmt.Println("Declined")
			transaksi, _ := c.TransaksiCutiRepo.FindDataTransaksiIDPengajuan(pengajuan_absen.IdPengajuanAbsen)
			for _, data_transaksi := range transaksi {
				fmt.Println("Transaksi")
				// tipe_absen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*pengajuan_absen.TipeAbsenId)
				// if tipe_absen.MaxAbsen != nil {
				saldo_cuti, err_saldo := c.SaldoCutiRepo.GetSaldoCutiPerTipe(pengajuan_absen.Nik, *pengajuan_absen.TipeAbsenId, data_transaksi.Periode)
				if err_saldo == nil {
					saldo_cuti.Saldo = saldo_cuti.Saldo + data_transaksi.JumlahCuti
					c.SaldoCutiRepo.Update(saldo_cuti)
				}
				// }
			}
			if req.Status != "" {
				pengajuan_absen.Status = &req.Status
			}
			if req.Keterangan != "" {
				pengajuan_absen.Keterangan = &req.Keterangan
			}

			c.PengajuanAbsenRepo.Update(pengajuan_absen)

			// c.TransaksiCutiRepo.Delete(pengajuan_absen.IdPengajuanAbsen)

			ctx.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"info":   "Success",
			})
		}
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

// Tipe Saldo Cuti
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
			saldoCutiPerTipe, err := c.SaldoCutiRepo.GetSaldoCutiPerTipe(pihc_mstr_krywn.EmpNo, dataCuti.IdTipeAbsen, req.Tahun)
			if err == nil {
				if dataCuti.NamaTipeAbsen == "Cuti Tahunan" {
					max_absen := &Authentication.MaxAbsenIndiv{
						TipeMaxAbsen: dataCuti.TipeMaxAbsen,
					}
					saldo := &Authentication.SaldoIndiv{
						Saldo:           saldoCutiPerTipe.Saldo,
						ValidFrom:       saldoCutiPerTipe.ValidFrom.Format(time.DateOnly),
						ValidTo:         saldoCutiPerTipe.ValidTo.Format(time.DateOnly),
						Periode:         saldoCutiPerTipe.Periode,
						MaxHutang:       saldoCutiPerTipe.MaxHutang,
						ValidFromHutang: saldoCutiPerTipe.ValidFromHutang.Format(time.DateOnly),
					}
					tipeSaldoCuti := Authentication.GetTipeAbsenSaldoIndiv{
						IdTipeAbsen:   dataCuti.IdTipeAbsen,
						NamaTipeAbsen: dataCuti.NamaTipeAbsen,
					}
					if dataCuti.MaxAbsen != nil && *dataCuti.MaxAbsen != 0 {
						max_absen.MaxAbsen = *dataCuti.MaxAbsen
						tipeSaldoCuti.MaxAbsenIndiv = max_absen
					}
					if saldo.Periode != "" {
						tipeSaldoCuti.SaldoIndiv = saldo
					}

					data = append(data, tipeSaldoCuti)
				} else {
					max_absen := &Authentication.MaxAbsenIndiv{
						TipeMaxAbsen: dataCuti.TipeMaxAbsen,
					}
					saldo := &Authentication.SaldoIndiv{
						Saldo:           saldoCutiPerTipe.Saldo,
						ValidFrom:       saldoCutiPerTipe.ValidFrom.Format(time.DateOnly),
						ValidTo:         saldoCutiPerTipe.ValidTo.Format(time.DateOnly),
						Periode:         saldoCutiPerTipe.Periode,
						MaxHutang:       saldoCutiPerTipe.MaxHutang,
						ValidFromHutang: saldoCutiPerTipe.ValidFromHutang.Format(time.DateOnly),
					}
					tipeSaldoCuti := Authentication.GetTipeAbsenSaldoIndiv{
						IdTipeAbsen:   dataCuti.IdTipeAbsen,
						NamaTipeAbsen: dataCuti.NamaTipeAbsen,
					}
					if saldo.Periode != "" {
						tipeSaldoCuti.SaldoIndiv = saldo
					}
					if dataCuti.MaxAbsen != nil && *dataCuti.MaxAbsen != 0 {
						max_absen.MaxAbsen = *dataCuti.MaxAbsen
						tipeSaldoCuti.MaxAbsenIndiv = max_absen
					}
					data2 = append(data2, tipeSaldoCuti)
				}
			} else {
				max_absen := &Authentication.MaxAbsenIndiv{
					TipeMaxAbsen: dataCuti.TipeMaxAbsen,
				}
				saldo := &Authentication.SaldoIndiv{
					Saldo:           saldoCutiPerTipe.Saldo,
					ValidFrom:       saldoCutiPerTipe.ValidFrom.Format(time.DateOnly),
					ValidTo:         saldoCutiPerTipe.ValidTo.Format(time.DateOnly),
					Periode:         saldoCutiPerTipe.Periode,
					MaxHutang:       saldoCutiPerTipe.MaxHutang,
					ValidFromHutang: saldoCutiPerTipe.ValidFromHutang.Format(time.DateOnly),
				}
				tipeSaldoCuti := Authentication.GetTipeAbsenSaldoIndiv{
					IdTipeAbsen:   dataCuti.IdTipeAbsen,
					NamaTipeAbsen: dataCuti.NamaTipeAbsen,
				}
				if saldo.Periode != "" {
					tipeSaldoCuti.SaldoIndiv = saldo
				}
				if dataCuti.MaxAbsen != nil && *dataCuti.MaxAbsen != 0 {
					max_absen.MaxAbsen = *dataCuti.MaxAbsen
					tipeSaldoCuti.MaxAbsenIndiv = max_absen
				}

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
		fmt.Println(pihc_mstr_krywn.Company)
		TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenSaldo(pihc_mstr_krywn.Company)

		for _, dataCuti := range TipeAbsen {
			if dataCuti.NamaTipeAbsen == "Cuti Tahunan" {
				TipeAbsenKaryawan := Authentication.GetTipeAbsenKaryawanSaldo{
					IdTipeAbsen:   dataCuti.IdTipeAbsen,
					NamaTipeAbsen: dataCuti.NamaTipeAbsen,
					TipeMaxAbsen:  dataCuti.TipeMaxAbsen,
					CreatedAt:     dataCuti.CreatedAt,
					UpdatedAt:     dataCuti.UpdatedAt,
				}
				if dataCuti.CompCode != nil && *dataCuti.CompCode != "" {
					TipeAbsenKaryawan.CompCode = *dataCuti.CompCode
				}
				if dataCuti.MaxAbsen != nil && *dataCuti.MaxAbsen != 0 {
					TipeAbsenKaryawan.MaxAbsen = *dataCuti.MaxAbsen
				}
				data = append(data, TipeAbsenKaryawan)
			} else {
				TipeAbsenKaryawan := Authentication.GetTipeAbsenKaryawanSaldo{
					IdTipeAbsen:   dataCuti.IdTipeAbsen,
					NamaTipeAbsen: dataCuti.NamaTipeAbsen,
					TipeMaxAbsen:  dataCuti.TipeMaxAbsen,
					CreatedAt:     dataCuti.CreatedAt,
					UpdatedAt:     dataCuti.UpdatedAt,
				}
				if dataCuti.CompCode != nil && *dataCuti.CompCode != "" {
					TipeAbsenKaryawan.CompCode = *dataCuti.CompCode
				}
				if dataCuti.MaxAbsen != nil && *dataCuti.MaxAbsen != 0 {
					TipeAbsenKaryawan.MaxAbsen = *dataCuti.MaxAbsen
				}
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

// Saldo Cuti
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

	kebenaran := false
	var keterangan string
	if req.IDSaldo == 0 {
		sc.TipeAbsenId = req.TipeAbsenId

		// if str, ok := req.Nik.(string); ok {
		// 	sc.Nik = str
		// } else if num, ok := req.Nik.(float64); ok {
		// 	sc.Nik = strconv.Itoa(int(num))
		// }
		sc.Nik = req.Nik
		sc.Saldo = req.Saldo
		sc.ValidFrom, _ = time.Parse(time.DateOnly, req.ValidFrom)
		sc.ValidTo, _ = time.Parse(time.DateOnly, req.ValidTo)

		// if str, ok := req.CreatedBy.(string); ok {
		// 	sc.CreatedBy = str
		// } else if num, ok := req.CreatedBy.(float64); ok {
		// 	createdBy := strconv.Itoa(int(num))
		// 	sc.CreatedBy = createdBy
		// }

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
				CreatedAt:          saldoCuti.CreatedAt,
				UpdatedAt:          saldoCuti.UpdatedAt,
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
				CreatedAt:          saldoCuti.CreatedAt,
				UpdatedAt:          saldoCuti.UpdatedAt,
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

	// var nik string
	// if str, ok := req.Nik.(string); ok {
	// 	nik = str
	// } else if num, ok := req.Nik.(float64); ok {
	// 	nik = strconv.Itoa(int(num))
	// }

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

func (c *CutiKrywnController) GetAdminSaldoCuti(ctx *gin.Context) {
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

func (c *CutiKrywnController) DeleteAdminSaldoCuti(ctx *gin.Context) {
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
