package util

import (
	"errors"
)

var ErrorCodeMap = map[int]error{
	BlankAuthErrCode:       BlankAuthErr,
	WrongAuthFormatErrCode: WrongAuthFormatErr,
	InValidTokenErrCode:    InvalidTokenErr,

	UnauthorizedErrCode:  UnauthorizedErr,
	InternalServeErrCode: InternalServeErr,

	IdNotIntegral:             IdNotIntegralErr,
	NoRecordErrCode:           NoRecordErr,
	RepeatedUsernameErrCode:   RepeatedUsernameErr,
	WrongPasswordErrCode:      WrongPasswordErr,
	UpdateFailErrCode:         UpdateFailErr,
	RepeatedTitleErrCode:      RepeatedTitleErr,
	RepeatedSubmissionErrCode: RepeatedSubmissionErr,

	BindingQueryErrCode: BindingQueryErr,

	//room
	RoomNotExistErrCode: RoomNotExistErr,
}

var (
	BlankAuthErr       = errors.New("请求头中auth为空")
	WrongAuthFormatErr = errors.New("请求头中auth格式有误")
	InvalidTokenErr    = errors.New("无效的Token")

	UnauthorizedErr  = errors.New("unauthorized")
	InternalServeErr = errors.New("internal serve error")

	IdNotIntegralErr      = errors.New("id is not a integral")
	NoRecordErr           = errors.New("no record")
	RepeatedUsernameErr   = errors.New("repeated username")
	WrongPasswordErr      = errors.New("wrong password")
	UpdateFailErr         = errors.New("update failed")
	RepeatedTitleErr      = errors.New("repeated title")
	RepeatedSubmissionErr = errors.New("repeat submit")

	BindingQueryErr = errors.New("binding error")

	RoomNotExistErr = errors.New("room not exist")
)

const (
	NoErrCode = iota

	BlankAuthErrCode //1
	WrongAuthFormatErrCode
	InValidTokenErrCode

	UnauthorizedErrCode //4
	InternalServeErrCode

	IdNotIntegral //6
	NoRecordErrCode
	RepeatedUsernameErrCode
	WrongPasswordErrCode
	UpdateFailErrCode //update or delete error   10
	RepeatedTitleErrCode
	RepeatedSubmissionErrCode

	BindingQueryErrCode //13
	//room

	RoomNotExistErrCode //14
)
