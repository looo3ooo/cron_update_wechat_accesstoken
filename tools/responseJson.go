package tools

import (
	"net/http"
	"encoding/json"
	"fmt"
	 "updatetoken/res"
)

/**
 返回json类型数据
 */
func ResponseJson(w http.ResponseWriter,res response.BaseJsonBean){
	jsonData,_ := json.Marshal(res)
	fmt.Fprintln(w,string(jsonData))
}
