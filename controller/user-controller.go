package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"

	// erroroauth "github.com/go-oauth2/oauth2/v4/errors"

	"github.com/go-playground/validator/v10"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	role "github.com/yusufwira/lern-golang-gin/entity/public/role"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl"
	users "github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type UsersController struct {
	UserRepo               *users.UserRepo
	OauthClientRepo        *users.OauthClientRepo
	PihcMasterKaryRtDbRepo *pihc.PihcMasterKaryRtDbRepo
	KegiatanKaryawanRepo   *tjsl.KegiatanKaryawanRepo
	ModelHasRoleRepo       *role.ModelHasRoleRepo
	RolesRepo              *role.RolesRepo
}

func NewUserController(db *gorm.DB, StorageClient *storage.Client) *UsersController {
	return &UsersController{UserRepo: users.NewUserRepo(db),
		OauthClientRepo:        users.NewOauthClientRepo(db),
		PihcMasterKaryRtDbRepo: pihc.NewPihcMasterKaryRtDbRepo(db),
		ModelHasRoleRepo:       role.NewModelHasRoleRepo(db),
		RolesRepo:              role.NewRolesRepo(db),
		KegiatanKaryawanRepo:   tjsl.NewKegiatanKaryawanRepo(db, StorageClient)}
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

func (c *UsersController) Index() []users.User {
	var user users.User
	data := c.UserRepo.GetAll(user)
	return data
}

func (c *UsersController) GetData(ctx *gin.Context) []users.User {
	var user users.User
	id := ctx.Param("id")
	data := c.UserRepo.GetUsersID(user, id, ctx)
	return data
}

func (c *UsersController) GetDataKaryawanNameAll(ctx *gin.Context) {
	var req Authentication.ValidationGetName

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

	data, err := c.PihcMasterKaryRtDbRepo.FindUserByNameArr(req.Name, req.Nik)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"data":   "Data Tidak Ditemukan!!",
		})
	}
}
func (c *UsersController) GetDataKaryawanAll(ctx *gin.Context) {
	data, err := c.PihcMasterKaryRtDbRepo.FindUserArr()

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"data":   "Data Tidak Ditemukan!!",
		})
	}
}
func (c *UsersController) GetAtasanPegawai(ctx *gin.Context) {
	// var req Authentication.ValidasiKonfirmasiNik
	Nik := ctx.Param("nik")

	if Nik != "" {
		url := "https://api-pismart-dev.pupuk-indonesia.com/oauth_api/api/delegasi/tryFindSuperios/" + Nik
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

		temp := Authentication.DataAtasan{
			EmpData: data_atasan,
			EmpNo:   data_atasan.EmpNo,
			Photo:   foto,
		}

		var atasan []Authentication.DataAtasan
		atasan = append(atasan, temp)

		ctx.JSON(http.StatusOK, gin.H{
			"status": 200,
			"data":   atasan,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

}

func (c *UsersController) GetDataKaryawanNameIndiv(ctx *gin.Context) {
	var req Authentication.ValidationGetName

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

	data, err := c.PihcMasterKaryRtDbRepo.FindUserByNameIndiv(req.Name, req.Nik)

	if err == nil {
		foto := "https://storage.googleapis.com/lumen-oauth-storage/DataKaryawan/Foto/" + data.Company + "/" + data.EmpNo + ".jpg"
		respons, err := http.Get(foto)
		if err != nil || respons.StatusCode != http.StatusOK {
			foto = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data":   data,
			"photos": foto,
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"data":   "Data Tidak Ditemukan!!",
		})
	}

}

func (c *UsersController) DelData(ctx *gin.Context) []users.User {
	var user users.User
	id := ctx.Param("id")
	data := c.UserRepo.DelUsersID(user, id, ctx)
	return data
}

func (c *UsersController) UpData(ctx *gin.Context) []users.User {
	var user users.User
	id := ctx.Param("id")
	ctx.BindJSON(&user)
	data := c.UserRepo.UpUsersID(user, id, ctx)
	return data
}

func (c *UsersController) Store(ctx *gin.Context) users.User {
	var user users.User
	user.Username = ctx.PostForm("Username")
	user.Password = ctx.PostForm("Password")
	ctx.BindJSON(&user)
	c.UserRepo.Create(user)
	return user
}

func (c *UsersController) Login(ctx *gin.Context) {
	var input Authentication.ValidationLogin

	if err := ctx.ShouldBind(&input); err != nil {
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

	user, err := c.UserRepo.LoginCheck(input.Username, input.Password)

	if err == nil {
		// FROM DATABASE
		clients, _ := c.OauthClientRepo.FindOauthClient(user.Id)
		client_id := strconv.FormatUint(uint64(clients.Id), 10)

		values := url.Values{}
		values.Set("grant_type", "password")
		values.Set("client_id", client_id)
		values.Set("client_secret", clients.Secret)
		values.Set("username", user.Username)
		values.Set("password", user.Password)
		values.Set("key", input.Password)

		resp, err1 := http.Get("http://localhost:9096/api/token?" + values.Encode())
		if err1 != nil {
			fmt.Println("ERRORRR")
		}
		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println("ERRORRR2")
		}

		trimmedBody := bytes.TrimSpace(body)
		var data Authentication.Token
		if err := json.Unmarshal(trimmedBody, &data); err != nil {
			return
		}

		karyawan, _ := c.PihcMasterKaryRtDbRepo.FindUserByNIK(user.Nik)
		is_superior := c.PihcMasterKaryRtDbRepo.IsSuperior(karyawan.PosID)

		roles, _ := c.RolesRepo.FindRoleByUser(karyawan.EmpNo)
		defaultRole := role.MyRole{
			Name:     "KARYAWAN",
			CompCode: nil,
		}
		if roles == nil {
			roles = append(roles, defaultRole)
		}

		var foto *string
		url := "https://storage.googleapis.com/lumen-oauth-storage/DataKaryawan/Foto/" + karyawan.Company + "/" + karyawan.EmpNo + ".jpg"
		foto = &url
		respons, err := http.Get(*foto)
		if err != nil || respons.StatusCode != http.StatusOK {
			default_foto := "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
			foto = &default_foto
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":         http.StatusOK,
			"success":        "Login Success",
			"user_key":       clients.Secret,
			"user_id":        user.Id,
			"user_name":      user.Name,
			"comp_code":      karyawan.Company,
			"email":          user.Email,
			"hp":             karyawan.HP,
			"user_org_name":  karyawan.OrgTitle,
			"user_dept_name": karyawan.DeptTitle,
			"user_komp_name": karyawan.KompTitle,
			"user_pos_name":  karyawan.PosTitle,
			"model_type":     user.UserType,
			"nik":            user.Nik,
			"position":       karyawan.PosID,
			"roles":          roles,
			"is_superior":    is_superior,
			"token":          data,
			"photo_karyawan": foto,
		})
	} else {
		if len(input.Password) < 8 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "The password must be at least 8 characters."},
			)
		} else if user.Id != 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Username dan password kurang benar"},
			)
		} else if user.Id == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Data karyawan belum terdapat pada database PISMART"},
			)
		}
	}
}

func (c *UsersController) Register(ctx *gin.Context) {
	var input Authentication.ValidationRegister

	if err := ctx.ShouldBind(&input); err != nil {
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
	user, err := c.UserRepo.RegisterCheck(input.Username, input.Password)
	if err != nil {
		user.Username = input.Username
		user.Password = input.Password
		user.Email = input.Email
		user.Nik = input.Nik
		user.UserType = input.Type

		c.UserRepo.Create(user)

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   "Data Berhasil Ditambahkan",
		})
	} else {
		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"status": http.StatusServiceUnavailable,
			"data":   nil,
		})
	}
}
