package tools

import (
	"net/http"
	"io/ioutil"
	"strings"
)

func HttpGet(url string)(string, error) {
	LogInfo("httpGet:", url)
	res, err := http.Get(url)
	if err != nil {
		LogError("httpGet Error:", err.Error())
		return "0",err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		LogError("httpGet ioutil Error:", err.Error())
		return "0",err
	}
	LogInfo(string(body))

	return string(body),err

}


func HttpPost(url string, reqbody string) (string, error) {
	LogInfo("HttpPost:", url)
	LogInfo("HttpPost:", reqbody)
	postReq, err := http.NewRequest("POST",
		url, //post链接
		strings.NewReader(reqbody)) //post内容

	if err != nil {
		LogError("HttpPost Error:", err.Error())
		return "0", err
	}

	//增加header
	//postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	var body []byte
	//执行请求
	httpClient := &http.Client{}
	resp, err := httpClient.Do(postReq)
	if err != nil {
		LogError("HttpPost Error:", err.Error())
		return "0", err
	} else {
		//读取响应
		body, err = ioutil.ReadAll(resp.Body) //此处可增加输入过滤
		if err != nil {
			LogError("POST请求:读取body失败:", err.Error())
			return "0", err
		}
		LogInfo("POST请求:创建成功:", string(body))
	}
	defer resp.Body.Close()

	return string(body), nil
}
