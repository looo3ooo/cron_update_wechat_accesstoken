package main

import (
	"updatetoken/tools"
	"updatetoken/model"
	"github.com/robfig/cron"
	"updatetoken/crontab"
)

func main() {

	//开启日志
	defer tools.LogFlush()
	tools.InitLog()

	//初始化mysql连接
	model.InitModel()
	crontab.InitModel()

	c := cron.New()
	spec := "*/3 * * * * *"
	c.AddFunc(spec, func() {
		tools.LogInfo("cron running autoUpdateToken:")
		updateTokenController := new(crontab.AutoUpdateToken)
		updateTokenController.Index()
	})


	c.Start()

	select{}//阻塞主线程不退出
}
