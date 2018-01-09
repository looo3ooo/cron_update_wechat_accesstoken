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
}

//初始化创建连接池
func InitPool()(m *SqlModel){
	m = new(SqlModel)
	var err error
	ConfigCentor := goini.SetConfig("./config/config.ini")
	ip := ConfigCentor.GetValue("mysql", "ip")
	uid := ConfigCentor.GetValue("mysql", "uid")
	pwd := ConfigCentor.GetValue("mysql", "pwd")
	dbname := ConfigCentor.GetValue("mysql", "databasename")
	data_str := fmt.Sprintf("%s:%s@(%s:3306)/%s?charset=utf8", uid, pwd, ip, dbname)
	tools.LogInfo("-----数据库连接----" + data_str)
	m.db, err = sql.Open("mysql", data_str)

	if err != nil {
		tools.LogError("mysql InitSql error:" + err.Error())
	}

	m.db.SetMaxOpenConns(200)
	m.db.SetMaxIdleConns(100)
	err = m.db.Ping()

	if err != nil {
		tools.LogError("mysql InitSql error:" + err.Error())
	}
	return m

}

//模型初始化
func (m *SqlModel) InitModel()*SqlModel{
	m.tablename = ""
	m.columnstr = "*"
	m.where = " where 1 "
	m.whereParam = make([]interface{},0)
	m.pk = ""
	m.orderby = ""
	m.limit = ""
	m.join = ""
	return m
}

//设置数据表
func (m *SqlModel) Table(tablename string)  *SqlModel{
	m.tablename = tablename
	return m
}

/**
设置where条件
接收类型 string、map、struct 三种类型数据
 */
func (m *SqlModel) Where(params interface{},args ... interface{}) *SqlModel {
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Map {
		p := params.(map[string]interface{})
		for key,val := range p{
			m.where += " AND `" + key + "`=? "
			m.whereParam = append(m.whereParam,val)
		}
	} else if v.Kind() == reflect.Struct {

	} else if v.Kind() == reflect.String {
		m.where += " AND `" + params.(string) + "`" + args[0].(string) + "? "
		m.whereParam = append(m.whereParam,args[1])
	}
	return m
}

/**
设置查询字段
 */
func (m *SqlModel) Field(columnstr ... string) *SqlModel{
	if len(columnstr) == 0 {
		m.columnstr = "*"
	}else {
		m.columnstr = strings.Join(columnstr,",")
	}
	return m
}

/**
设置主键
 */
func (m *SqlModel) SetPk(pk string) *SqlModel {
	m.pk = pk
	return m
}

/**
设置排序
 */
func (m *SqlModel) OrderBy(params ... string) *SqlModel {
	m.orderby = " ORDER BY "
	m.orderby += strings.Join(params,",")
	return m
}

/**
原生sql排序
 */
func (m *SqlModel) OrderByRaw(orderRawSql string) *SqlModel {
	m.orderby = " ORDER BY " + orderRawSql
	return m
}

/**
联表
 */
func (m *SqlModel) Join(table,condition string,joinType ... string) *SqlModel {
	if len(joinType) == 0 {
		m.join += fmt.Sprintf(" LEFT JOIN %v ON %v ", table, condition)
		return m
	}
	joinType[0] = strings.ToUpper(joinType[0])
	switch joinType[0] {
	case "LEFT":
		m.join += fmt.Sprintf(" LEFT JOIN %v ON %v ", table, condition)
	case "RIGHT":
		m.join += fmt.Sprintf("RIGHT JOIN %v ON %v", table, condition)
	case "INNER":
		m.join += fmt.Sprintf("INNER JOIN %v ON %v", table, condition)
	case "FULLJOIN":
		m.join += fmt.Sprintf("FULL JOIN %v ON %v", table, condition)
	default:
		m.join += fmt.Sprintf(" LEFT JOIN %v ON %v ", table, condition)
	}
	return m
}

/**
结果条数限制
 */
func (m *SqlModel) Limit(size ...int64) *SqlModel {
	var end int64
	start := size[0]
	if len(size) > 1 {
		end = size[1]
		m.limit = fmt.Sprintf(" LIMIT %d,%d ",start,end)
		return m
	}
	m.limit = fmt.Sprintf(" LIMIT %d ",start)
	return m
}

func (m *SqlModel) Create(param map[string] interface{}) (num int64, err error){
	if m.db == nil {
		m.InitModel()
		return 0, errors.New("mysql not connect")
	}
	var keys []interface{}
	var values []interface{}
	var preFlagArr []interface{}

	/*if len(m.pk) != 0 {
		delete(param,m.pk)
	}*/
	for k,v := range param{
		keys = append(keys,k)
		values = append(values,v)
		preFlagArr = append(preFlagArr,"?")
	}
	flag := ","
	columns := tools.SliceToString(keys,flag)
	preFlag := tools.SliceToString(preFlagArr,flag)


	sql := "INSERT INTO " + m.tablename + "(" + columns + ") values(" + preFlag + ")"
	// INSERT INTO table(colume1,colume2,colume3) values(?,?,?)
	num, err = m.Insert(sql,values ...)

	m.InitModel()
	return num, err
}

func (m *SqlModel) Insert(sql string,values ... interface{})(int64, error) {
	if m.db == nil {
		m.InitModel()
		return 0,errors.New("db is Null")
	}

	//sql INSERT INTO table(colume1,colume2,colume3) values(?,?,?)
	dbi, err := m.db.Prepare(sql)
	defer dbi.Close()
	dbi.Close()

	tools.LogInfo(sql)
	tools.LogInfo(values ...)
	if err != nil {
		tools.LogError("Insert error: ", err.Error())
		m.InitModel()
		return 0, err
	}


	res, err := dbi.Exec(values ...)

	if err != nil {
		tools.LogError("Insert error: ", err.Error())
		m.InitModel()
		return 0, err
	}

	id, err := res.LastInsertId()

	m.InitModel()
	return id, err
}

