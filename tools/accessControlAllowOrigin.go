package tools

import "net/http"

func AccessControlAllowOrigin(w http.ResponseWriter)  {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//w.Header().Set("content-type", "application/json")             //返回数据格式是json
}
