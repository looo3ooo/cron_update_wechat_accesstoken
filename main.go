package main

import (
	"updatetoken/tools"
	"github.com/robfig/cron"
	"updatetoken/crontab"
	"goini"
	"fmt"
)

func main() {

	//开启日志
	defer tools.LogFlush()
	tools.InitLog()

	pool := crontab.InitModel()
	defer pool.DbClose()

	c := cron.New()
	ConfigCentor := goini.SetConfig("./config/config.ini")
	spec := ConfigCentor.GetValue("cron", "spec")
	fmt.Println(spec)
	c.AddFunc(spec, func() {
		tools.LogInfo("cron running autoUpdateToken:")
		updateTokenController := new(crontab.AutoUpdateToken)
		updateTokenController = updateTokenController.AttrInit()
		updateTokenController.Index()
	})


	c.Start()

	select{}//阻塞主线程不退出
}
