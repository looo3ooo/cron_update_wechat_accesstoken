package model

import "updatetoken/mysql"

type Model struct {
	Wechat Wechat
}

var pool *gomysql.SqlModel

func InitModel(){
	pool = gomysql.InitPool()
	pool.InitModel()
}
