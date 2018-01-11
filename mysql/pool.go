package gomysql

import (
	"fmt"
	"database/sql"
	"goini"
	"test/tools"
)

type Pool struct {
	db 						*sql.DB
}

//初始化创建连接池
func InitPool()(*Pool){
	this := new(Pool)
	var err error
	ConfigCentor := goini.SetConfig("./config/config.ini")
	ip := ConfigCentor.GetValue("mysql", "ip")
	uid := ConfigCentor.GetValue("mysql", "uid")
	pwd := ConfigCentor.GetValue("mysql", "pwd")
	dbname := ConfigCentor.GetValue("mysql", "databasename")
	data_str := fmt.Sprintf("%s:%s@(%s:3306)/%s?charset=utf8", uid, pwd, ip, dbname)
	tools.LogInfo("-----数据库连接----" + data_str)
	this.db, err = sql.Open("mysql", data_str)

	if err != nil {
		tools.LogError("mysql InitSql error:" + err.Error())
	}

	this.db.SetMaxOpenConns(30)
	this.db.SetMaxIdleConns(20)
	err = this.db.Ping()

	if err != nil {
		tools.LogError("mysql InitSql error:" + err.Error())
	}
	return this

}

func (this *Pool) DbClose() {
	this.db.Close()
}

func (this *Pool) NewModel() *SqlModel{
	NewSqlModel := new(SqlModel)
	NewSqlModel.db = this.db
	NewSqlModel.Clear()
	return NewSqlModel
}
