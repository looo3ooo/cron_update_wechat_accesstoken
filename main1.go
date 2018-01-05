package main

import (
	"test/tools"
	"net/http"
	"test/routers"
	"test/utest"
	"test/model"
)

func main()  {
	//开启日志
	defer tools.LogFlush()
	tools.InitLog()

	utest.InitModel()
	//utest.MysqlCreateTest()
	//utest.MysqlInsertTest()
	//utest.MysqlSaveTest()
	//utest.MysqlUpdateTest()
	//utest.MysqlDeleteTest()
	//utest.MysqlSqlDeleteTest()
	//utest.MysqlFindAllTest()
	//utest.MysqlQueryTest()
	//utest.MysqlQueryRowTest()
	//utest.MysqlFindOneTest()

	model.InitModel()

	//读取路由
	routers.Init()
	err := http.ListenAndServe(":8021",nil)
	if err != nil {
		tools.LogError("ListenAndServe Error:" + err.Error())
	}
}


