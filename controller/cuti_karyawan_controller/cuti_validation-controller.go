package cuti_karyawan_controller

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/cuti"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
)

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return (fe.Field() + " wajib di isi")
	}
	return "Unknown error"
}

// Perhitungan jmlhHariKalender && JmlHariKerja
func perhitungan(mulai time.Time, akhir time.Time) [2]int {
	jmlhHariKalender := 0
	JmlHariKerja := 0
	for currentDate := mulai; !currentDate.After(akhir); currentDate = currentDate.AddDate(0, 0, 1) {
		jmlhHariKalender++
		if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
			JmlHariKerja++
		}
	}
	return [2]int{jmlhHariKalender, JmlHariKerja}
}

// Convert Db to struct
func convertSourceTargetDataKaryawan(source pihc.PihcMasterKaryRtDb) pihc.PihcMasterKaryRt {
	return pihc.PihcMasterKaryRt{
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
		Lokasi:         source.Lokasi,
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
func HistorySaldoCutiSet(source cuti.SaldoCuti) cuti.HistorySaldoCuti {
	return cuti.HistorySaldoCuti{
		IdHistorySaldoCuti: source.IdSaldoCuti,
		TipeAbsenId:        source.TipeAbsenId,
		Nik:                source.Nik,
		Saldo:              source.Saldo,
		ValidFrom:          source.ValidFrom,
		ValidTo:            source.ValidTo,
		CreatedBy:          source.CreatedBy,
		Periode:            source.Periode,
		MaxHutang:          source.MaxHutang,
		ValidFromHutang:    source.ValidFromHutang,
	}
}
func HistoryPengajuanCutiSet(source cuti.PengajuanAbsen) cuti.HistoryPengajuanAbsen {
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
func SaldoCutiKaryawanSet(source cuti.SaldoCuti) Authentication.SaldoCutiKaryawan {
	return Authentication.SaldoCutiKaryawan{
		IdSaldoCuti:     source.IdSaldoCuti,
		TipeAbsenId:     source.TipeAbsenId,
		Nik:             source.Nik,
		Saldo:           source.Saldo,
		ValidFrom:       source.ValidFrom.Format(time.DateOnly),
		ValidTo:         source.ValidTo.Format(time.DateOnly),
		CreatedBy:       source.CreatedBy,
		CreatedAt:       source.CreatedAt,
		UpdatedAt:       source.UpdatedAt,
		Periode:         source.Periode,
		MaxHutang:       source.MaxHutang,
		ValidFromHutang: source.ValidFromHutang.Format(time.DateOnly),
	}
}
func GetTipeAbsenKaryawanSaldoSet(source cuti.TipeAbsen) Authentication.GetTipeAbsenKaryawanSaldo {
	return Authentication.GetTipeAbsenKaryawanSaldo{
		IdTipeAbsen:   source.IdTipeAbsen,
		NamaTipeAbsen: source.NamaTipeAbsen,
		MaxAbsen:      source.MaxAbsen,
		TipeMaxAbsen:  source.TipeMaxAbsen,
		CompCode:      source.CompCode,
		CreatedAt:     source.CreatedAt,
		UpdatedAt:     source.UpdatedAt}
}
func CompanyKaryawanSet(source pihc.PihcMasterCompany) Authentication.CompanyKaryawan {
	return Authentication.CompanyKaryawan{
		Code: source.Code,
		Name: source.Name,
	}
}
func ListSaldoCutiKaryawanSet(source1 cuti.SaldoCuti, source2 pihc.PihcMasterCompany, source3 cuti.TipeAbsen) Authentication.ListSaldoCutiKaryawan {
	return Authentication.ListSaldoCutiKaryawan{
		IdSaldoCuti:               source1.IdSaldoCuti,
		GetTipeAbsenKaryawanSaldo: GetTipeAbsenKaryawanSaldoSet(source3),
		Nik:                       source1.Nik,
		CompanyKaryawan:           CompanyKaryawanSet(source2),
		Saldo:                     source1.Saldo,
		ValidFrom:                 source1.ValidFrom.Format(time.DateOnly),
		ValidTo:                   source1.ValidTo.Format(time.DateOnly),
		CreatedBy:                 source1.CreatedBy,
		CreatedAt:                 source1.CreatedAt,
		UpdatedAt:                 source1.UpdatedAt,
		Periode:                   source1.Periode,
		MaxHutang:                 source1.MaxHutang,
		ValidFromHutang:           source1.ValidFromHutang.Format(time.DateOnly),
	}
}

// Interface Type Data
func ConvertInterfaceTypeDataToInt(x interface{}) (result int) {
	dataType := reflect.TypeOf(x)
	switch dataType.Kind() {
	case reflect.String:
		// If it's a string, convert to int
		if strValue, ok := x.(string); ok {
			intValue, err := strconv.Atoi(strValue)
			if err != nil {
				// Handle error if conversion fails
				fmt.Println("Error converting string to int:", err)
				return
			}
			result = intValue
		}
	case reflect.Float64:
		// If it's a float64, convert to int
		if floatValue, ok := x.(float64); ok {
			result = int(floatValue)
		}
	}
	return result
}
