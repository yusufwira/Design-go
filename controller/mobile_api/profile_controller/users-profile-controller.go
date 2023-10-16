package profile_controller

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/mobile/profile"
	"github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type UsersProfileController struct {
	UserProfileRepo          *users.UserProfileRepo
	ProfileRepo              *profile.ProfileRepo
	AboutUsRepo              *profile.AboutUsRepo
	ProfileSkillRepo         *profile.ProfileSkillRepo
	PengalamanKerjaRepo      *profile.PengalamanKerjaRepo
	PhotoProfileRepo         *profile.PhotoProfileRepo
	PihcKaryawanMutasiPIRepo *pihc.PihcKaryawanMutasiPIRepo
	PihcMasterKaryDbRepo     *pihc.PihcMasterKaryDbRepo
	PihcMasterCompanyRepo    *pihc.PihcMasterCompanyRepo
	ViewOrganisasiRepo       *pihc.ViewOrganisasiRepo
}

func NewUsersProfileController(Db *gorm.DB, StorageClient *storage.Client) *UsersProfileController {
	return &UsersProfileController{UserProfileRepo: users.NewUserProfileRepo(Db),
		ProfileRepo:              profile.NewProfileRepo(Db),
		AboutUsRepo:              profile.NewAboutUsRepo(Db),
		ProfileSkillRepo:         profile.NewProfileSkillRepo(Db),
		PengalamanKerjaRepo:      profile.NewPengalamanKerjaRepo(Db),
		PhotoProfileRepo:         profile.NewPhotoProfileRepo(Db, StorageClient),
		PihcKaryawanMutasiPIRepo: pihc.NewPihcKaryawanMutasiPIRepo(Db),
		PihcMasterKaryDbRepo:     pihc.NewPihcMasterKaryDbRepo(Db),
		PihcMasterCompanyRepo:    pihc.NewPihcMasterCompanyRepo(Db),
		ViewOrganisasiRepo:       pihc.NewViewOrganisasiRepo(Db)}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return (fe.Field() + " wajib di isi")
	case "validyear":
		return ("Field has an invalid value: " + fe.Field() + fe.Tag())
	}
	return "Unknown error"
}

func (c *UsersProfileController) StoreData(ctx *gin.Context) {
	var req Authentication.ValidationSavePersonalInformationEmployee
	var data Authentication.PersonalInformationEmployee

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

	personalInformation, err := c.UserProfileRepo.FindProfileUsers(req.EmployeeId)
	personalInformation.Nik = req.EmployeeId
	personalInformation.Alamat = &req.Alamat
	personalInformation.Kelurahan = &req.Kelurahan
	personalInformation.Kecamatan = &req.Kecamatan
	personalInformation.Kabupaten = &req.Kabupaten
	personalInformation.Provinsi = &req.Provinsi
	personalInformation.Kodepos = &req.KodePos
	personalInformation.TipeDomisili = &req.TipeDomisili
	domisili := c.UserProfileRepo.FindKetDomisili(req.TipeDomisili)
	personalInformation.KetDomisili = domisili.KetDomisili
	personalInformation.PosisiMap = &req.PosisiMap
	personalInformation.Lat = &req.Lat
	personalInformation.Long = &req.Long

	if err != nil {
		result, _ := c.UserProfileRepo.Create(personalInformation)

		data.Nik = result.Nik
		data.Alamat = result.Alamat
		data.Kelurahan = result.Kelurahan
		data.Kecamatan = result.Kecamatan
		data.Kabupaten = result.Kabupaten
		data.Provinsi = result.Provinsi
		data.KodePos = result.Kodepos
		data.TipeDomisili = result.TipeDomisili
		data.KetDomisili = result.KetDomisili
		data.PosisiMap = result.PosisiMap
		data.Lat = result.Lat
		data.Long = result.Long
		data.UpdatedFrom = result.UpdatedFrom
		data.UpdatedDate = result.UpdatedDate.Format(time.DateTime)
	} else {
		result, _ := c.UserProfileRepo.Update(personalInformation)

		data.Nik = result.Nik
		data.Alamat = result.Alamat
		data.Kelurahan = result.Kelurahan
		data.Kecamatan = result.Kecamatan
		data.Kabupaten = result.Kabupaten
		data.Provinsi = result.Provinsi
		data.KodePos = result.Kodepos
		data.TipeDomisili = result.TipeDomisili
		data.KetDomisili = result.KetDomisili
		data.PosisiMap = result.PosisiMap
		data.Lat = result.Lat
		data.Long = result.Long
		data.UpdatedFrom = result.UpdatedFrom
		data.UpdatedDate = result.UpdatedDate.Format(time.DateTime)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Sukses update data profile",
		"data":    data,
	})
}

func (c *UsersProfileController) GetData(ctx *gin.Context) {
	var req Authentication.ValidationGetPersonalInformation
	var data Authentication.PersonalInformationEmployee

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

	personalInformation, err := c.UserProfileRepo.FindProfileUsers(req.NIK)

	if err == nil {
		data.Nik = personalInformation.Nik
		data.Alamat = personalInformation.Alamat
		data.Kelurahan = personalInformation.Kelurahan
		data.Kecamatan = personalInformation.Kecamatan
		data.Kabupaten = personalInformation.Kabupaten
		data.Provinsi = personalInformation.Provinsi
		data.KodePos = personalInformation.Kodepos
		data.TipeDomisili = personalInformation.TipeDomisili
		data.KetDomisili = personalInformation.KetDomisili
		data.PosisiMap = personalInformation.PosisiMap
		data.Lat = personalInformation.Lat
		data.Long = personalInformation.Long
		data.UpdatedFrom = personalInformation.UpdatedFrom
		data.UpdatedDate = personalInformation.UpdatedDate.Format(time.DateTime)

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    data,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    nil,
		})
	}
}

func (c *UsersProfileController) GetCategory(ctx *gin.Context) {
	var data []Authentication.DataDomisili
	var empty Authentication.DataDomisili
	Domisili := c.UserProfileRepo.FindDomisili()

	for _, dm := range Domisili {
		empty.Id = *dm.TipeDomisili
		empty.Name = *dm.KetDomisili
		data = append(data, empty)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data,
	})
}

