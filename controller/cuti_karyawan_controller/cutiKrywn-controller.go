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
}

func NewCutiKrywnController(Db *gorm.DB) *CutiKrywnController {
	return &CutiKrywnController{PengajuanAbsenRepo: cuti.NewPengajuanAbsenRepo(Db),
		HistoryPengajuanAbsenRepo: cuti.NewHistoryPengajuanAbsenRepo(Db),
		SaldoCutiRepo:             cuti.NewSaldoCutiRepo(Db),
		HistorySaldoCutiRepo:      cuti.NewHistorySaldoCutiRepo(Db),
		TipeAbsenRepo:             cuti.NewTipeAbsenRepo(Db),
		FileAbsenRepo:             cuti.NewFileAbsenRepo(Db),
		TransaksiCutiRepo:         cuti.NewTransaksiCutiRepo(Db),
		PihcMasterKaryDbRepo:      pihc.NewPihcMasterKaryDbRepo(Db)}
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

func (c *CutiKrywnController) StoreCutiKaryawan(ctx *gin.Context) {
	var req Authentication.ValidasiStoreCutiKaryawan
	var sck cuti.PengajuanAbsen

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

	if req.IDPengajuanAbsen == 0 {
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

		tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(*sck.TipeAbsenId)
		existSaldo, saldoCuti, _ := c.SaldoCutiRepo.FindExistSaldo(*sck.TipeAbsenId, sck.Nik, req.MulaiAbsen, req.AkhirAbsen)

		if existSaldo {
			isMax := false
			// Menghitung jumlah hari kerja dan hari kalender
			jumlahHariKerja := 0
			jmlhHariKalender := 0

			for currentDate := sck.MulaiAbsen; !currentDate.After(sck.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
				jmlhHariKalender++
				if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
					jumlahHariKerja++
				}
			}
			fmt.Println("Hari Kalender: ", jmlhHariKalender, ", Hari Kerja:", jumlahHariKerja)

			isMaxPeriode := false
			var indexKerja int
			var indexKalender int
			indexHutang := 0
			var newPeriode time.Time

			if *tipeAbsen.MaxAbsen == 0 {
				for _, dataSaldo := range saldoCuti {
					hariKerja := 0
					hariKalender := 0

					fmt.Println("MASUK")

					if (sck.MulaiAbsen.Before(dataSaldo.ValidTo) || sck.MulaiAbsen.Equal(dataSaldo.ValidTo)) && (sck.AkhirAbsen.After(dataSaldo.ValidTo) || sck.AkhirAbsen.Equal(dataSaldo.ValidTo)) {
						// MulaiAbsen <= ValidTo && AkhirAbsen>=ValidTo
						for currentDate := sck.MulaiAbsen; !currentDate.After(dataSaldo.ValidTo); currentDate = currentDate.AddDate(0, 0, 1) {
							hariKalender++
							if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
								hariKerja++
							}
						}
						newPeriode = dataSaldo.ValidTo
					} else if (sck.MulaiAbsen.After(dataSaldo.ValidFrom) || sck.MulaiAbsen.Equal(dataSaldo.ValidFrom)) && (sck.AkhirAbsen.Before(dataSaldo.ValidTo) || sck.AkhirAbsen.Equal(dataSaldo.ValidTo)) {
						// MulaiAbsen >= ValidFrom && AkhirAbsen<=ValidTo
						for currentDate := sck.MulaiAbsen; !currentDate.After(sck.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
							hariKalender++
							if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
								hariKerja++
							}
						}
					} else if (sck.AkhirAbsen.After(dataSaldo.ValidFrom) || sck.AkhirAbsen.Equal(dataSaldo.ValidFrom)) && (newPeriode.After(dataSaldo.ValidFrom) || newPeriode.Equal(dataSaldo.ValidFrom)) {
						// AkhirAbsen >= ValidFrom && newPeriode>=ValidFrom (periode ke-2)
						for currentDate := newPeriode; !currentDate.After(dataSaldo.ValidFrom); currentDate = currentDate.AddDate(0, 0, 1) {
							hariKalender++
							if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
								hariKerja++
							}
						}
						for currentDate := dataSaldo.ValidFrom; !currentDate.After(sck.AkhirAbsen); currentDate = currentDate.AddDate(0, 0, 1) {
							hariKalender++
							if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
								hariKerja++
							}
						}
					}

					if isMaxPeriode {
						fmt.Println("INDEX: ", indexKerja)
						hariKerja = hariKerja - indexKerja
						hariKalender = hariKalender - indexKalender
						fmt.Println("Kerja: ", hariKerja)
						fmt.Println("Kalender: ", hariKalender)
					}
					fmt.Println(hariKalender, hariKerja, dataSaldo.Saldo)

					// if hariKerja <= *tipeAbsen.MaxAbsen && hariKerja != 0 {
					// 	if hariKerja <= dataSaldo.MaxHutang && dataSaldo.Saldo == 0 && dataSaldo.MaxHutang != 0 {
					// 		isMax = true
					// 		isMaxPeriode = isMax
					// 		dataSaldo.MaxHutang = dataSaldo.MaxHutang - hariKerja
					// 		indexHutang = dataSaldo.MaxHutang
					// 		indexKerja = hariKerja
					// 		indexKalender = hariKalender
					// 	}
					// 	if hariKerja <= dataSaldo.Saldo && dataSaldo.Saldo != 0 {
					// 		if indexHutang != 0 {
					// 			if dataSaldo.Saldo-indexHutang >= 0 {
					// 				isMax = true
					// 				isMaxPeriode = isMax
					// 				dataSaldo.Saldo = dataSaldo.Saldo - hariKerja - indexHutang
					// 			} else {
					// 				if dataSaldo.MaxHutang != 0 {

					// 				}
					// 			}

					// 		} else {
					// 			isMax = true
					// 			isMaxPeriode = isMax
					// 			dataSaldo.Saldo = dataSaldo.Saldo - hariKerja
					// 		}

					// 		indexKerja = hariKerja
					// 		indexKalender = hariKalender
					// 	} else {
					// 		fmt.Println("MASUKKK FALSE")
					// 		isMax = false
					// 	}
					// }
				}
			}

			if isMax {
				fmt.Println(isMax)
			} else {
				fmt.Println(isMax)
			}

			// saldoDigunakan := (sck.AkhirAbsen.Sub(sck.MulaiAbsen).Hours() / 24) + 1
			// fmt.Println(saldoDigunakan, saldoCuti.Saldo)
			// if saldoCuti.Saldo >= int(saldoDigunakan) && saldoCuti.Saldo != 0 {
			// 	saldoCuti.Saldo = saldoCuti.Saldo - int(saldoDigunakan)
			// 	// c.SaldoCutiRepo.Update(saldoCuti)
			// }
			// c.PengajuanAbsenRepo.
		}
	}
	// if err != nil {
	// 	ctx.JSON(http.StatusNotFound, gin.H{
	// 		"status": http.StatusNotFound,
	// 		"info":   "Data Karyawan Tidak Ada",
	// 		"Data":   nil})
	// 	return
	// }

}

