package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	secKillConf = &SecSkillConf{
		secProductInfo: make(map[int]*secProductInfoConf, 1024),
	}
)

type RedisConf struct {
	redisAddr        string
	redisMaxIdle     int
	redisMaxActive   int
	redisIdleTimeout int
}

type EtcdConf struct {
	etcdAddr   string
	etcdSecKey string
	timeout    int
}

type SecSkillConf struct {
	redisConf      RedisConf
	etcdConf       EtcdConf
	logPath        string
	logLevel       string
	secProductInfo map[int]*secProductInfoConf
}

type SecInfoConf struct {
	ProductId int
	StartTime int
	EndTime   int
	Status    int
	Iotal     int
	Left      int
}

func initConfig() (err error) {

	redisAddr := beego.AppConfig.String("redis_addr")
	etcdAddr := beego.AppConfig.String("etcd_addr")
	logs.Debug("read config succ, redis addr:%v", redisAddr)
	logs.Debug("read config succ, etcd addr:%v", etcdAddr)

	secKillConf.etcdConf.etcdAddr = etcdAddr
	secKillConf.redisConf.redisAddr = redisAddr

	if len(redisAddr) == 0 || len(etcdAddr) == 0 {
		err = fmt.Errorf("init config faild, redis[%s] or etcd[%s] config is null", redisAddr, etcdAddr)
		println(redisAddr, "-----", etcdAddr)
		return err
	}

	redisMaxIdle, err := beego.AppConfig.Int("redis_max_idle")
	if err != nil {
		err = fmt.Errorf("init config faild, redisMaxIdle[%s] config is null", redisMaxIdle)

		return err
	}

	redisMaxActive, err := beego.AppConfig.Int("redis_max_active")
	if err != nil {
		err = fmt.Errorf("init config faild, redisMaxActive[%s] config is null", redisMaxActive)

		return err
	}

	redisIdleTimeout, err := beego.AppConfig.Int("redis_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config faild, redisIdTimeout[%s] config is null", redisIdleTimeout)

		return err
	}

	secKillConf.redisConf.redisMaxIdle = redisMaxIdle
	secKillConf.redisConf.redisMaxActive = redisMaxActive
	secKillConf.redisConf.redisIdleTimeout = redisIdleTimeout

	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		err = fmt.Errorf("init config")
		return
	}
	secKillConf.etcdConf.timeout = etcdTimeout
	secKillConf.etcdConf.etcdSecKey = beego.AppConfig.String("etcd_sec")

	secKillConf.logPath = beego.AppConfig.String("log_path")
	secKillConf.logLevel = beego.AppConfig.String("log_level")
	return
}