func (c *UsersProfileController) StoreInformationContact(ctx *gin.Context) {
	var req Authentication.ValidationStoreContactInformation
	var data Authentication.ContactInformation

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

	personalInformation, err := c.UserProfileRepo.FindProfileUsers(req.Nik)
	personalInformation.Nik = req.Nik

	if req.NoTelp1 == "" {
		personalInformation.NoTelp1 = nil
	} else {
		personalInformation.NoTelp1 = &req.NoTelp1
	}

	if req.NoTelp2 == "" {
		personalInformation.NoTelp2 = nil
	} else {
		personalInformation.NoTelp2 = &req.NoTelp2
	}

	if req.Email1 == "" {
		personalInformation.Email1 = nil
	} else {
		personalInformation.Email1 = &req.Email1
	}

	if req.Email2 == "" {
		personalInformation.Email2 = nil
	} else {
		personalInformation.Email2 = &req.Email2
	}

	if err != nil {
		result, _ := c.UserProfileRepo.Create(personalInformation)
		data.NoTelp1 = result.NoTelp1
		data.NoTelp2 = result.NoTelp2
		data.Email1 = result.Email1
		data.Email2 = result.Email2
	} else {
		result, _ := c.UserProfileRepo.Update(personalInformation)
		data.NoTelp1 = result.NoTelp1
		data.NoTelp2 = result.NoTelp2
		data.Email1 = result.Email1
		data.Email2 = result.Email2
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data,
	})
}

func (c *UsersProfileController) GetContactInformation(ctx *gin.Context) {
	Nik := ctx.Param("nik")
	var data Authentication.ContactInformation

	personalInformation, err := c.UserProfileRepo.FindProfileUsers(Nik)

	if err == nil {
		data.NoTelp1 = personalInformation.NoTelp1
		data.NoTelp2 = personalInformation.NoTelp2
		data.Email1 = personalInformation.Email1
		data.Email2 = personalInformation.Email2

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    data,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    nil,
		})
	}
}

func (c *UsersProfileController) StoreProfile(ctx *gin.Context) {
	var req Authentication.ValidationStoreProfile
	var data Authentication.GetStoreProfile

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

	personalInformation, err := c.ProfileRepo.FindProfile(req.NIK)

	personalInformation.Nik = req.NIK
	personalInformation.Bio = req.Bio
	personalInformation.LinkTwitter = req.LinkTwitter
	personalInformation.LinkInstagram = req.LinkInstagram
	personalInformation.LinkWebsite = req.LinkWebsite
	personalInformation.LinkFacebook = req.LinkFacebook
	personalInformation.LinkTiktok = req.LinkTiktok

	if err != nil {
		result, _ := c.ProfileRepo.Create(personalInformation)

		data.ID = result.ID
		data.NIK = result.Nik
		data.Bio = result.Bio
		data.LinkFacebook = result.LinkFacebook
		data.LinkInstagram = result.LinkInstagram
		data.LinkTiktok = result.LinkTiktok
		data.LinkTwitter = result.LinkTwitter
		data.LinkWebsite = result.LinkWebsite
		data.CreatedAt = result.CreatedAt
		data.UpdatedAt = result.UpdatedAt

		ctx.JSON(http.StatusOK, gin.H{
			"ResponseCode":   0,
			"ResponseString": "OK",
			"data":           data,
		})
	} else {
		result, _ := c.ProfileRepo.Update(personalInformation)
		ctx.JSON(http.StatusOK, gin.H{
			"ResponseCode":   0,
			"ResponseString": "OK",
			"data":           result,
		})
	}
}

func (c *UsersProfileController) StoreAboutUs(ctx *gin.Context) {
	var req Authentication.ValidationStoreAboutUs
	// var data Authentication.GetStoreProfile

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

	personalInformation, err := c.AboutUsRepo.FindProfileAboutUs(req.NIK)

	personalInformation.Nik = req.NIK
	if req.Desc == "" {
		personalInformation.AboutUsDesc = nil
	} else {
		personalInformation.AboutUsDesc = &req.Desc
	}

	if req.Hobby == "" {
		personalInformation.AboutUsHobby = nil
	} else {
		personalInformation.AboutUsHobby = &req.Hobby
	}

	if err != nil {
		result, _ := c.AboutUsRepo.Create(personalInformation)

		// data.ID = result.ID
		// data.NIK = result.Nik
		// data.Bio = result.Bio
		// data.LinkFacebook = result.LinkFacebook
		// data.LinkInstagram = result.LinkInstagram
		// data.LinkTiktok = result.LinkTiktok
		// data.LinkTwitter = result.LinkTwitter
		// data.LinkWebsite = result.LinkWebsite
		// data.CreatedAt = result.CreatedAt
		// data.UpdatedAt = result.UpdatedAt

		ctx.JSON(http.StatusOK, gin.H{
			"ResponseCode":   0,
			"ResponseString": "OK",
			"data":           result,
		})
	} else {
		result, _ := c.AboutUsRepo.Update(personalInformation)
		ctx.JSON(http.StatusOK, gin.H{
			"ResponseCode":   0,
			"ResponseString": "OK",
			"data":           result,
		})
	}
}

func (c *UsersProfileController) GetShowAboutUs(ctx *gin.Context) {
	Nik := ctx.Param("nik")

	var data Authentication.ShowAboutUs

	personalInformation, err := c.AboutUsRepo.FindProfileAboutUs(Nik)

	if err == nil {
		data.Id = personalInformation.ID
		data.Nik = personalInformation.Nik
		data.AboutUsDesc = personalInformation.AboutUsDesc
		data.AboutUsHobby = personalInformation.AboutUsHobby
		data.CreatedAt = personalInformation.CreatedAt
		data.UpdatedAt = personalInformation.UpdatedAt

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    data,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"success": "Data tidak ditemukan",
		})
	}
}

