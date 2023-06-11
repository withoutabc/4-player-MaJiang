package client

import "regexp"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	ID        int    `json:"id"`
	LoginTime string `json:"login_time"`
	Token     string `json:"token"`
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func IsMatchChu(s string) bool {
	pattern := "^.*打出\\d+,\\d+$"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		// 处理编译错误
		return false
	}
	return regex.MatchString(s)
}
