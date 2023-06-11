package main

import (
	"majiang/api"
	"majiang/dao"
)

func main() {
	dao.InitDB()
	api.InitRouter()
}
