package crontab

import "updatetoken/mysql"

type Cron struct {
	AutoUpdateToken AutoUpdateToken
}
var pool *gomysql.SqlModel

func InitModel(){
	pool = gomysql.InitPool()
	pool.InitModel()
}
