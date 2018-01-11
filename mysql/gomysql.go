package gomysql

import (
	"database/sql"
	"goini"
	"fmt"
	_ "mysql-master"
	"errors"
	"reflect"
	"strings"
	"test/tools"
)

type SqlModel struct {
	db 						*sql.DB
	tablename 				string
	columnstr  				string
	where 					string
	whereParam				[]interface{}
	pk 						string
	orderby 				string
	limit 					string
	join 					string
	clear					int64
}

//初始化创建连接池
func InitPool()(this *SqlModel){
	this = new(SqlModel)
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

	this.db.SetMaxOpenConns(11)
	this.db.SetMaxIdleConns(9)
	err = this.db.Ping()

	if err != nil {
		tools.LogError("mysql InitSql error:" + err.Error())
	}
	return this

}

//模型初始化
func (this *SqlModel) Clear()*SqlModel{
	this.tablename = ""
	this.columnstr = "*"
	this.where = " where 1 "
	this.whereParam = make([]interface{},0)
	this.pk = ""
	this.orderby = ""
	this.limit = ""
	this.join = ""
	this.clear = 0
	return this
}

//设置数据表
func (this *SqlModel) Table(tablename string)  *SqlModel{
	if this.clear > 0 {
		this.Clear()
	}
	this.tablename = tablename
	return this
}

/**
设置where条件
接收类型 string、map、struct 三种类型数据
 */
func (this *SqlModel) Where(params interface{},args ... interface{}) *SqlModel {
	if this.clear > 0 {
		this.Clear()
	}
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Map {
		p := params.(map[string]interface{})
		for key,val := range p{
			this.where += " AND `" + key + "`=? "
			this.whereParam = append(this.whereParam,val)
		}
	} else if v.Kind() == reflect.Struct {

	} else if v.Kind() == reflect.String {
		this.where += " AND `" + params.(string) + "`" + args[0].(string) + "? "
		this.whereParam = append(this.whereParam,args[1])
	}
	return this
}

/**
设置查询字段
 */
func (this *SqlModel) Field(columnstr ... string) *SqlModel{
	if this.clear > 0 {
		this.Clear()
	}
	if len(columnstr) == 0 {
		this.columnstr = "*"
	}else {
		this.columnstr = strings.Join(columnstr,",")
	}
	return this
}

/**
设置主键
 */
func (this *SqlModel) SetPk(pk string) *SqlModel {
	if this.clear > 0 {
		this.Clear()
	}
	this.pk = pk
	return this
}

/**
设置排序
 */
func (this *SqlModel) OrderBy(params ... string) *SqlModel {
	if this.clear > 0 {
		this.Clear()
	}
	this.orderby = " ORDER BY "
	this.orderby += strings.Join(params,",")
	return this
}

/**
原生sql排序
 */
func (this *SqlModel) OrderByRaw(orderRawSql string) *SqlModel {
	if this.clear > 0 {
		this.Clear()
	}
	this.orderby = " ORDER BY " + orderRawSql
	return this
}

/**
联表
 */
func (this *SqlModel) Join(table,condition string,joinType ... string) *SqlModel {
	if this.clear > 0 {
		this.Clear()
	}
	if len(joinType) == 0 {
		this.join += fmt.Sprintf(" LEFT JOIN %v ON %v ", table, condition)
		return this
	}
	joinType[0] = strings.ToUpper(joinType[0])
	switch joinType[0] {
	case "LEFT":
		this.join += fmt.Sprintf(" LEFT JOIN %v ON %v ", table, condition)
	case "RIGHT":
		this.join += fmt.Sprintf("RIGHT JOIN %v ON %v", table, condition)
	case "INNER":
		this.join += fmt.Sprintf("INNER JOIN %v ON %v", table, condition)
	case "FULLJOIN":
		this.join += fmt.Sprintf("FULL JOIN %v ON %v", table, condition)
	default:
		this.join += fmt.Sprintf(" LEFT JOIN %v ON %v ", table, condition)
	}
	return this
}

