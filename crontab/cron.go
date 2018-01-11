package crontab

import "updatetoken/mysql"

type Cron struct {
	AutoUpdateToken AutoUpdateToken
}
var pool *gomysql.SqlModel

func InitModel() *gomysql.SqlModel{
	pool = gomysql.InitPool()
	pool.Clear()
	return pool
}
