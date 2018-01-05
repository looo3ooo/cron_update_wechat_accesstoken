package tools

import (
	"strconv"
	"reflect"
	"encoding/json"
)

/**
 字符串截取
 */
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}



/**
 * 传入切片和分割符，返回分割后的字符串
 * @param slice arr := []{"a","b","c",1,2}
 * @param string flag := ","
 * @return string "a,b,c,1,2"
 */
func SliceToString(slice []interface{},flag string) string{
	s := ""
	l := len(slice)
	for k,v :=range slice{
		if IsString(v){
			s += v.(string)
		}else if IsInt(v){
			str := strconv.Itoa(v.(int))
			s += str
		}
		if k!=(l-1){
			s += flag
		}
	}
	return s
}

/**
 interface类型的slice转换为 slice
 */
func InterfaceToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}


/**
 判断是否为字符串
 */
func IsString(data interface{}) bool{
	if _,ok := data.(string); ok{
		return true
	}else {
		return false
	}
}

/**
 判断是否为int
 */
func IsInt(data interface{}) bool{
	if _,ok := data.(int); ok{
		return true
	}else {
		return false
	}
}


// 函　数：Obj2map
// 概　要：
// 参　数：
//      obj: 传入Obj
// 返回值：
//      mapObj: map对象
//      err: 错误
func Obj2mapObj(obj string) (mapObj map[string]interface{}, err error) {
	// 结构体转json
	/*b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}*/

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(obj), &result); err != nil {
		return nil, err
	}
	return result, nil
}