// Tipe Saldo Cuti
func (c *CutiKrywnController) GetTipeAbsenSaldoPengajuan(ctx *gin.Context) {
	nik := ctx.Query("nik")
	periode := ctx.Query("tahun")
	data := []Authentication.GetTipeAbsenSaldoIndiv{}
	data2 := []Authentication.GetTipeAbsenSaldoIndiv{}

	if nik == "" {
		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": "Nik wajib di isi"})

		return
	}

	pihc_mstr_krywn, err_pihc_mstr_krywn := c.PihcMasterKaryDbRepo.FindUserByNIK(nik)

	if err_pihc_mstr_krywn == nil {
		fmt.Println(pihc_mstr_krywn.Company)
		TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsen(pihc_mstr_krywn.Company)

		for _, dataCuti := range TipeAbsen {
			saldoCutiPerTipe, _ := c.SaldoCutiRepo.GetSaldoCutiPerTipe(dataCuti.IdTipeAbsen, pihc_mstr_krywn.EmpNo, periode)
			for _, saldoCuti := range saldoCutiPerTipe {
				if dataCuti.NamaTipeAbsen == "Cuti Tahunan" {
					tipeSaldoCuti := Authentication.GetTipeAbsenSaldoIndiv{
						IdTipeAbsen:   dataCuti.IdTipeAbsen,
						NamaTipeAbsen: dataCuti.NamaTipeAbsen,
						Saldo:         saldoCuti.Saldo,
						ValidFrom:     saldoCuti.ValidFrom.Format(time.DateOnly),
						ValidTo:       saldoCuti.ValidTo.Format(time.DateOnly),
					}
					data = append(data, tipeSaldoCuti)
				} else {
					tipeSaldoCuti := Authentication.GetTipeAbsenSaldoIndiv{
						IdTipeAbsen:   dataCuti.IdTipeAbsen,
						NamaTipeAbsen: dataCuti.NamaTipeAbsen,
						Saldo:         saldoCuti.Saldo,
						ValidFrom:     saldoCuti.ValidFrom.Format(time.DateOnly),
						ValidTo:       saldoCuti.ValidTo.Format(time.DateOnly),
					}
					data2 = append(data2, tipeSaldoCuti)
				}
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
	data := []Authentication.GetTipeAbsenKaryawan{}
	data2 := []Authentication.GetTipeAbsenKaryawan{}

	if nik == "" {
		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": "Nik wajib di isi"})

		return
	}

	pihc_mstr_krywn, err_pihc_mstr_krywn := c.PihcMasterKaryDbRepo.FindUserByNIK(nik)

	if err_pihc_mstr_krywn == nil {
		fmt.Println(pihc_mstr_krywn.Company)
		TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsen(pihc_mstr_krywn.Company)

		for _, dataCuti := range TipeAbsen {
			if dataCuti.NamaTipeAbsen == "Cuti Tahunan" {
				TipeAbsenKaryawan := Authentication.GetTipeAbsenKaryawan{
					IdTipeAbsen:   dataCuti.IdTipeAbsen,
					NamaTipeAbsen: dataCuti.NamaTipeAbsen,
				}
				data = append(data, TipeAbsenKaryawan)
			} else {
				TipeAbsenKaryawan := Authentication.GetTipeAbsenKaryawan{
					IdTipeAbsen:   dataCuti.IdTipeAbsen,
					NamaTipeAbsen: dataCuti.NamaTipeAbsen,
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
			TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(dataSaldoo.TipeAbsenId)
			dataSaldoCuti := Authentication.ListSaldoCutiKaryawan{
				IdSaldoCuti:     dataSaldoo.IdSaldoCuti,
				TipeAbsenId:     dataSaldoo.TipeAbsenId,
				NamaTipeAbsen:   TipeAbsen.NamaTipeAbsen,
				Nik:             dataSaldoo.Nik,
				Saldo:           dataSaldoo.Saldo,
				ValidFrom:       dataSaldoo.ValidFrom.Format(time.DateOnly),
				ValidTo:         dataSaldoo.ValidTo.Format(time.DateOnly),
				CreatedBy:       dataSaldoo.CreatedBy,
				CreatedAt:       dataSaldoo.CreatedAt,
				UpdatedAt:       dataSaldoo.UpdatedAt,
				Periode:         dataSaldoo.Periode,
				MaxHutang:       dataSaldoo.MaxHutang,
				ValidFromHutang: dataSaldoo.ValidFromHutang.Format(time.DateOnly),
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
		TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(saldoCuti.TipeAbsenId)
		dataSaldoCuti := Authentication.ListSaldoCutiKaryawan{
			IdSaldoCuti:     saldoCuti.IdSaldoCuti,
			TipeAbsenId:     saldoCuti.TipeAbsenId,
			NamaTipeAbsen:   TipeAbsen.NamaTipeAbsen,
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
