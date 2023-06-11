package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"majiang/client"
	"majiang/dao"
	"majiang/model"
	"net/http"
)

func main() {
	dao.InitDB()
	client.InitWsSocket()
	//用于接收消息的死循环
	var forever chan struct{}
	go client.WaitResponse()
	var (
		username string
		password string
	)

	fmt.Printf("请输入用户名\n")
	fmt.Scan(&username)
	fmt.Printf("请输入密码\n")
	fmt.Scan(&password)
	resp, err := Login(username, password)
	if err != nil {
		log.Println(err)
		log.Println("不想处理错误，游戏结束")

	}
	resp.ID = 1
	log.Printf("登录成功\n")
	//不做登录验证
	log.Printf("1.创建房间\n2.加入房间\n3.查看房间列表\n")
	var x int
	fmt.Scan(&x)
	//	common = iota
	//	create
	//	join
	//	leave
	//	changeReady

	switch x {
	case 1:
		data := map[string]interface{}{
			"userID": resp.ID,
		}
		err = client.SendMessage(1, data)
		if err != nil {
			log.Println("不想处理错误，游戏结束")
			return
		}

	case 2:
	}
	<-forever
}

func Login(username string, password string) (model.RespLogin, error) {
	user := client.User{
		Username: username,
		Password: password,
	}
	log.Println(user)
	jsonStr, err := json.Marshal(user)
	log.Println(jsonStr)
	if err != nil {
		return model.RespLogin{}, err
	}
	req, err := http.NewRequest("POST", "http://localhost:2022/user/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		return model.RespLogin{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return model.RespLogin{}, err
	}

	defer resp.Body.Close()
	var response model.RespLogin
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return model.RespLogin{}, err
	}
	log.Println(response)
	return model.RespLogin{}, err
}
