package service

import (
	"majiang/dao"
	"majiang/model"
)

func CreateUser(user *model.User) error {
	return dao.CreateUser(user)
}

func SearchUserById(uid int) (user model.User, err error) {
	return dao.SearchUserById(uid)
}

func SearchUserByName(username string) (user model.User, err error) {
	return dao.SearchUserByName(username)
}
