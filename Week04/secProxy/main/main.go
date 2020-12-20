package main

import (
	_ "awesomeProject/seckillSystem/secProxy/router"

	"github.com/astaxie/beego"
)

const (
	EtcdKey = "/oldboy/backend/secKill/product"
)

type SecInfoConf struct {
	ProductId int
	StartTime int
	EndTime   int
	Status    int
	Iotal     int
	Left      int
}

func SetLogConfToEtcd() {

}

func main() {
	err := initConfig()
	if err != nil {
		println(err)
		panic("initConfig faild")
		return
	}

	err = initSec()
	if err != nil {
		panic("initSec faild")
		return
	}

	beego.Run()
}
