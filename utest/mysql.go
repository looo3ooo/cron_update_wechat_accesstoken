package utest

import (
	"fmt"
	"updatetoken/mysql"
	"updatetoken/tools"
	"updatetoken/res"
	"encoding/json"
)


var pool *gomysql.SqlModel

func InitModel(){
	pool = gomysql.InitPool()
	defer pool.DbClose()
	pool.Clear()
}

func MysqlCreateTest(){
	Users := make(map[string]interface{})
	Users["name"] = "test"
	Users["email"] = "test@qq.com"
	Users["password"] = "$10$iflHbfXNyE1SeW3MsROMt.PuCHOuW7IclIpaFtRMo8RrTS7Qi7HMa"
	Users["activated"] = 1
	insert_id,err := pool.Table("users").Create(Users)

	if err != nil {
		tools.LogError("MysqlInsertTest Error:",err.Error())
	}
	fmt.Println(insert_id)

}


func MysqlInsertTest()  {
	var Inserts []interface{}
	Inserts = append(Inserts,"test")
	Inserts = append(Inserts,"test1@qq.com")
	Inserts = append(Inserts,"$10$iflHbfXNyE1SeW3MsROMt.PuCHOuW7IclIpaFtRMo8RrTS7Qi7HMa")
	Inserts = append(Inserts,1)
	sql := "insert into users(name,email,password,activated) values(?,?,?,?)"
	insert_id,err := pool.Insert(sql,Inserts ...)

	if err != nil {
		tools.LogError("MysqlInsertTest Error:",err.Error())
	}
	fmt.Println(insert_id)
}

func MysqlSaveTest(){
	Data := make(map[string]interface{})
	Data["password"] = "$10$576e.hZNopkmpFrpYHZRoeidgXx331ZB8WqPZ9rJj6UXbpCCCcCHe"
	Data["activated"] = 1

Where := make(map[string]interface{})
	Where["name"]  = "test"
	//affected_rows,err := pool.InitModel().Table("users").Where(Where).Save(Data)


	affected_rows,err := pool.Table("users").Where("name","=","test").Where("email","<>","test@qq.com").Save(Data)

	if err != nil {
		tools.LogError("MysqlSaveTest Error:",err.Error())
	}
	fmt.Println(affected_rows)
}

func MysqlUpdateTest(){
	var Updates []interface{}
	Updates = append(Updates,"$10$iflHbfXNyE1SeW3MsROMt.PuCHOuW7IclIpaFtRMo8RrTS7Qi7HMa1")
	Updates = append(Updates,"1")
	Updates = append(Updates,"test")


	sql := "UPDATE users set password=?,activated=? where name=?"
	affected_rows,err := pool.Update(sql,Updates ...)
	if err != nil {
		tools.LogError("MysqlSaveTest Error:",err.Error())
	}
	fmt.Println(affected_rows)
}

func MysqlDeleteTest() {
	Users := make(map[string]interface{})
	Users["name"] = "test"
	Users["email"] = "test@qq.com"
	Users["password"] = "$10$iflHbfXNyE1SeW3MsROMt.PuCHOuW7IclIpaFtRMo8RrTS7Qi7HMa"
	Users["activated"] = 1
	affected_rows,err := pool.Table("users").Where(Users).Delete()

	if err != nil {
		tools.LogError("MysqlDeleteTest Error:",err.Error())
	}
	fmt.Println(affected_rows)
}


func MysqlSqlDeleteTest()  {
	sql := "DELETE FROM users where name=? and email=?"
	var vals []interface{}
	vals = append(vals,"test")
	vals = append(vals,"test1@qq.com")
	affected_rows,err := pool.Table("users").SqlDelete(sql,vals...)

	if err != nil {
		tools.LogError("MysqlDeleteTest Error:",err.Error())
	}
	fmt.Println(affected_rows)
}

func MysqlFindAllTest()  {
	var orderBy []string
	orderBy = append(orderBy,"created_at desc")
	orderBy = append(orderBy,"id")
	rows_list, err := pool.Table("users").Where("activated","=",0).OrderBy(orderBy...).Limit(1,1).FindAll()
	res := response.NewBaseJsonBean()
	if err != nil {
		res.Code = 99
		res.Data = rows_list
		res.Message = "mysql error: " + err.Error()
		jsonData,_ := json.Marshal(res)
		fmt.Println(string(jsonData))
		return
	}

	for _,v := range rows_list{  //slice
		fmt.Println(v)
		fmt.Println(v["name"])
	}

	res.Code = 100
	res.Data = rows_list
	res.Message = "查询成功"
	jsonData,_ := json.Marshal(res)
	fmt.Println(string(jsonData))
}

func MysqlQueryTest(){
	sql := "select id,name,email,activated from users where activated=? limit 3"
	var vals []interface{}
	vals = append(vals,1)
	rows_list, err := pool.Query(sql,vals ...)

	res := response.NewBaseJsonBean()
	if err != nil {
		res.Code = 99
		res.Data = rows_list
		res.Message = "mysql error: " + err.Error()
		jsonData,_ := json.Marshal(res)
		fmt.Println(string(jsonData))
		return
	}

	for _,v := range rows_list{  //slice
		fmt.Println(v)
		fmt.Println(v["name"])
	}

	res.Code = 100
	res.Data = rows_list
	res.Message = "查询成功"
	jsonData,_ := json.Marshal(res)
	fmt.Println(string(jsonData))

}

func MysqlFindOneTest(){
	orderBy := make(map [string]string)
	orderBy["created_at"] = "desc"
	orderBy["id"] = "asc"
	row, err := pool.Table("users").Field("id","name","email").Where("activated","=",1).OrderByRaw("created_at desc,id asc").FindOne()

	res := response.NewBaseJsonBean()
	if err != nil{
		res.Code = 100
		res.Data = row
		res.Message = "mysql error:" + err.Error()
		jsonData,_ := json.Marshal(res)
		fmt.Println(string(jsonData))
		return
	}


	res.Code = 100
	res.Data = row
	res.Message = "查询成功"
	jsonData,_ := json.Marshal(res)
	fmt.Println(string(jsonData))
}


func MysqlQueryRowTest(){
	sql := "select id,name,email,activated from users where activated=? limit 3"
	var vals []interface{}
	vals = append(vals,1)
	row, err := pool.QueryRow(sql,vals ...)

	res := response.NewBaseJsonBean()
	if err != nil{
		res.Code = 100
		res.Data = row
		res.Message = "mysql error:" + err.Error()
		jsonData,_ := json.Marshal(res)
		fmt.Println(string(jsonData))
		return
	}


	res.Code = 100
	res.Data = row
	res.Message = "查询成功"
	jsonData,_ := json.Marshal(res)
	fmt.Println(string(jsonData))
}