/**
结果条数限制
 */
func (this *SqlModel) Limit(size ...int64) *SqlModel {
	if this.clear > 0 {
		this.Clear()
	}
	var end int64
	start := size[0]
	if len(size) > 1 {
		end = size[1]
		this.limit = fmt.Sprintf(" LIMIT %d,%d ",start,end)
		return this
	}
	this.limit = fmt.Sprintf(" LIMIT %d ",start)
	return this
}

func (this *SqlModel) Create(param map[string] interface{}) (num int64, err error){
	this.clear = 1
	if this.db == nil {
		return 0, errors.New("mysql not connect")
	}
	var keys []interface{}
	var values []interface{}
	var preFlagArr []interface{}

	/*if len(this.pk) != 0 {
		delete(param,this.pk)
	}*/
	for k,v := range param{
		keys = append(keys,k)
		values = append(values,v)
		preFlagArr = append(preFlagArr,"?")
	}
	flag := ","
	columns := tools.SliceToString(keys,flag)
	preFlag := tools.SliceToString(preFlagArr,flag)


	sql := "INSERT INTO " + this.tablename + "(" + columns + ") values(" + preFlag + ")"
	// INSERT INTO table(colume1,colume2,colume3) values(?,?,?)
	num, err = this.Insert(sql,values ...)

	return num, err
}