func (c *UsersProfileController) GetSocialMediaInformation(ctx *gin.Context) {
	Nik := ctx.Param("nik")
	var data Authentication.GetSocialMedia

	personalInformation, err := c.ProfileRepo.FindProfile(Nik)

	if err == nil {
		data.NIK = personalInformation.Nik
		data.Bio = personalInformation.Bio
		data.LinkTwitter = personalInformation.LinkTwitter
		data.LinkInstagram = personalInformation.LinkInstagram
		data.LinkWebsite = personalInformation.LinkWebsite
		data.LinkFacebook = personalInformation.LinkFacebook
		data.LinkTiktok = personalInformation.LinkTiktok

		ctx.JSON(http.StatusOK, gin.H{
			"ResponseCode":   0,
			"ResponseString": "OK",
			"data":           data,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"success": "Data tidak ditemukan",
		})
	}
}

func (c *UsersProfileController) StoreSkill(ctx *gin.Context) {
	var req Authentication.ValidationStoreSkill

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
	for _, category := range req.Category {
		var skillProfile []profile.ProfileSkill
		keahlian, err := c.ProfileSkillRepo.FindProfileCategorySkill(category.ID)

		if err != nil {
			profilSkillCategory := profile.ProfileSkill{
				Nik:  req.NIK,
				Name: category.Name,
				Type: "category_skill",
			}
			skillProfile = append(skillProfile, profilSkillCategory)

			fmt.Println("Category:", category.Name)
			for _, skill := range category.Skill {
				profilMainSkill := profile.ProfileSkill{
					Nik:  req.NIK,
					Name: skill.Name,
					Type: "main_skill",
				}
				skillProfile = append(skillProfile, profilMainSkill)

				for _, subSkill := range skill.SubSkill {
					profilSubSkill := profile.ProfileSkill{
						Nik:  req.NIK,
						Name: subSkill.Name,
						Type: "sub_skill",
					}
					skillProfile = append(skillProfile, profilSubSkill)
				}
			}
			kebenaran = false
		} else {
			keahlian.Name = category.Name

			skillProfile = append(skillProfile, keahlian)

			for _, skil := range category.Skill {
				mainSkill, _ := c.ProfileSkillRepo.FindProfileSkill(skil.ID, keahlian.ID)

				mainSkill.Name = skil.Name

				skillProfile = append(skillProfile, mainSkill)

				for _, subskil := range skil.SubSkill {
					subSkill, _ := c.ProfileSkillRepo.FindProfileSkill(subskil.ID, mainSkill.ID)

					subSkill.Name = subskil.Name
					skillProfile = append(skillProfile, subSkill)
				}
			}
			kebenaran = true
		}
		if kebenaran {
			c.ProfileSkillRepo.UpdateC(skillProfile)
		} else {
			c.ProfileSkillRepo.CreateC(skillProfile)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
	})
}

func (c *UsersProfileController) UpdateSkill(ctx *gin.Context) {
	var req Authentication.ValidationUpdateSkill

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

	typeCat := "category_skill"
	typeMainSkill := "main_skill"
	typeSubSkill := "sub_skill"

	if req.Type == typeCat {
		catSkill, err := c.ProfileSkillRepo.FindProfileCategorySkill(req.ID)
		if err != nil {
			catSkill.ID = req.ID
			catSkill.Name = req.Name
			catSkill.Nik = req.NIK
			catSkill.Type = req.Type
			c.ProfileSkillRepo.Create(catSkill)
		} else {
			catSkill.Name = req.Name
			c.ProfileSkillRepo.Update(catSkill)
		}
	}
	if req.Type == typeMainSkill {
		mainSkill, err := c.ProfileSkillRepo.FindProfileSkill(req.ID, req.IdParentSkill)
		if err != nil {
			mainSkill.ID = req.ID
			mainSkill.IdParentSkill = &req.IdParentSkill
			mainSkill.Name = req.Name
			mainSkill.Nik = req.NIK
			mainSkill.Type = req.Type
			c.ProfileSkillRepo.Create(mainSkill)
		} else {
			mainSkill.Name = req.Name
			c.ProfileSkillRepo.Update(mainSkill)
		}
	}
	if req.Type == typeSubSkill {
		subSkill, err := c.ProfileSkillRepo.FindProfileSkill(req.ID, req.IdParentSkill)
		if err != nil {
			subSkill.ID = req.ID
			subSkill.IdParentSkill = &req.IdParentSkill
			subSkill.Name = req.Name
			subSkill.Nik = req.NIK
			subSkill.Type = req.Type
			c.ProfileSkillRepo.Create(subSkill)
		} else {
			subSkill.Name = req.Name
			c.ProfileSkillRepo.Update(subSkill)
		}
	}

	var data []Authentication.ShowSkills

	personalCategory, _ := c.ProfileSkillRepo.GetProfileSkillArr(req.NIK, typeCat)

	personalMainSkill, _ := c.ProfileSkillRepo.GetProfileSkillArr(req.NIK, typeMainSkill)

	personalSubSkill, _ := c.ProfileSkillRepo.GetProfileSkillArr(req.NIK, typeSubSkill)

	for _, cat := range personalCategory {
		var mainskill []Authentication.ProfileMainSkill

		for _, mainskll := range personalMainSkill {
			var subskill []Authentication.ProfileSubSkill

			for _, subskll := range personalSubSkill {
				if subskll.IdParentSkill != nil {
					if mainskll.ID == *subskll.IdParentSkill {
						subskill = append(subskill, struct{ profile.ProfileSkill }{subskll})
					}
				}
			}

			if mainskll.IdParentSkill != nil {
				if cat.ID == *mainskll.IdParentSkill {
					mainSkills := Authentication.ProfileMainSkill{
						ProfileSkill: mainskll,
						SubSkill:     subskill,
					}
					mainskill = append(mainskill, mainSkills)
				}
			}
		}

		catSkills := Authentication.ShowSkills{
			ProfileSkill: cat,
			Skill:        mainskill,
		}
		data = append(data, catSkills)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data,
	})
}

// func (c *UsersProfileController) StoreSkill(ctx *gin.Context) {
// 	var req Authentication.ValidationStoreSkill

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

