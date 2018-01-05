package crontab

import (
	"updatetoken/tools"
	"time"
)

type GetToken struct {
	SIGNKEY string
	ACCESS_TOKEN_INVALID float64
	ACCESS_TOKEN_EXPIRES_HINT float64
	COMPONENTA_ACCESS_TOKEN_INVALID float64
}

func (this *GetToken) AttrInit(){
	this.SIGNKEY = "0a5aabec2a2b11e786d30025b3a90ab6"
	this.ACCESS_TOKEN_INVALID = 40001
	this.ACCESS_TOKEN_EXPIRES_HINT = 42001
	this.COMPONENTA_ACCESS_TOKEN_INVALID = 41001
}

/**
获取公众号access_token
 */
func (this *GetToken) GetAccessToken(appId,secret string) (string,error){
	tools.LogInfo("---------------GetAccessToken----------------")
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appId + "&secret=" + secret
	res,err := tools.HttpGet(url)
	if err != nil {
		tools.LogError("GetAccessToken Error:" ,err.Error())
	}
	return res,err
}

func (this *GetToken) GetAppAccessToken(appAppId,appAppSecret string)(string,error){
	tools.LogInfo("---------------GetAppAccessToken----------------")
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appAppId + "&secret=" + appAppSecret
	res,err := tools.HttpGet(url)
	if err != nil {
		tools.LogError("GetAppAccessToken Error:" ,err.Error())
	}
	return res,err
}

/**
加密得到jsapi_ticket
 */
func (this *GetToken) GetJsApiTicket(accessToken,appId,appSecret string) (string,error){
	tools.LogInfo("---------------GetJsApiTicket----------------")
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=" + accessToken
	res,err := tools.HttpGet(url)
	if err != nil {
		tools.LogError("GetJsApiTicket Error:" ,err.Error())
		return "",err
	}

	//请求数据转换成map
	resObj,err := tools.Obj2mapObj(res)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return "",err
	}

	// access_token失效，更新access_token再重新获取jsapi_tiket
	if resObj["errcode"] != nil && (resObj["errcode"].(float64) == this.ACCESS_TOKEN_INVALID || resObj["errcode"].(float64) == this.ACCESS_TOKEN_EXPIRES_HINT) && appId != "" && appSecret != "" {
		accessToken,err = this.updateAccessToken(appId,appSecret)
		if err != nil {
			tools.LogError("updateAccessToken Error:",err.Error())
			return "",err
		}
		url = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=" + accessToken
		res,err = tools.HttpGet(url)
		if err != nil {
			tools.LogError("GetJsApiTicket Error:" ,err.Error())
			return "",err
		}
	}
	return res,err
}

/**
获取卡券apiticket
 */
func (this *GetToken) GetApiTicket(accessToken,appId,appSecret string)(string,error){
	tools.LogInfo("---------------GetApiTicket----------------")
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=wx_card&access_token=" + accessToken
	res,err := tools.HttpGet(url)
	if err != nil {
		tools.LogError("GetApiTicket Error:" ,err.Error())
		return "",err
	}

	//请求数据转换成map
	resObj,err := tools.Obj2mapObj(res)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return "",err
	}

	if resObj["errcode"] != nil && (resObj["errcode"].(float64) == this.ACCESS_TOKEN_INVALID || resObj["errcode"].(float64) == this.ACCESS_TOKEN_EXPIRES_HINT) && appId != "" && appSecret != "" {
		accessToken,err = this.updateAccessToken(appId,appSecret)
		if err != nil {
			tools.LogError("updateAccessToken Error:",err.Error())
			return "",err
		}
		url = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=wx_card&access_token=" + accessToken
		res,err = tools.HttpGet(url)
		if err != nil {
			tools.LogError("GetApiTicket Error:" ,err.Error())
			return "",err
		}
	}
	return res,err
}

/**
更新公众号access_token
 */
func (this *GetToken) updateAccessToken(appId,appSecret string)(string,error) {
	accessToken,err := this.GetAccessToken(appId,appSecret)
	if err != nil {
		tools.LogError("updateAccessToken Error:",err.Error())
		return "",err
	}

	accessTokenObj,err := tools.Obj2mapObj(accessToken)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return "",err
	}

	data := make(map[string]interface{})
	data["access_token"] = accessTokenObj["access_token"]
	data["access_token_expires_time"] = time.Now().Format("2006-01-02 15:04:05")
	_,err = pool.Table("wechat").Where("appid","=",appId).Save(data)
	if err != nil {
		tools.LogError("updateAccessToken Error:",err.Error())
		return "",err
	}

	return data["access_token"].(string),err

}