func (this *SqlModel) Insert(sql string,values ... interface{})(int64, error) {
	this.clear = 1
	if this.db == nil {
		return 0,errors.New("db is Null")
	}

	//sql INSERT INTO table(colume1,colume2,colume3) values(?,?,?)
	dbi, err := this.db.Prepare(sql)
	defer dbi.Close()
	dbi.Close()

	tools.LogInfo(sql)
	tools.LogInfo(values ...)
	if err != nil {
		tools.LogError("Insert error: ", err.Error())
		return 0, err
	}


	res, err := dbi.Exec(values ...)

	if err != nil {
		tools.LogError("Insert error: ", err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()

	return id, err
}

func (this *SqlModel) Save(data map[string] interface{}) (int64, error) {
	this.clear = 1
	//构造sql语句
	// UPDATE table set colume1=?,colume2=?,colume3=? + this.where
	var setFileds string
	var vals []interface{}
	num := 0
	for k,v := range data{
		num++
		setFileds += k + "=? "
		if num != len(data){
			setFileds += ","
		}
		vals = append(vals,v)
	}

	//接上where参数
	for _,v1 := range this.whereParam {
		vals = append(vals,v1)
	}

	sql := "UPDATE " + this.tablename + " SET " + setFileds + this.where
	RowsAffected, err := this.Update(sql,vals ...)

	return RowsAffected, err
}

func (this *SqlModel) Update(sql string, values ...interface{}) (int64, error) {
	this.clear = 1
	if this.db == nil {
		return 0, errors.New("db is Null")
	}
	dbi, err := this.db.Prepare(sql)
	defer dbi.Close()

	tools.LogInfo(sql)
	tools.LogInfo(values ...)
	if err != nil {
		tools.LogError("UPDATE error:", err.Error())
		return 0, err
	}

	res, err := dbi.Exec(values ...)

	if err != nil {
		tools.LogError("UPDATE error:", err.Error())
		return 0, err
	}

	RowsAffected, err := res.RowsAffected()

	return RowsAffected, err
}

func (this *SqlModel)Delete() (int64, error) {
	this.clear = 1
	//构造sql语句
	// DELETE FROM  table  WHERE colume4=? and colume5 = ?

	if this.where == " where 1 "{   //防止整个表删除
		tools.LogError("Delete error:", "delete语句必须带有有效where条件")
		return 0,nil
	}

	sql := "DELETE FROM " + this.tablename + this.where
	RowsAffected, err := this.SqlDelete(sql,this.whereParam ...)

	return RowsAffected, err
}

func (this *SqlModel) SqlDelete(sql string,vals ... interface{}) (int64, error){
	this.clear = 1
	dbi, err := this.db.Prepare(sql)
	defer dbi.Close()

	tools.LogInfo(sql)
	tools.LogInfo(vals ...)
	if err != nil {
		tools.LogError("DELETE error:", err.Error())
		return 0, err
	}

	res, err := dbi.Exec(vals ...)

	if err != nil {
		tools.LogError("DELETE error:", err.Error())
		return 0, err
	}

	RowsAffected, err := res.RowsAffected()

	return RowsAffected, err
}

func(this *SqlModel) FindAll()([]map[string]interface {}, error) {
	this.clear = 1
	sql := "SELECT " + this.columnstr + " FROM " + this.tablename + this.join + this.where + this.orderby + this.limit
	rowslist,err := this.Query(sql,this.whereParam...)
	tools.LogInfo(rowslist)

	return rowslist,err
}

func(this *SqlModel) FindOne()(map[string]interface{}, error) {
	this.clear = 1
	sql := "SELECT " + this.columnstr + " FROM " + this.tablename + this.join + this.where + this.orderby + this.limit
	row,err := this.QueryRow(sql,this.whereParam...)

	return row,err
}

func (this *SqlModel) Query(sql string,vals ... interface{}) ([]map[string]interface {}, error) {
	this.clear = 1
	tools.LogInfo(sql)
	tools.LogInfo(vals ...)
	dbi,err := this.db.Prepare(sql)

	rowslist := make([]map[string] interface{}, 0)
	if err != nil {
		tools.LogInfo(err.Error())
		return rowslist, err
	}

	rows,err := dbi.Query(vals ...) //
	defer dbi.Close()

	if err != nil {
		tools.LogInfo(err.Error())
		return rowslist, err
	}

	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, err := rows.Columns()
	if err != nil {
		return rowslist, err
	}
	valuePtrs := make([]interface{}, len(columns))  //地址
	values := make([]interface{}, len(columns))  //值

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	rowslist = make([]map[string] interface{}, 0)

	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(valuePtrs...)
		record := make(map[string]interface{})
		for i, col := range values {
			if col != nil {
				b,ok := col.([]byte)
				if ok{
					record[columns[i]] = string(b)
				}else{
					record[columns[i]] = col
				}

			}

		}

		rowslist = append(rowslist,record)
	}

	return rowslist, err

}

func (this *SqlModel) QueryRow(sql string,vals ... interface{})(map[string]interface{}, error){
	this.clear = 1
	tools.LogInfo(sql)
	tools.LogInfo(vals ...)
	dbi,err := this.db.Prepare(sql)
	if err != nil {
		tools.LogInfo(err.Error())
		return nil, err
	}

	rows,err := dbi.Query(vals ...) //
	defer dbi.Close()

	if err != nil {
		tools.LogInfo(err.Error())
		return nil, err
	}

	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	valuePtrs := make([]interface{}, len(columns))  //地址
	values := make([]interface{}, len(columns))  //值

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	record := make(map[string]interface{})
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil,err
		}
		for i, col := range values {
			if col != nil {
				b,ok := col.([]byte)
				if ok{
					record[columns[i]] = string(b)
				}else{
					record[columns[i]] = col
				}

			}

		}
		rows.Close()
	}

	return record, err
}

func (this *SqlModel) Count()(int64, error){
	this.clear = 1
	this.columnstr = "count(*)"
	sql := "SELECT " + this.columnstr + " FROM " + this.tablename + this.join + this.where
	dbi,err := this.db.Prepare(sql)
	defer dbi.Close()
	if err != nil {
		tools.LogInfo(err.Error())
		return 0, err
	}

	var count int64
	err = dbi.QueryRow(this.whereParam ...).Scan(&count)

	if err != nil {
		tools.LogInfo(err.Error())
		return 0, err
	}
	return count,err
}

func (this *SqlModel) DbClose() {
	this.db.Close()
}


