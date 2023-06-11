package main

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"majiang/api"
	"majiang/client"
	"majiang/dao"
	"majiang/model"
	"majiang/service"
	"majiang/util"
	"time"
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
	for {
		fmt.Printf("请输入用户名\n")
		fmt.Scan(&username)
		fmt.Printf("请输入密码\n")
		fmt.Scan(&password)
		resp, err := Login(username, password)
		if err != nil {
			continue
		}
		fmt.Printf("登录成功\n")
		//不做登录验证
	start:
		for {
			time.Sleep(1 * time.Second)
			fmt.Printf("1.创建房间\n2.加入房间\n3.查看房间列表\n")
			var x int
			fmt.Scan(&x)
			switch x {
			case 1:
				data := map[string]interface{}{
					"userID": resp.ID,
				}
				err = client.SendMessage(int(api.Create), data)
				if err != nil {
					continue
				}
				for {
					time.Sleep(1 * time.Second)
					var action int
					fmt.Printf("1.准备\n2.离开房间\n3.发言\n")
					fmt.Scan(&action)
					//准备
					if action == 1 {
						roomID := FindRoomID(resp.ID)
						data = map[string]interface{}{
							"userID": resp.ID,
							"roomID": roomID,
							"ready":  true,
						}
						err = client.SendMessage(int(api.ChangeReady), data)
						if err != nil {
							continue
						}
					} else if action == 2 { //离开
						roomID := FindRoomID(resp.ID)
						data = map[string]interface{}{
							"userID": resp.ID,
							"roomID": roomID,
						}
						err = client.SendMessage(int(api.Leave), data)
						if err != nil {
							continue
						}
						continue start
					} else if action == 3 {
						fmt.Printf("输入你想说的话\n")
						var sentence string
						fmt.Scan(&sentence)
						roomID := FindRoomID(resp.ID)
						data = map[string]interface{}{
							"userID":   resp.ID,
							"roomID":   roomID,
							"sentence": sentence,
						}
						err = client.SendMessage(int(api.Common), data)
						if err != nil {
							continue
						}

					} else {
						continue
					}
				}
			case 2:
				fmt.Println("请输入要加入的房间号(输入0返回上一步)")
				var roomID string
				fmt.Scan(&roomID)
				if roomID == "0" {
					continue
				}
				data := map[string]interface{}{
					"userID": resp.ID,
					"roomID": roomID,
				}
				err = client.SendMessage(2, data)
				if err != nil {
					continue
				}
				for {

				}
			case 3:
				err = client.SendMessage(5, nil)
				if err != nil {
					continue
				}
				continue start
			}
		}

	}

	<-forever
}

func Login(username string, password string) (model.RespLogin, error) {
	user := client.User{
		Username: username,
		Password: password,
	}
	mysqlUser, err := service.SearchUserByName(user.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.RespLogin{}, err
		} else {
			return model.RespLogin{}, err
		}
	}
	//if password right
	if string(util.HashWithSalt(user.Password, mysqlUser.Salt)) != mysqlUser.Password {
		return model.RespLogin{}, err
	}
	respLogin := model.RespLogin{
		ID: mysqlUser.ID,
	}
	return respLogin, nil
}

func FindRoomID(UID int) string {
	log.Println(api.Rooms)
	for roomID, room := range api.Rooms {
		for userID := range room.Users {
			if userID == UID {
				return roomID
			}
		}
	}
	return ""
}
