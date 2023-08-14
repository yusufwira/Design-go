package service

import (
	"github.com/yusufwira/lern-golang-gin/connection"
	users "github.com/yusufwira/lern-golang-gin/entity/users"
)

type UserService interface {
	Save(users.User) users.User
	FindAll() []users.User
	GetAll() []users.User
	GetUsersID(string) []users.User
}

type userService struct {
	users []users.User
}

func New() UserService {
	return &userService{}
}

func (service *userService) Save(user users.User) users.User {
	db := connection.Database()
	db.Table("User").Create(user)
	return user
}

func (service *userService) FindAll() []users.User {
	return nil
}

func (service *userService) GetAll() []users.User {
	var user []users.User
	db := connection.Database()
	db.Table("public.users").Find(&user)
	return user
}

func (service *userService) GetUsersID(id string) []users.User {
	var user []users.User
	db := connection.Database()
	db.Table("public.users").Where("id = ?", id).Take(&user)
	return user
}