// 	for i, cat := range req.Category {
// 		fmt.Println("MASUK CATEGORY: ", i)

// 		var subskillProfil []profile.ProfileSkill
// 		keahlian, err := c.ProfileSkillRepo.FindProfileCategorySkill(cat.ID)
// 		if err != nil {
// 			profilSkillCategory := profile.ProfileSkill{
// 				Nik:  req.NIK,
// 				Name: cat.Name,
// 				Type: "category_skill",
// 			}
// 			itemsCategory, _ := c.ProfileSkillRepo.Create(profilSkillCategory)

// 			for j, skil := range cat.Skill {
// 				fmt.Println("MASUK SKILL: ", j, "-", i)
// 				profilSkillMain := profile.ProfileSkill{
// 					Nik:           req.NIK,
// 					IdParentSkill: &itemsCategory.ID,
// 					Name:          skil.Name,
// 					Type:          "main_skill",
// 				}
// 				itemsSkill, _ := c.ProfileSkillRepo.Create(profilSkillMain)

// 				for k, subskil := range skil.SubSkill {
// 					fmt.Println("MASUK SUBSKILL: ", k, "-", j, "-", i)
// 					profilSubSkill := profile.ProfileSkill{
// 						Nik:           req.NIK,
// 						IdParentSkill: &itemsSkill.ID,
// 						Name:          subskil.Name,
// 						Type:          "sub_skill",
// 					}
// 					subskillProfil = append(subskillProfil, profilSubSkill)
// 				}
// 				c.ProfileSkillRepo.CreateC(subskillProfil)
// 			}
// 		} else {
// 			keahlian.Name = cat.Name

// 			itemsCategory, _ := c.ProfileSkillRepo.Update(keahlian)

// 			for _, skil := range cat.Skill {
// 				mainSkill, _ := c.ProfileSkillRepo.FindProfileSkill(skil.ID, itemsCategory.ID)

// 				mainSkill.Name = skil.Name

// 				itemsSkill, _ := c.ProfileSkillRepo.Update(mainSkill)

// 				for _, subskil := range skil.SubSkill {
// 					subSkill, _ := c.ProfileSkillRepo.FindProfileSkill(subskil.ID, itemsSkill.ID)

// 					subSkill.Name = subskil.Name
// 					subskillProfil = append(subskillProfil, subSkill)
// 				}
// 				c.ProfileSkillRepo.UpdateC(subskillProfil)
// 			}
// 		}
// 	}
// }

// func (c *UsersProfileController) GetSkill(ctx *gin.Context) {
// 	Nik := ctx.Param("nik")
// 	var data []Authentication.ShowSkills

// 	personalCategory, err := c.ProfileSkillRepo.FindProfileCategorySkillArr(Nik)
// 	for _, cat := range personalCategory {
// 		var mainskill []Authentication.ProfileMainSkill

// 		personalMainSkill, _ := c.ProfileSkillRepo.FindProfileSkillArr(cat.Nik, cat.ID)
// 		for _, mainSkill := range personalMainSkill {
// 			var subskill []Authentication.ProfileSubSkill

// 			personalSubSkill, _ := c.ProfileSkillRepo.FindProfileSkillArr(mainSkill.Nik, mainSkill.ID)

// 			for _, dataSubSkill := range personalSubSkill {
// 				subskill = append(subskill, struct{ profile.ProfileSkill }{dataSubSkill})
// 			}

// 			mainSkills := Authentication.ProfileMainSkill{
// 				ProfileSkill: mainSkill,
// 				SubSkill:     subskill,
// 			}
// 			mainskill = append(mainskill, mainSkills)
// 		}
// 		catSkills := Authentication.ShowSkills{
// 			ProfileSkill: cat,
// 			Skill:        mainskill,
// 		}
// 		data = append(data, catSkills)
// 	}

//		if err == nil {
//			ctx.JSON(http.StatusOK, gin.H{
//				"status":  http.StatusOK,
//				"success": "Success",
//				"data":    data,
//			})
//		} else {
//			ctx.JSON(http.StatusOK, gin.H{
//				"status":  http.StatusOK,
//				"success": "Data Tidak Ditemukan!!",
//				"data":    nil,
//			})
//		}
//	}

func (c *UsersProfileController) GetSkill(ctx *gin.Context) {
	Nik := ctx.Param("nik")
	data := []Authentication.ShowSkills{}

	typeCat := "category_skill"
	personalCategory, _ := c.ProfileSkillRepo.GetProfileSkillArr(Nik, typeCat)

	if personalCategory != nil {
		for _, cat := range personalCategory {
			typeMainSkill := "main_skill"
			personalMainSkill, _ := c.ProfileSkillRepo.GetProfileSkillArr(Nik, typeMainSkill)

			typeSubSkill := "sub_skill"
			personalSubSkill, _ := c.ProfileSkillRepo.GetProfileSkillArr(Nik, typeSubSkill)
			mainskill := []Authentication.ProfileMainSkill{}

			for _, mainskll := range personalMainSkill {
				subskill := []Authentication.ProfileSubSkill{}

				for _, subskll := range personalSubSkill {
					if subskll.IdParentSkill != nil {
						if mainskll.ID == *subskll.IdParentSkill {
							subskill = append(subskill, struct{ profile.ProfileSkill }{subskll})
						}
					}
				}

				if mainskll.IdParentSkill != nil {
					if cat.ID == *mainskll.IdParentSkill {
						mainSkills := Authentication.ProfileMainSkill{
							ProfileSkill: mainskll,
							SubSkill:     subskill,
						}
						mainskill = append(mainskill, mainSkills)
					}
				}
			}

			catSkills := Authentication.ShowSkills{
				ProfileSkill: cat,
				Skill:        mainskill,
			}
			data = append(data, catSkills)
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    data,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Data Tidak Ditemukan!!",
			"data":    data,
		})
	}
}

