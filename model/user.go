package model

import "github.com/dgrijalva/jwt-go"

// User 用户信息
type User struct {
	ID       int    `json:"id" form:"id" binding:"-" gorm:"primarykey" `
	Username string `json:"username" form:"username" binding:"required" gorm:"type:varchar(40);not null"`
	Password string `json:"password" form:"password" binding:"required" gorm:"not null;type:longblob"`
	Salt     []byte `json:"salt" form:"salt" binding:"-" gorm:"not null"`
}

type PlayUser struct {
	ID       int    `json:"id" form:"id" binding:"-" gorm:"primarykey" `
	Username string `json:"username" form:"username" binding:"required" gorm:"type:varchar(40);not null"`
}

type MyClaims struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// RespLogin 登录响应
type RespLogin struct {
	ID int `json:"id"`
}
