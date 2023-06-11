package dao

import (
	"majiang/model"
)

func CreateUser(user *model.User) error {
	result := DB.Create(&user)
	return result.Error
}

func SearchUserById(uid int) (user model.User, err error) {
	result := DB.Where(&model.User{ID: uid}).First(&user)
	return user, result.Error
}

func SearchUserByName(username string) (user model.User, err error) {
	result := DB.Where(&model.User{Username: username}).First(&user)
	return user, result.Error
}