func (c *UsersProfileController) DeleteSkill(ctx *gin.Context) {
	var req Authentication.ValidationDeleteSkill

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

	var skillProfile []profile.ProfileSkill
	var data []profile.ProfileSkill
	if req.Type == "Kategori" {
		catSkill, err := c.ProfileSkillRepo.FindProfileCategorySkill(req.ID)
		if err == nil {
			skillProfile = append(skillProfile, catSkill)

			mainSkill, err2 := c.ProfileSkillRepo.FindProfileSkillArr(catSkill.ID)
			if err2 == nil {
				for _, dataMainSkill := range mainSkill {
					skillProfile = append(skillProfile, dataMainSkill)

					subSkill, err3 := c.ProfileSkillRepo.FindProfileSkillArr(dataMainSkill.ID)
					if err3 == nil {
						skillProfile = append(skillProfile, subSkill...)
					}
				}
			}
			c.ProfileSkillRepo.DeleteC(skillProfile)

			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
				"data":    data,
			})
		}
	}
	if req.Type == "Skill" {
		mainSkill, err := c.ProfileSkillRepo.FindProfileSkillIndiv(req.ID)

		if err == nil {
			skillProfile = append(skillProfile, mainSkill)

			subSkill, err2 := c.ProfileSkillRepo.FindProfileSkillArr(mainSkill.ID)

			if err2 == nil {
				skillProfile = append(skillProfile, subSkill...)
			}
			c.ProfileSkillRepo.DeleteC(skillProfile)

			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
				"data":    data,
			})
		}
	}
	if req.Type == "Sub" {
		subSkill, err := c.ProfileSkillRepo.FindProfileSkillIndiv(req.ID)

		if err == nil {
			skillProfile = append(skillProfile, subSkill)

			c.ProfileSkillRepo.DeleteC(skillProfile)

			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"success": "Success",
				"data":    data,
			})
		}
	}
}

type ByValidFrom []Authentication.PengalamanKerja

func (a ByValidFrom) Len() int           { return len(a) }
func (a ByValidFrom) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValidFrom) Less(i, j int) bool { return a[i].ValidFrom > a[j].ValidFrom }

