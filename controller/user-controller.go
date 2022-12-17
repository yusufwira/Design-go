package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/yusufwira/lern-golang-gin/entity"
)

type UserController interface {
	Index() []entity.User
	Store(ctx *gin.Context) entity.User
}

type controller struct {
}

func New() UserController {
	return &controller{}
}

func (c *controller) Index() []entity.User {
	var user entity.User
	data := user.GetAll()
	return data
}

func (c *controller) Store(ctx *gin.Context) entity.User {
	var user entity.User
	user.Username = ctx.PostForm("Username")
	user.Password = ctx.PostForm("Password")
	user.Name = ctx.PostForm("Name")
	user.Email = ctx.PostForm("Email")
	user.Save()
	return user
}
