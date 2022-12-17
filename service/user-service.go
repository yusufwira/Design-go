package service

import (
	"github.com/yusufwira/lern-golang-gin/connection"
	"github.com/yusufwira/lern-golang-gin/entity"
)

type UserService interface {
	Save(entity.User) entity.User
	FindAll() []entity.User
	GetAll() []entity.User
}

type userService struct {
	users []entity.User
}

func New() UserService {
	return &userService{}
}

func (service *userService) Save(user entity.User) entity.User {
	db := connection.Database()
	db.Table("User").Create(user)
	return user
}

func (service *userService) FindAll() []entity.User {
	return nil
}

func (service *userService) GetAll() []entity.User {
	var user []entity.User
	db := connection.Database()
	db.Table("User").Find(&user)
	return user
}