func (c *UsersProfileController) GetPengalamanKerja(ctx *gin.Context) {
	Nik := ctx.Param("nik")
	var data []Authentication.PengalamanKerja
	var merged []profile.PengalamanKerja

	if string(Nik[0]) == "8" {
		nik_lama := Nik[1:]
		result, err := c.PengalamanKerjaRepo.FindRiwayatJabatan(nik_lama)
		if err == nil {
			merged = append(merged, result...)
		}
	}
	if string(Nik[0]) == "1" {
		krywn, err := c.PihcKaryawanMutasiPIRepo.FindPihcKaryawanMutasiPI(Nik)
		if err == nil {
			result, err1 := c.PengalamanKerjaRepo.FindRiwayatJabatan(krywn.EmpNo)
			if err1 == nil {
				merged = append(merged, result...)
			}
		}
	}
	result, err := c.PengalamanKerjaRepo.FindRiwayatJabatan(Nik)
	if err == nil {
		merged = append(merged, result...)
	}

	for _, riwayat := range merged {
		// Membuang tanda kurung
		riwayat.RiwayatJabatan = strings.Trim(riwayat.RiwayatJabatan, "()")

		// Membuat pembaca CSV
		reader := csv.NewReader(strings.NewReader(riwayat.RiwayatJabatan))
		reader.Comma = ','
		// Membaca data
		values, err := reader.Read()
		fmt.Println(values)
		if err != nil {
			fmt.Println("Gagal membaca data:", err)
			return
		}

		// Menggunakan regular expression untuk membersihkan nilai-nilai yang dalam tanda kutip
		re := regexp.MustCompile(`"(.+?)"`)
		for i, value := range values {
			matches := re.FindStringSubmatch(value)
			if len(matches) > 1 {
				values[i] = matches[1]
			}
		}

		// Pastikan values memiliki setidaknya 8 elemen (sesuaikan jika perlu)
		if len(values) >= 8 {
			if values[4] != "" {
				// Hapus tanda kutip dari string yang dibungkus dalam tanda kutip
				result := Authentication.PengalamanKerja{
					Grade:        strings.Trim(values[2], " \""),
					PositionId:   strings.Trim(values[3], " \""),
					PositionName: strings.Trim(values[4], " \""),
					Unit1:        strings.Trim(values[6], " \""),
					Unit2:        strings.Trim(values[7], " \""),
				}
				// Get the current year
				currentYear := time.Now().Year()

				// Parse the date string
				parsedTimeValidFrom, _ := time.Parse("2006-01-02 15:04:05", values[0])
				parsedTimeValidTo, _ := time.Parse("2006-01-02 15:04:05", values[1])
				// Replace the year with the current year
				// updatedTime := parsedTimeValidTo.AddDate(9999-currentYear, 0, 0)

				// Format the updated time as a string
				validFrom := parsedTimeValidFrom.Format("2006-01-02")
				validTo := parsedTimeValidTo.Format("2006-01-02")
				validTo = strings.Replace(validTo, "9999", fmt.Sprintf("%d", currentYear), -1)
				result.ValidFrom = strings.Trim(validFrom, " \"")
				result.ValidTo = strings.Trim(validTo, " \"")

				data = append(data, result)
			}
		} else {
			fmt.Println("Data tidak lengkap")
		}
	}

	// Menggunakan sort.Sort dengan jenis khusus ByValidFrom
	sort.Sort(ByValidFrom(data))

	// Menggunakan sort.Slice dengan fungsi penilaian khusus
	// sort.Slice(data, func(i, j int) bool {
	// 	return data[i].ValidFrom > data[j].ValidFrom
	// })

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"success": "Success",
		"data":    data,
	})
}
func (c *UsersProfileController) ShowProfile(ctx *gin.Context) {
	Nik := ctx.Param("nik")
	var data Authentication.ProfilePribadi

	data_karyawan, err := c.PihcMasterKaryDbRepo.FindUserByNIK(Nik)
	if err == nil {
		data.PihcMasterKary.EmpNo = data_karyawan.EmpNo
		data.PihcMasterKary.Nama = data_karyawan.Nama
		data.PihcMasterKary.Gender = data_karyawan.Gender
		data.PihcMasterKary.Agama = data_karyawan.Agama
		data.PihcMasterKary.StatusKawin = data_karyawan.StatusKawin
		data.PihcMasterKary.Anak = data_karyawan.Anak
		data.PihcMasterKary.Mdg = strconv.Itoa(data_karyawan.Mdg)
		data.PihcMasterKary.EmpGrade = data_karyawan.EmpGrade
		data.PihcMasterKary.EmpGradeTitle = data_karyawan.EmpGradeTitle
		data.PihcMasterKary.Area = data_karyawan.Area
		data.PihcMasterKary.AreaTitle = data_karyawan.AreaTitle
		data.PihcMasterKary.SubArea = data_karyawan.SubArea
		data.PihcMasterKary.SubAreaTitle = data_karyawan.SubAreaTitle
		data.PihcMasterKary.Contract = data_karyawan.Contract
		data.PihcMasterKary.Pendidikan = data_karyawan.Pendidikan
		data.PihcMasterKary.Company = data_karyawan.Company
		if data_karyawan.Lokasi != "" {
			data.PihcMasterKary.Lokasi = &data_karyawan.Lokasi
		}
		data.PihcMasterKary.EmployeeStatus = data_karyawan.EmployeeStatus
		data.PihcMasterKary.Email = data_karyawan.Email
		data.PihcMasterKary.HP = data_karyawan.HP
		data.PihcMasterKary.TglLahir = data_karyawan.TglLahir.Format(time.DateOnly)
		data.PihcMasterKary.PosID = data_karyawan.PosID
		data.PihcMasterKary.PosTitle = data_karyawan.PosTitle
		data.PihcMasterKary.SupPosID = data_karyawan.SupPosID
		data.PihcMasterKary.PosGrade = data_karyawan.PosGrade
		data.PihcMasterKary.PosKategori = data_karyawan.PosKategori
		data.PihcMasterKary.OrgID = data_karyawan.OrgID
		data.PihcMasterKary.OrgTitle = data_karyawan.OrgTitle
		data.PihcMasterKary.DeptID = data_karyawan.DeptID
		data.PihcMasterKary.DeptTitle = data_karyawan.DeptTitle
		data.PihcMasterKary.KompID = data_karyawan.KompID
		data.PihcMasterKary.KompTitle = data_karyawan.KompTitle
		data.PihcMasterKary.DirID = data_karyawan.DirID
		data.PihcMasterKary.DirTitle = data_karyawan.DirTitle
		data.PihcMasterKary.PosLevel = data_karyawan.PosLevel
		data.PihcMasterKary.SupEmpNo = data_karyawan.SupEmpNo
		data.PihcMasterKary.BagID = data_karyawan.BagID
		data.PihcMasterKary.BagTitle = data_karyawan.BagTitle
		if data_karyawan.SeksiID != "" {
			data.PihcMasterKary.SeksiID = &data_karyawan.SeksiID
		}
		if data_karyawan.SeksiTitle != "" {
			data.PihcMasterKary.SeksiTitle = &data_karyawan.SeksiTitle
		}
		if data_karyawan.PreNameTitle != "" {
			data.PihcMasterKary.PreNameTitle = &data_karyawan.PreNameTitle
		}
		if data_karyawan.PostNameTitle != "" {
			data.PihcMasterKary.PostNameTitle = &data_karyawan.PostNameTitle
		}
		if data_karyawan.NoNPWP != "" {
			data.PihcMasterKary.NoNPWP = &data_karyawan.NoNPWP
		}
		if data_karyawan.BankAccount != "" {
			data.PihcMasterKary.BankAccount = &data_karyawan.BankAccount
		}
		if data_karyawan.BankName != "" {
			data.PihcMasterKary.BankName = &data_karyawan.BankName
		}
		data.PihcMasterKary.MdgDate = data_karyawan.MdgDate
		if data_karyawan.PayScale != "" {
			data.PihcMasterKary.PayScale = &data_karyawan.PayScale
		}
		data.PihcMasterKary.CCCode = data_karyawan.CCCode
		data.PihcMasterKary.Nickname = data_karyawan.Nickname

		domisili, _ := c.UserProfileRepo.FindProfileUsers(data.EmpNo)
		if domisili.Nik != "" {
			data_domisili := users.UserProfile{
				Nik:         domisili.Nik,
				Alamat:      domisili.Alamat,
				Kelurahan:   domisili.Kelurahan,
				Kecamatan:   domisili.Kecamatan,
				Kabupaten:   domisili.Kabupaten,
				Provinsi:    domisili.Provinsi,
				Kodepos:     domisili.Kodepos,
				Domisili:    domisili.Domisili,
				PosisiMap:   domisili.PosisiMap,
				Email2:      domisili.Email2,
				UpdatedBy:   domisili.UpdatedBy,
				NoTelp1:     domisili.NoTelp1,
				NoTelp2:     domisili.NoTelp2,
				Lat:         domisili.Lat,
				Long:        domisili.Long,
				Email1:      domisili.Email1,
				UpdatedFrom: domisili.UpdatedFrom,
				UpdatedDate: domisili.UpdatedDate.Format(time.DateTime),
				IsAdmin:     domisili.IsAdmin,
			}

			data.Domisili = &data_domisili
		}

		data_profile, _ := c.ProfileRepo.FindProfile(domisili.Nik)
		if data_profile.ID != 0 {
			profileMobile := &Authentication.MobileProfile{
				Profile:     data_profile,
				UserProfile: *data.Domisili,
			}
			data.ProfileMobile = profileMobile
		}

		about, _ := c.AboutUsRepo.FindProfileAboutUs(data_profile.Nik)
		if about.ID != 0 {
			data.AboutUs = &about
		}

		company, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(data_karyawan.Company)
		data.Companys = company

		typeCat := "category_skill"
		personalCategory, _ := c.ProfileSkillRepo.GetProfileSkillArr(data_karyawan.EmpNo, typeCat)
		if personalCategory != nil {
			typeMainSkill := "main_skill"
			personalMainSkill, _ := c.ProfileSkillRepo.GetProfileSkillArr(data_karyawan.EmpNo, typeMainSkill)

			typeSubSkill := "sub_skill"
			personalSubSkill, _ := c.ProfileSkillRepo.GetProfileSkillArr(data_karyawan.EmpNo, typeSubSkill)
			for _, cat := range personalCategory {
				mainskill := []Authentication.ProfileMainSkill{}

				for _, mainskll := range personalMainSkill {
					subskill := []Authentication.ProfileSubSkill{}

					for _, subskll := range personalSubSkill {
						if subskll.IdParentSkill != nil {
							if mainskll.ID == *subskll.IdParentSkill {
								subskill = append(subskill, struct{ profile.ProfileSkill }{subskll})
							}
						}
					}

					if mainskll.IdParentSkill != nil {
						if cat.ID == *mainskll.IdParentSkill {
							mainSkills := Authentication.ProfileMainSkill{
								ProfileSkill: mainskll,
								SubSkill:     subskill,
							}
							mainskill = append(mainskill, mainSkills)
						}
					}
				}

				catSkills := Authentication.ShowSkills{
					ProfileSkill: cat,
					Skill:        mainskill,
				}
				data.Skill = append(data.Skill, catSkills)
			}
		} else {
			data.Skill = []Authentication.ShowSkills{}
		}

		data.CompanyLogo = "https://storage.googleapis.com/lumen-oauth-storage/company/logo-pi-full.png"
		photoProfile, _ := c.PhotoProfileRepo.FindPhotoProfile(data_karyawan.EmpNo)
		if photoProfile.Url != "" {
			data.PhotoProfile = photoProfile.Url
		} else {
			data.PhotoProfile = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
		}
		data.PhotoProfileDefault = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"

		organization, _ := c.ViewOrganisasiRepo.FindViewOrganization(data_karyawan.EmpNo)

		data.Organisasi = append(data.Organisasi, organization.Unit1)
		data.Organisasi = append(data.Organisasi, organization.Unit2)
		data.Organisasi = append(data.Organisasi, organization.Org3)
		data.Organisasi = append(data.Organisasi, organization.Org4)

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    data,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
func (c *UsersProfileController) UpdatePhotoProfile(ctx *gin.Context) {
	var req Authentication.ValidationPhotoProfile
	var data Authentication.ProfilePribadi
	var photoProfile profile.PhotoProfile

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

	data_karyawan, err := c.PihcMasterKaryDbRepo.FindUserByNIK(req.NIK)
	if err == nil {
		file, _ := ctx.FormFile("photo")

		originalFileName := file.Filename
		fmt.Println(originalFileName)

		fileToUpload, err := file.Open()
		if err != nil {
			// Handle the error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not open file",
			})
			return
		}

		imageURL, file_name, err := c.PhotoProfileRepo.UploadFilePhotoProfile(data_karyawan.EmpNo, originalFileName, fileToUpload)
		// file_url, file_name, err := c.EventNotulenRepo.UploadFile(originalFileName, fileToUpload)
		if err == nil {
			pp, _ := c.PhotoProfileRepo.FindPhotoProfile(data_karyawan.EmpNo)
			fmt.Println(pp.Id)
			if pp.Id != 0 {
				fmt.Println("UPDATE")
				pp.Name = file_name
				pp.Url = imageURL

				updatePP, _ := c.PhotoProfileRepo.Update(pp)
				photoProfile = updatePP
			} else {
				fmt.Println("CREATE")
				pp.Name = file_name
				pp.Url = imageURL
				pp.EmpNo = data_karyawan.EmpNo

				createPP, _ := c.PhotoProfileRepo.Create(pp)
				photoProfile = createPP
			}
		}
		if err != nil {
			// Handle the error
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not upload file",
			})
			return
		}
		data.PihcMasterKary.EmpNo = data_karyawan.EmpNo
		data.PihcMasterKary.Nama = data_karyawan.Nama
		data.PihcMasterKary.Gender = data_karyawan.Gender
		data.PihcMasterKary.Agama = data_karyawan.Agama
		data.PihcMasterKary.StatusKawin = data_karyawan.StatusKawin
		data.PihcMasterKary.Anak = data_karyawan.Anak
		data.PihcMasterKary.Mdg = strconv.Itoa(data_karyawan.Mdg)
		data.PihcMasterKary.EmpGrade = data_karyawan.EmpGrade
		data.PihcMasterKary.EmpGradeTitle = data_karyawan.EmpGradeTitle
		data.PihcMasterKary.Area = data_karyawan.Area
		data.PihcMasterKary.AreaTitle = data_karyawan.AreaTitle
		data.PihcMasterKary.SubArea = data_karyawan.SubArea
		data.PihcMasterKary.SubAreaTitle = data_karyawan.SubAreaTitle
		data.PihcMasterKary.Contract = data_karyawan.Contract
		data.PihcMasterKary.Pendidikan = data_karyawan.Pendidikan
		data.PihcMasterKary.Company = data_karyawan.Company
		if data_karyawan.Lokasi != "" {
			data.PihcMasterKary.Lokasi = &data_karyawan.Lokasi
		}
		data.PihcMasterKary.EmployeeStatus = data_karyawan.EmployeeStatus
		data.PihcMasterKary.Email = data_karyawan.Email
		data.PihcMasterKary.HP = data_karyawan.HP
		data.PihcMasterKary.TglLahir = data_karyawan.TglLahir.Format(time.DateOnly)
		data.PihcMasterKary.PosID = data_karyawan.PosID
		data.PihcMasterKary.PosTitle = data_karyawan.PosTitle
		data.PihcMasterKary.SupPosID = data_karyawan.SupPosID
		data.PihcMasterKary.PosGrade = data_karyawan.PosGrade
		data.PihcMasterKary.PosKategori = data_karyawan.PosKategori
		data.PihcMasterKary.OrgID = data_karyawan.OrgID
		data.PihcMasterKary.OrgTitle = data_karyawan.OrgTitle
		data.PihcMasterKary.DeptID = data_karyawan.DeptID
		data.PihcMasterKary.DeptTitle = data_karyawan.DeptTitle
		data.PihcMasterKary.KompID = data_karyawan.KompID
		data.PihcMasterKary.KompTitle = data_karyawan.KompTitle
		data.PihcMasterKary.DirID = data_karyawan.DirID
		data.PihcMasterKary.DirTitle = data_karyawan.DirTitle
		data.PihcMasterKary.PosLevel = data_karyawan.PosLevel
		data.PihcMasterKary.SupEmpNo = data_karyawan.SupEmpNo
		data.PihcMasterKary.BagID = data_karyawan.BagID
		data.PihcMasterKary.BagTitle = data_karyawan.BagTitle
		if data_karyawan.SeksiID != "" {
			data.PihcMasterKary.SeksiID = &data_karyawan.SeksiID
		}
		if data_karyawan.SeksiTitle != "" {
			data.PihcMasterKary.SeksiTitle = &data_karyawan.SeksiTitle
		}
		if data_karyawan.PreNameTitle != "" {
			data.PihcMasterKary.PreNameTitle = &data_karyawan.PreNameTitle
		}
		if data_karyawan.PostNameTitle != "" {
			data.PihcMasterKary.PostNameTitle = &data_karyawan.PostNameTitle
		}
		if data_karyawan.NoNPWP != "" {
			data.PihcMasterKary.NoNPWP = &data_karyawan.NoNPWP
		}
		if data_karyawan.BankAccount != "" {
			data.PihcMasterKary.BankAccount = &data_karyawan.BankAccount
		}
		if data_karyawan.BankName != "" {
			data.PihcMasterKary.BankName = &data_karyawan.BankName
		}
		data.PihcMasterKary.MdgDate = data_karyawan.MdgDate
		if data_karyawan.PayScale != "" {
			data.PihcMasterKary.PayScale = &data_karyawan.PayScale
		}
		data.PihcMasterKary.CCCode = data_karyawan.CCCode
		data.PihcMasterKary.Nickname = data_karyawan.Nickname

		domisili, _ := c.UserProfileRepo.FindProfileUsers(data.EmpNo)

		if domisili.Nik != "" {
			data_domisili := users.UserProfile{
				Nik:         domisili.Nik,
				Alamat:      domisili.Alamat,
				Kelurahan:   domisili.Kelurahan,
				Kecamatan:   domisili.Kecamatan,
				Kabupaten:   domisili.Kabupaten,
				Provinsi:    domisili.Provinsi,
				Kodepos:     domisili.Kodepos,
				Domisili:    domisili.Domisili,
				PosisiMap:   domisili.PosisiMap,
				Email2:      domisili.Email2,
				UpdatedBy:   domisili.UpdatedBy,
				NoTelp1:     domisili.NoTelp1,
				NoTelp2:     domisili.NoTelp2,
				Lat:         domisili.Lat,
				Long:        domisili.Long,
				Email1:      domisili.Email1,
				UpdatedFrom: domisili.UpdatedFrom,
				UpdatedDate: domisili.UpdatedDate.Format(time.DateTime),
				IsAdmin:     domisili.IsAdmin,
			}

			data.Domisili = &data_domisili
		}

		data_profile, _ := c.ProfileRepo.FindProfile(domisili.Nik)
		if data_profile.ID != 0 {
			profileMobile := &Authentication.MobileProfile{
				Profile:     data_profile,
				UserProfile: *data.Domisili,
			}
			data.ProfileMobile = profileMobile
		}

		about, _ := c.AboutUsRepo.FindProfileAboutUs(data_profile.Nik)
		if about.ID != 0 {
			data.AboutUs = &about
		}

		company, _ := c.PihcMasterCompanyRepo.FindPihcMsterCompany(data_karyawan.Company)
		data.Companys = company

		typeCat := "category_skill"
		personalCategory, _ := c.ProfileSkillRepo.GetProfileSkillArr(data_karyawan.EmpNo, typeCat)
		if personalCategory != nil {
			typeMainSkill := "main_skill"
			personalMainSkill, _ := c.ProfileSkillRepo.GetProfileSkillArr(data_karyawan.EmpNo, typeMainSkill)

			typeSubSkill := "sub_skill"
			personalSubSkill, _ := c.ProfileSkillRepo.GetProfileSkillArr(data_karyawan.EmpNo, typeSubSkill)
			for _, cat := range personalCategory {
				mainskill := []Authentication.ProfileMainSkill{}

				for _, mainskll := range personalMainSkill {
					subskill := []Authentication.ProfileSubSkill{}

					for _, subskll := range personalSubSkill {
						if subskll.IdParentSkill != nil {
							if mainskll.ID == *subskll.IdParentSkill {
								subskill = append(subskill, struct{ profile.ProfileSkill }{subskll})
							}
						}
					}

					if mainskll.IdParentSkill != nil {
						if cat.ID == *mainskll.IdParentSkill {
							mainSkills := Authentication.ProfileMainSkill{
								ProfileSkill: mainskll,
								SubSkill:     subskill,
							}
							mainskill = append(mainskill, mainSkills)
						}
					}
				}

				catSkills := Authentication.ShowSkills{
					ProfileSkill: cat,
					Skill:        mainskill,
				}
				data.Skill = append(data.Skill, catSkills)
			}
		} else {
			data.Skill = []Authentication.ShowSkills{}
		}

		data.CompanyLogo = "https://storage.googleapis.com/lumen-oauth-storage/company/logo-pi-full.png"
		if photoProfile.Url != "" {
			data.PhotoProfile = photoProfile.Url
		} else {
			data.PhotoProfile = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
		}
		data.PhotoProfileDefault = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"

		organization, _ := c.ViewOrganisasiRepo.FindViewOrganization(data_karyawan.EmpNo)

		data.Organisasi = append(data.Organisasi, organization.Unit1)
		data.Organisasi = append(data.Organisasi, organization.Unit2)
		data.Organisasi = append(data.Organisasi, organization.Org3)
		data.Organisasi = append(data.Organisasi, organization.Org4)

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": "Success",
			"data":    data,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
