package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Authentication "github.com/yusufwira/lern-golang-gin/entity/authentication"
	users "github.com/yusufwira/lern-golang-gin/entity/users"
	"gorm.io/gorm"
)

type UsersController struct {
	UserRepo *users.UserRepo
}

func NewUserController(db *gorm.DB) *UsersController {
	return &UsersController{UserRepo: users.NewUserRepo(db)}
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
	var input Authentication.AuthenticationInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Username / Password Tidak Boleh Kosong"})
		return
	}

	user, err := c.UserRepo.LoginCheck(input.Username, input.Password)

	if err == nil {
		ctx.JSON(http.StatusOK, user)
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Username / Password Salah"})
	}
}
