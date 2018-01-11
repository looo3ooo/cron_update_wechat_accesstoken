package crontab

import "updatetoken/mysql"

type Cron struct {
	AutoUpdateToken AutoUpdateToken
}
var pool *gomysql.Pool

func InitModel() *gomysql.Pool{
	pool = gomysql.InitPool()
	return pool
}
