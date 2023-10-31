package controller

import (
	"errors"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl"
	users "github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type UsersController struct {
	UserRepo               *users.UserRepo
	OauthClientRepo        *users.OauthClientRepo
	PihcMasterKaryRtDbRepo *pihc.PihcMasterKaryRtDbRepo
	KegiatanKaryawanRepo   *tjsl.KegiatanKaryawanRepo
}

func NewUserController(db *gorm.DB, StorageClient *storage.Client) *UsersController {
	return &UsersController{UserRepo: users.NewUserRepo(db),
		OauthClientRepo:        users.NewOauthClientRepo(db),
		PihcMasterKaryRtDbRepo: pihc.NewPihcMasterKaryRtDbRepo(db),
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

	// name := ctx.PostForm("name")
	// nik := ctx.PostForm("nik")
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

	// name := ctx.PostForm("name")
	// nik := ctx.PostForm("nik")
	data, err := c.PihcMasterKaryRtDbRepo.FindUserByNameIndiv(req.Name, req.Nik)

	if err == nil {
		photos, err1 := c.KegiatanKaryawanRepo.FindPhotosKaryawan(data.EmpNo, data.Company)

		if err1 != nil {
			photos = "https://t3.ftcdn.net/jpg/03/46/83/96/360_F_346839683_6nAPzbhpSkIpb8pmAwufkC7c5eD7wYws.jpg"
		} else {
			photos = "https://storage.googleapis.com/" + photos
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data":   data,
			"photos": photos,
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
	// user.Name = ctx.PostForm("Name")
	// user.Email = ctx.PostForm("Email")
	ctx.BindJSON(&user)
	c.UserRepo.Create(user)
	return user
}

func (c *UsersController) Login(ctx *gin.Context) {
	var input Authentication.ValidationLogin

	// if err := ctx.ShouldBindJSON(&input); err != nil {
	// 	ctx.JSON(http.StatusNotFound, gin.H{"error": "Username / Password Tidak Boleh Kosong"})
	// 	return
	// }
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
		// clients, _ := c.OauthClientRepo.FindOauthClient(user.Id)
		// client_id := strconv.FormatUint(uint64(clients.Id), 10)

		// manager := manage.NewDefaultManager()
		// fmt.Println("X1")
		// manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
		// fmt.Println("X2")

		// // token store
		// manager.MustTokenStorage(store.NewMemoryTokenStore())
		// fmt.Println("X3")

		// // generate jwt access token
		// manager.MapAccessGenerate(generates.NewAccessGenerate())
		// fmt.Println("X4")

		// // client store
		// clientStore := store.NewClientStore()
		// fmt.Println("X5")

		// clientInfo := &models.Client{
		// 	ID:     client_id,
		// 	Secret: clients.Secret,
		// 	Domain: "http://localhost:9094",
		// }
		// fmt.Println("X6")
		// clientStore.Set(client_id, clientInfo)
		// fmt.Println("X7")
		// manager.MapClientStorage(clientStore)
		// fmt.Println("X8")

		// // Define a custom grant type handler
		// manager.SetAuthorizeCodeTokenCfg(&manage.Config{AccessTokenExp: 0})
		// fmt.Println("X9")

		// // SetClientTokenCfg set the client grant token config
		// manager.SetClientTokenCfg(&manage.Config{
		// 	AccessTokenExp:    time.Hour * 24 * 365 * 10, // 10 years
		// 	RefreshTokenExp:   time.Hour * 24 * 365 * 10, // 10 years
		// 	IsGenerateRefresh: true,
		// })
		// fmt.Println("X10")

		// // Set the access token and refresh token configuration
		// manager.SetPasswordTokenCfg(&manage.Config{
		// 	AccessTokenExp:    time.Hour * 24 * 365 * 10, // 10 years
		// 	RefreshTokenExp:   time.Hour * 24 * 365 * 10, // 10 years
		// 	IsGenerateRefresh: true,
		// })
		// fmt.Println("X11")

		// // Initialize the oauth2 service
		// // ginserver.InitServerWithManager(newGinServer, manager)
		// ginserver.InitServer(manager)
		// fmt.Println("X12")
		// ginserver.SetAllowGetAccessRequest(true)
		// fmt.Println("X13")
		// ginserver.SetClientInfoHandler(server.ClientFormHandler)
		// fmt.Println("X14")
		// ginserver.SetAllowedGrantType("password")
		// fmt.Println("X15")
		// ginserver.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		// 	// Implement your authentication logic here and return the userID if valid
		// 	if username == user.Username && password == user.Password {
		// 		userID = client_id
		// 		fmt.Println("X16")
		// 	}
		// 	fmt.Println("X17")
		// 	return
		// })
		// fmt.Println("X18")

		// url := "http://localhost:9096/oauth2/token?grant_type=password&client_id=" + client_id + "&client_secret=" + clients.Secret + "&username=" + user.Username + "&password=" + user.Password
		// resp, err1 := http.Get(url)
		// fmt.Println("X19")
		// if err1 != nil {
		// 	fmt.Println("ERRORRR")
		// }
		// fmt.Println("X20")
		// body, err2 := ioutil.ReadAll(resp.Body)
		// if err2 != nil {
		// 	fmt.Println("ERRORRR2")
		// }
		// fmt.Println("X21")

		// var data Authentication.Token
		// if err := json.Unmarshal(body, &data); err != nil {
		// 	fmt.Println("X22")
		// 	fmt.Println("Error unmarshaling JSON:", err)
		// 	return
		// }
		// fmt.Println("X23")

		karyawan, _ := c.PihcMasterKaryRtDbRepo.FindUserByNIK(user.Nik)
		ctx.JSON(http.StatusOK, gin.H{
			"status":         http.StatusOK,
			"user_id":        user.Id,
			"user_name":      user.Name,
			"comp_code":      karyawan.Company,
			"email":          user.Email,
			"hp":             karyawan.HP,
			"user_org_name":  karyawan.OrgTitle,
			"user_dept_name": karyawan.DeptTitle,
			"model_type":     user.UserType,
			"nik":            user.Nik,
			"position":       karyawan.PosID,
			// "token":          data,
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