func (m *SqlModel) Save(data map[string] interface{}) (int64, error) {
	//构造sql语句
	// UPDATE table set colume1=?,colume2=?,colume3=? + m.where
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
	for _,v1 := range m.whereParam {
		vals = append(vals,v1)
	}

	sql := "UPDATE " + m.tablename + " SET " + setFileds + m.where
	RowsAffected, err := m.Update(sql,vals ...)

	m.InitModel()
	return RowsAffected, err
}

func (m *SqlModel) Update(sql string, values ...interface{}) (int64, error) {
	if m.db == nil {
		m.InitModel()
		return 0, errors.New("db is Null")
	}
	dbi, err := m.db.Prepare(sql)
	defer dbi.Close()

	tools.LogInfo(sql)
	tools.LogInfo(values ...)
	if err != nil {
		tools.LogError("UPDATE error:", err.Error())
		m.InitModel()
		return 0, err
	}

	res, err := dbi.Exec(values ...)

	if err != nil {
		tools.LogError("UPDATE error:", err.Error())
		m.InitModel()
		return 0, err
	}

	RowsAffected, err := res.RowsAffected()

	m.InitModel()
	return RowsAffected, err
}

func (m *SqlModel)Delete() (int64, error) {
	//构造sql语句
	// DELETE FROM  table  WHERE colume4=? and colume5 = ?

	if m.where == " where 1 "{   //防止整个表删除
		tools.LogError("Delete error:", "delete语句必须带有有效where条件")
		m.InitModel()
		return 0,nil
	}

	sql := "DELETE FROM " + m.tablename + m.where
	RowsAffected, err := m.SqlDelete(sql,m.whereParam ...)

	m.InitModel()
	return RowsAffected, err
}

func (m *SqlModel) SqlDelete(sql string,vals ... interface{}) (int64, error){
	dbi, err := m.db.Prepare(sql)
	defer dbi.Close()

	tools.LogInfo(sql)
	tools.LogInfo(vals ...)
	if err != nil {
		tools.LogError("DELETE error:", err.Error())
		m.InitModel()
		return 0, err
	}

	res, err := dbi.Exec(vals ...)

	if err != nil {
		tools.LogError("DELETE error:", err.Error())
		m.InitModel()
		return 0, err
	}

	RowsAffected, err := res.RowsAffected()

	m.InitModel()
	return RowsAffected, err
}

func(m *SqlModel) FindAll()([]map[string]interface {}, error) {
	sql := "SELECT " + m.columnstr + " FROM " + m.tablename + m.join + m.where + m.orderby + m.limit
	rowslist,err := m.Query(sql,m.whereParam...)

	m.InitModel()
	return rowslist,err
}

func(m *SqlModel) FindOne()(map[string]interface{}, error) {
	sql := "SELECT " + m.columnstr + " FROM " + m.tablename + m.join + m.where + m.orderby + m.limit
	row,err := m.QueryRow(sql,m.whereParam...)

	m.InitModel()
	return row,err
}

func (m *SqlModel) Query(sql string,vals ... interface{}) ([]map[string]interface {}, error) {
	tools.LogInfo(sql)
	tools.LogInfo(vals ...)
	dbi,err := m.db.Prepare(sql)
	if err != nil {
		tools.LogInfo(err.Error())
		m.InitModel()
		return nil, err
	}

	rows,err := dbi.Query(vals ...) //
	defer dbi.Close()

	if err != nil {
		tools.LogInfo(err.Error())
		m.InitModel()
		return nil, err
	}

	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, err := rows.Columns()
	if err != nil {
		m.InitModel()
		return nil, err
	}
	valuePtrs := make([]interface{}, len(columns))  //地址
	values := make([]interface{}, len(columns))  //值

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	rowslist := make([]map[string] interface{}, 0)

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

	m.InitModel()
	return rowslist, err

}

func (m *SqlModel) QueryRow(sql string,vals ... interface{})(map[string]interface{}, error){
	tools.LogInfo(sql)
	tools.LogInfo(vals ...)
	dbi,err := m.db.Prepare(sql)
	if err != nil {
		tools.LogInfo(err.Error())
		m.InitModel()
		return nil, err
	}

	rows,err := dbi.Query(vals ...) //
	defer dbi.Close()

	if err != nil {
		tools.LogInfo(err.Error())
		m.InitModel()
		return nil, err
	}

	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, err := rows.Columns()
	if err != nil {
		m.InitModel()
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
			m.InitModel()
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
		break
	}

	m.InitModel()
	return record, err
}

func (m *SqlModel) Count()(int64, error){
	m.columnstr = "count(*)"
	sql := "SELECT " + m.columnstr + " FROM " + m.tablename + m.join + m.where
	dbi,err := m.db.Prepare(sql)
	defer dbi.Close()
	if err != nil {
		tools.LogInfo(err.Error())
		m.InitModel()
		return 0, err
	}

	var count int64
	err = dbi.QueryRow(m.whereParam ...).Scan(&count)

	if err != nil {
		tools.LogInfo(err.Error())
		m.InitModel()
		return 0, err
	}
	m.InitModel()
	return count,err
}

func (m *SqlModel) DbClose() {
	m.db.Close()
}


