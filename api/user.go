package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"majiang/model"
	"majiang/service"
	"majiang/util"
)

func Register(c *gin.Context) {
	//receive
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		util.RespNormErr(c, util.BindingQueryErrCode)
		return
	}
	//检索数据库
	mysqlUser, err := service.SearchUserByName(user.Username)
	if mysqlUser.Username != "" {
		util.RespNormErr(c, util.RepeatedUsernameErrCode)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		util.RespInternalErr(c)
	}
	//生成盐值
	var salt []byte
	salt, err = util.GenerateSalt()
	if err != nil {
		log.Println(err)
		util.RespInternalErr(c)
	}
	//加密
	hashedPassword := util.HashWithSalt(user.Password, salt)
	//用户信息写入数据库
	user.Password = string(hashedPassword)
	user.Salt = salt
	if err = service.CreateUser(&user); err != nil {
		util.RespInternalErr(c)
	}
	//response
	util.RespOK(c)
}

func Login(c *gin.Context) {
	//receive
	var user model.User
	log.Println("0")
	if err := c.ShouldBind(&user); err != nil {
		log.Println(err)
		util.RespNormErr(c, util.BindingQueryErrCode)
		return
	}
	log.Println("1")
	mysqlUser, err := service.SearchUserByName(user.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			util.RespNormErr(c, util.NoRecordErrCode)
		} else {
			util.RespInternalErr(c)
		}
	}
	//if password right
	if string(util.HashWithSalt(user.Password, mysqlUser.Salt)) != mysqlUser.Password {
		util.RespNormErr(c, util.WrongPasswordErrCode)
	}

	respLogin := model.RespLogin{
		ID: mysqlUser.ID,
	}

	util.RespNormSuccess(c, respLogin)
}
