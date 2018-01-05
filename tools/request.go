package tools

import (
	"net/http"
	"io/ioutil"
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
