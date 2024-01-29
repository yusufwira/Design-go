package cuti_karyawan_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	PihcMasterKaryRtDbRepo    *pihc.PihcMasterKaryRtDbRepo
	PihcMasterCompanyRepo     *pihc.PihcMasterCompanyRepo
}

func NewCutiKrywnController(Db *gorm.DB) *CutiKrywnController {
	return &CutiKrywnController{
		PengajuanAbsenRepo:        cuti.NewPengajuanAbsenRepo(Db),
		HistoryPengajuanAbsenRepo: cuti.NewHistoryPengajuanAbsenRepo(Db),
		SaldoCutiRepo:             cuti.NewSaldoCutiRepo(Db),
		HistorySaldoCutiRepo:      cuti.NewHistorySaldoCutiRepo(Db),
		TipeAbsenRepo:             cuti.NewTipeAbsenRepo(Db),
		FileAbsenRepo:             cuti.NewFileAbsenRepo(Db),
		TransaksiCutiRepo:         cuti.NewTransaksiCutiRepo(Db),
		PihcMasterKaryRtDbRepo:    pihc.NewPihcMasterKaryRtDbRepo(Db),
		PihcMasterCompanyRepo:     pihc.NewPihcMasterCompanyRepo(Db)}
}

// Pengajuan Cuti (DONE)
func (c *CutiKrywnController) StoreCutiKaryawan(ctx *gin.Context) {
	var req Authentication.ValidasiStoreCutiKaryawan
	var sck cuti.PengajuanAbsen
	var fsc []cuti.FileAbsen
	var trsc []Authentication.SaldoCutiTransaksiPengajuan
	transaksi_cuti := cuti.TransaksiCuti{}

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

	PIHC_MSTR_KRY, _ := c.PihcMasterKaryRtDbRepo.FindUserByNIK(req.Nik)
	comp_code := PIHC_MSTR_KRY.Company

	sck.Nik = req.Nik
	sck.TipeAbsenId = &req.TipeAbsenId
	sck.CompCode = comp_code
	sck.Deskripsi = &req.Deskripsi
	sck.MulaiAbsen, _ = time.Parse(time.DateOnly, req.MulaiAbsen)
	sck.AkhirAbsen, _ = time.Parse(time.DateOnly, req.AkhirAbsen)
	sck.TglPengajuan, _ = time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	sck.Status = &req.Status
	periode := strconv.Itoa(time.Now().Year())
	sck.Periode = &periode
	sck.CreatedBy = &req.CreatedBy

	// Mencari Atasan
	url := "https://api-pismart-dev.pupuk-indonesia.com/oauth_api/api/delegasi/tryFindSuperios/" + sck.Nik
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIzMTIiLCJqdGkiOiI0NGIxOWE2NzFlNjU4OGYwMzZkZjQzZjBlYWVkNDc5NzY2NDM4YzgxNDM2NmI3YmRiZmM2N2NiYzk3MTQ1MmI4ZjBkY2Y4NDk3OWMzZThhNCIsImlhdCI6MTcwNDc4MTA1MC4xODM3NjEsIm5iZiI6MTcwNDc4MTA1MC4xODM3NywiZXhwIjoxNzM2NDAzNDUwLjE3NDMsInN1YiI6IjU0OCIsInNjb3BlcyI6W119.cTNgAPfuW3Cxg-wAliHGxb8zZOLqq5ym6n2SM6yLHv431WRVWC_YFsw3eQ673hiK63GqDSs6hbDHJEzZxmhpxWEIaXRl6L_Zg9htPSqykJKG5uY9MoHbfs9BXxVKLQMntzfwc6K5S48incJFhWLVKVOopLuPQJ2e9-m25E9f2-mE_AfeHr0KvQhrmUHkB-T-GxDFBqn5YP9xQU4DKpDb-baPAUY_7wlDcTgSAJrGuy8LgLT96wNnDJA9rGXvTDJ7JRpF2SrsZEwIkMk9RF3MWPE4SLOnWC8xMkrM6WK4sTHEloBhfzH4_CMwegsIFN7XfHCavBmTU4PdSbg1dc_8tHJZ0zJALlZZ44UhdlfsMT1uGNcJAam_M7lo9Qnn3knl06pe3gWkm7uo2B7262Dvs-jHP5gLmdRmFZJcvyPNGBdjfdwvOwW3hALzGeDwypnHI-UXS5lQQvAC0kn0PTvmGWQl16Z7JK9He9B75fpn6Cb0StGq_xy9vL6MQ0QkinnmzwHfQuNwOWolz-2LvRFjz7MKUSged5q6r0sz0N-62daNIgGB3GIItS9taLmGfSbbghPuB3y8wHm5yqgadCmyGe2x4OuuzWpnRBxDXFVqnLo-iph90Squ0RgIxhGt5dSN94z9lA3vqUJbt7Em88cjmLZ8po8XAEFKBtLua2CQ6y8"

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading the response bytes:", err)
	}

	var data_atasan pihc.PihcMasterKaryRt
	json.Unmarshal(body, &data_atasan)

	foto := "https://storage.googleapis.com/lumen-oauth-storage/DataKaryawan/Foto/" + data_atasan.Company + "/" + data_atasan.EmpNo + ".jpg"
	respons, err := http.Get(foto)
	if err != nil || respons.StatusCode != http.StatusOK {
		default_foto := "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
		foto = default_foto
	}

	statusApproved1 := "Approver 1"
	var status1 *string
	if req.Status == "Submitted" {
		status1 = new(string)
		*status1 = "WaitApv"
	}

	temp := cuti.AtasanApproved{
		Nik:          data_atasan.EmpNo,
		Name:         data_atasan.Nama,
		Position:     data_atasan.PosTitle,
		TypeApprover: &statusApproved1,
		Status:       status1,
		Photo:        foto,
	}

	var atasan []cuti.AtasanApproved
	atasan = append(atasan, temp)
	// Marshal atasan into JSON
	jsonDataAtasan, _ := json.Marshal(atasan)
	sck.ApprovedBy = json.RawMessage(jsonDataAtasan)

	// dataKaryawan, _ := c.PihcMasterKaryRtDbRepo.FindUserByNIK(sck.Nik)
	// if dataKaryawan.PosTitle != "Wakil Direktur Utama" {
	// 	for dataKaryawan.PosTitle != "Wakil Direktur Utama" {
	// 		dataKaryawan, _ = c.PihcMasterKaryRtDbRepo.FindUserAtasanBySupPosID(dataKaryawan.SupPosID)
	// 		if dataKaryawan.SupPosID == "" {
	// 			break
	// 		}
	// 	}
	// } else {
	// 	for dataKaryawan.PosTitle != "Direktur Utama" {
	// 		dataKaryawan, _ = c.PihcMasterKaryRtDbRepo.FindUserAtasanBySupPosID(dataKaryawan.SupPosID)
	// 		if dataKaryawan.SupPosID == "" {
	// 			break
	// 		}
	// 	}
	// }
	// approvedBy := dataKaryawan.EmpNo
	// if approvedBy != "" {
	// 	sck.ApprovedBy = &approvedBy
	// }

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
	checkMaxAbsen := false

	if req.IdPengajuanAbsen == nil {
		if tipeAbsen.MaxAbsen != nil {
			JmlHariKerja := 0
			jmlhHariKalender := 0

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
			transaksi_cuti.TipeAbsenId = tipeAbsen.IdTipeAbsen

			sck.JmlHariKalendar = &jmlhHariKalender
			sck.JmlHariKerja = &JmlHariKerja

			create = true
			checkMaxAbsen = true
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
		if (*pengajuan_absen.Status == "Submitted") || (*pengajuan_absen.Status == "Drafted") || (*pengajuan_absen.Status == "Rejected") {
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
				transaksi_cuti.TipeAbsenId = tipeAbsen.IdTipeAbsen

				sck.JmlHariKalendar = &jmlhHariKalender
				sck.JmlHariKerja = &JmlHariKerja

				checkMaxAbsen = true
				update = true
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
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	if checkSaldo || checkMaxAbsen {
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

			// Transaksi Cuti
			if len(trsc) != 0 {
				for _, transaction := range trsc {
					transaksi_cuti := cuti.TransaksiCuti{}
					transaksi_cuti.PengajuanAbsenId = sckData.IdPengajuanAbsen
					transaksi_cuti.TipeAbsenId = *sckData.TipeAbsenId
					transaksi_cuti.Nik = sckData.Nik
					transaksi_cuti.Periode = transaction.Periode
					transaksi_cuti.JumlahCuti = transaction.JmlhCuti
					transaksi_cuti.TipeHari = tipe_hari
					c.TransaksiCutiRepo.Create(transaksi_cuti)
				}
			} else {
				if sckData.Periode != nil {
					transaksi_cuti.Periode = *sckData.Periode
				}
				transaksi_cuti.PengajuanAbsenId = sckData.IdPengajuanAbsen
				transaksi_cuti.TipeAbsenId = *sckData.TipeAbsenId
				transaksi_cuti.Nik = sckData.Nik
				if sckData.Periode != nil {
					transaksi_cuti.Periode = *sckData.Periode
				}
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

				// Transaksi Cuti
				if len(trsc) != 0 {
					for _, transaction := range trsc {
						transaksi_cuti := cuti.TransaksiCuti{}
						transaksi_cuti.PengajuanAbsenId = sckData.IdPengajuanAbsen
						transaksi_cuti.TipeAbsenId = *sckData.TipeAbsenId
						transaksi_cuti.Nik = sckData.Nik
						transaksi_cuti.Periode = transaction.Periode
						transaksi_cuti.JumlahCuti = transaction.JmlhCuti
						transaksi_cuti.TipeHari = tipe_hari
						c.TransaksiCutiRepo.Create(transaksi_cuti)
					}
				} else {
					transaksi_cuti.PengajuanAbsenId = sckData.IdPengajuanAbsen
					transaksi_cuti.TipeAbsenId = *sckData.TipeAbsenId
					transaksi_cuti.Nik = sckData.Nik
					if sckData.Periode != nil {
						transaksi_cuti.Periode = *sckData.Periode
					}

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
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"status":     http.StatusServiceUnavailable,
				"keterangan": keterangan_x,
			})
		} else {
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
		if *pengajuanAbsen.Status == "Submitted" || *pengajuanAbsen.Status == "Drafted" || *pengajuanAbsen.Status == "Rejected" {
			c.PengajuanAbsenRepo.DelPengajuanCuti(pengajuanAbsen.IdPengajuanAbsen)
			transaksi_cuti, _ := c.TransaksiCutiRepo.FindDataTransaksiIDPengajuan(pengajuanAbsen.IdPengajuanAbsen)
			file_cuti, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuan(pengajuanAbsen.IdPengajuanAbsen)

			for _, tr := range transaksi_cuti {
				c.TransaksiCutiRepo.Delete(tr)
			}
			for _, fc := range file_cuti {
				c.FileAbsenRepo.Delete(fc)
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
		if *pengajuan_absen.Status == "Submitted" {
			var approved []cuti.AtasanApproved
			json.Unmarshal(pengajuan_absen.ApprovedBy, &approved)

			var jumlahApprove = 0
			var jumlahApproval = 0
			for i, atasan := range approved {
				if *atasan.Status != "WaitApv" {
					continue
				}

				if req.IsManager {
					if atasan.Nik == req.Nik {
						if req.Status == "Approved" {
							approved[i].Status = &req.Status
							if i+1 < len(approved) {
								approved[i+1].Status = new(string)
								*approved[i+1].Status = "WaitApv"
							}
						} else if req.Status == "Rejected" {
							approved[i].Status = new(string)
							approved[i].Status = &req.Status
							approved[i].Keterangan = &req.Keterangan
							pengajuan_absen.Status = &req.Status
							pengajuan_absen.Keterangan = &req.Keterangan
						}
					}
					break
				} else {
					if req.Status == "Approved" {
						approved[i].Status = &req.Status
						if i+1 < len(approved) {
							approved[i+1].Status = new(string)
							*approved[i+1].Status = "WaitApv"
						}
					} else if req.Status == "Rejected" {
						approved[i].Status = new(string)
						approved[i].Status = &req.Status
						approved[i].Keterangan = &req.Keterangan
						pengajuan_absen.Status = &req.Status
						pengajuan_absen.Keterangan = &req.Keterangan
					}
					break
				}
			}

			for _, atasan := range approved {
				if atasan.Status != nil && *atasan.Status == "Approved" {
					jumlahApprove++
				}
				jumlahApproval++
			}

			if jumlahApprove == jumlahApproval {
				pengajuan_absen.Status = new(string)
				*pengajuan_absen.Status = "Completed"
			}

			if *pengajuan_absen.Status == "Rejected" {
				approvedBytes, _ := json.Marshal(approved)
				pengajuan_absen.ApprovedBy = json.RawMessage(approvedBytes)

				updated_pengajuan, _ := c.PengajuanAbsenRepo.Update(pengajuan_absen)
				history_pengajuan := HistoryPengajuanCutiSet(updated_pengajuan)
				c.HistoryPengajuanAbsenRepo.Create(history_pengajuan)

				ctx.JSON(http.StatusOK, gin.H{
					"status":     http.StatusOK,
					"keterangan": "Success",
				})
			} else if *pengajuan_absen.Status == "Completed" {
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
					approvedBytes, _ := json.Marshal(approved)
					pengajuan_absen.ApprovedBy = json.RawMessage(approvedBytes)

					pengajuan_absen.Keterangan = new(string)
					pengajuan_absen.Keterangan = nil
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
					pengajuan_absen.Keterangan = new(string)
					pengajuan_absen.Keterangan = nil
					*pengajuan_absen.Status = "Submitted"
					updated_pengajuan, _ := c.PengajuanAbsenRepo.Update(pengajuan_absen)
					history_pengajuan := HistoryPengajuanCutiSet(updated_pengajuan)
					c.HistoryPengajuanAbsenRepo.Create(history_pengajuan)

					ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
						"status":     http.StatusServiceUnavailable,
						"keterangan": "Maaf Pengajuan tidak dapat di Approve, Saldo Tidak Cukup, Silahkan Mengajukan kembali",
					})
				}
			} else if *pengajuan_absen.Status == "Submitted" {
				approvedBytes, _ := json.Marshal(approved)
				pengajuan_absen.ApprovedBy = json.RawMessage(approvedBytes)

				pengajuan_absen.Keterangan = new(string)
				pengajuan_absen.Keterangan = nil
				updated_pengajuan, _ := c.PengajuanAbsenRepo.Update(pengajuan_absen)
				history_pengajuan := HistoryPengajuanCutiSet(updated_pengajuan)
				c.HistoryPengajuanAbsenRepo.Create(history_pengajuan)

				for _, saldo := range eksekusi_saldo {
					updated_saldo, _ := c.SaldoCutiRepo.Update(saldo)
					history_saldo := HistorySaldoCutiSet(updated_saldo)
					hsc, errr := c.HistorySaldoCutiRepo.Create(history_saldo)
					if errr == nil {
						fmt.Println(hsc)
					}
				}
				ctx.JSON(http.StatusOK, gin.H{
					"status":     http.StatusOK,
					"keterangan": "Success",
				})
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

	fmt.Println(req.IsManager, req.Status)

	var dataDB []cuti.PengajuanAbsen
	var err error

	if req.Status == "Submitted" {
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
		karyawan, _ := c.PihcMasterKaryRtDbRepo.FindUserByNIKArray(arrNIK)
		tipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByIDArray(arrTipeAbsenID)
		files, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuanArray(arrIdPengajuanAbsen)
		for _, myKrywn := range karyawan {
			arrCompany = append(arrCompany, myKrywn.Company)
		}
		companys, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompanyArray(arrCompany)

		for _, myCuti := range dataDB {
			myFiles := []cuti.FileAbsen{}
			list_pengajuan := Authentication.ListApprovalCuti{}
			// fmt.Printf("Tipe data dari integerVariable: %T\n", myCuti.ApprovedBy)

			var approved []cuti.AtasanApproved
			json.Unmarshal(myCuti.ApprovedBy, &approved)

			list_pengajuan.ApprovedBy = approved
			// Karyawan
			for _, myKaryawan := range karyawan {
				if myCuti.Nik == myKaryawan.EmpNo {
					for _, myCompany := range companys {
						if myKaryawan.Company == myCompany.Code {
							list_pengajuan.PihcMasterKaryRt = convertSourceTargetDataKaryawan(myKaryawan)
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
		karyawan, _ := c.PihcMasterKaryRtDbRepo.FindUserByNIK(dataDB.Nik)
		companys, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(karyawan.Company)
		files, _ := c.FileAbsenRepo.FindFileAbsenByIDPengajuan(dataDB.IdPengajuanAbsen)
		if files == nil {
			files = []cuti.FileAbsen{}
		}

		data_karyawan_convert := convertSourceTargetDataKaryawan(karyawan)
		result := convertSourceTargetMyPengajuanAbsen(dataDB, tipeAbsen)

		list_aprvl.IdPengajuanAbsen = result.IdPengajuanAbsen
		list_aprvl.PihcMasterKaryRt = data_karyawan_convert
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

		list_pengajuan := Authentication.ListApprovalCuti{}
		var approved []cuti.AtasanApproved
		json.Unmarshal(result.ApprovedBy, &approved)

		list_pengajuan.ApprovedBy = approved

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

	pihc_mstr_krywn, err_pihc_mstr_krywn := c.PihcMasterKaryRtDbRepo.FindUserByNIK(req.NIK)

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

	pihc_mstr_krywn, err_pihc_mstr_krywn := c.PihcMasterKaryRtDbRepo.FindUserByNIK(nik)

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
	var req1 Authentication.ValidasiStoreSaldoCuti
	var req2 []Authentication.ValidasiStoreSaldoCuti
	var data []Authentication.SaldoCutiKaryawan
	var sc cuti.SaldoCuti

	if err := ctx.ShouldBindBodyWith(&req1, binding.JSON); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]Authentication.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = Authentication.ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": out})
		}
		if err := ctx.ShouldBindBodyWith(&req2, binding.JSON); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				out := make([]Authentication.ErrorMsg, len(ve))
				for i, fe := range ve {
					out[i] = Authentication.ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
				}
				ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": out})
			}
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"errorcode_": http.StatusServiceUnavailable, "errormsg_": "Pastikan field tidak kosong"})
			return
		}
	}

	if len(req2) == 0 {
		req2 = append(req2, req1)
	}

	kebenaran := false
	var keterangan string
	for _, temp := range req2 {
		if temp.IDSaldo != nil {
			temp.IDSaldo = ConvertInterfaceTypeDataToInt(temp.IDSaldo)
		}

		if temp.IDSaldo == nil {
			sc.TipeAbsenId = temp.TipeAbsenId
			sc.Nik = temp.Nik
			sc.Saldo = temp.Saldo
			sc.ValidFrom, _ = time.Parse(time.DateOnly, temp.ValidFrom)
			sc.ValidTo, _ = time.Parse(time.DateOnly, temp.ValidTo)
			sc.CreatedBy = temp.CreatedBy

			if temp.Periode != "" {
				sc.Periode = temp.Periode
			}

			sc.MaxHutang = temp.MaxHutang
			sc.ValidFromHutang, _ = time.Parse(time.DateOnly, temp.ValidFromHutang)

			saldoCuti, err := c.SaldoCutiRepo.Create(sc)

			if err == nil {
				data = append(data, SaldoCutiKaryawanSet(saldoCuti))
				dataHistorySaldoCuti := HistorySaldoCutiSet(saldoCuti)
				c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)

				kebenaran = true
				keterangan = "Success"
			} else {
				data = []Authentication.SaldoCutiKaryawan{}

				kebenaran = false
				keterangan = "Gagal Store Saldo Cuti"
			}
		} else {
			sc, _ := c.SaldoCutiRepo.GetSaldoCutiByID(temp.IDSaldo)
			sc.Saldo = temp.Saldo
			sc.ValidFrom, _ = time.Parse(time.DateOnly, temp.ValidFrom)
			sc.ValidTo, _ = time.Parse(time.DateOnly, temp.ValidTo)
			sc.CreatedBy = temp.CreatedBy

			if temp.Periode != "" {
				sc.Periode = temp.Periode
			}
			// sc.Periode = strconv.Itoa(sc.ValidTo.Year())
			sc.MaxHutang = temp.MaxHutang
			sc.ValidFromHutang, _ = time.Parse(time.DateOnly, temp.ValidFromHutang)

			saldoCuti, err := c.SaldoCutiRepo.Update(sc)
			if err == nil {
				data = append(data, SaldoCutiKaryawanSet(saldoCuti))
				dataHistorySaldoCuti := HistorySaldoCutiSet(saldoCuti)
				c.HistorySaldoCutiRepo.Create(dataHistorySaldoCuti)

				kebenaran = true
				keterangan = "Success"
			} else {
				data = []Authentication.SaldoCutiKaryawan{}

				kebenaran = false
				keterangan = "Gagal Update Saldo Cuti"
			}
		}
	}

	var output interface{}
	if len(data) == 1 {
		output = data[0]
	} else {
		output = data
	}
	if kebenaran {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"Success": keterangan,
			"data":    output,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"status":     http.StatusServiceUnavailable,
			"keterangan": keterangan,
			"data":       output,
		})
	}
}
func (c *CutiKrywnController) ListAdminSaldoCutiKaryawan(ctx *gin.Context) {
	var req Authentication.ValidasiListSaldoCuti
	fmt.Println(req.Perusahaan, req.Kompartemen, req.Departemen, req.Direktorat)
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

	saldoCuti, err := c.SaldoCutiRepo.FindSaldoCutiKaryawanAdmin(req.Key, req.Perusahaan, req.Direktorat, req.Departemen, req.Kompartemen, req.Nik, req.Tahun)

	if err == nil {
		for _, dataSaldoo := range saldoCuti {
			karyawan, _ := c.PihcMasterKaryRtDbRepo.FindUserByNIK(dataSaldoo.Nik)
			company, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(karyawan.Company)
			TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(dataSaldoo.TipeAbsenId)
			dataSaldoCuti := ListSaldoCutiKaryawanSet(dataSaldoo, company, TipeAbsen)

			if karyawan.Nama != "" {
				dataSaldoCuti.Nama = karyawan.Nama
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
		karyawan, _ := c.PihcMasterKaryRtDbRepo.FindUserByNIK(saldoCuti.Nik)
		company, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(karyawan.Company)
		TipeAbsen, _ := c.TipeAbsenRepo.FindTipeAbsenByID(saldoCuti.TipeAbsenId)
		data = ListSaldoCutiKaryawanSet(saldoCuti, company, TipeAbsen)
		if karyawan.Nama != "" {
			data.Nama = karyawan.Nama
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

func (c *CutiKrywnController) GetCompany(ctx *gin.Context) {
	company, _ := c.PihcMasterCompanyRepo.FindAllCompany()
	if company != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    company,
			"status":  http.StatusOK,
			"success": "Success",
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":       company,
			"status":     http.StatusBadRequest,
			"keterangan": "Data Company tidak Ditemukan!!",
		})
	}
}
func (c *CutiKrywnController) GetDirektorat(ctx *gin.Context) {
	var req Authentication.ValidationCompany

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

	direktorat, _ := c.PihcMasterKaryRtDbRepo.FindDirektoratCompany(req.Company)
	if direktorat != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    direktorat,
			"status":  http.StatusOK,
			"success": "Success",
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":       direktorat,
			"status":     http.StatusBadRequest,
			"keterangan": "Data Direktorat tidak Ditemukan!!",
		})
	}
}
func (c *CutiKrywnController) GetKompartemen(ctx *gin.Context) {
	var req Authentication.ValidationKompartemen

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

	kompartemen, _ := c.PihcMasterKaryRtDbRepo.FindKompartemenCompany(req.Company, req.Direktorat)
	if kompartemen != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    kompartemen,
			"status":  http.StatusOK,
			"success": "Success",
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":       kompartemen,
			"status":     http.StatusBadRequest,
			"keterangan": "Data Kompartemen tidak Ditemukan!!",
		})
	}
}
func (c *CutiKrywnController) GetDepartemen(ctx *gin.Context) {
	var req Authentication.ValidationDepartemen

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

	departemen, _ := c.PihcMasterKaryRtDbRepo.FindDepartemenCompany(req.Company, req.Kompartemen)
	if departemen != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    departemen,
			"status":  http.StatusOK,
			"success": "Success",
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data":       departemen,
			"status":     http.StatusBadRequest,
			"keterangan": "Data Departemen tidak Ditemukan!!",
		})
	}
}
