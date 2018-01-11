package crontab

import (
	"updatetoken/tools"
	"time"
	"goini"
	"encoding/json"
	"updatetoken/mysql"
)

type GetToken struct {
	SIGNKEY string
	ACCESS_TOKEN_INVALID float64
	ACCESS_TOKEN_EXPIRES_HINT float64
	COMPONENTA_ACCESS_TOKEN_INVALID float64
	Model *gomysql.SqlModel
}

func (this *GetToken) AttrInit() *GetToken{
	ConfigCentor := goini.SetConfig("./config/config.ini")
	signkey := ConfigCentor.GetValue("gettoken", "signkey")
	this.SIGNKEY = signkey
	this.ACCESS_TOKEN_INVALID = 41001
	this.ACCESS_TOKEN_EXPIRES_HINT = 42001
	this.COMPONENTA_ACCESS_TOKEN_INVALID = 41001

	this.Model = pool.NewModel()
	return this
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

	if accessTokenObj["access_token"] != nil {
		data := make(map[string]interface{})
		data["access_token"] = accessTokenObj["access_token"]
		data["access_token_expires_time"] = time.Now().Format("2006-01-02 15:04:05")
		_,err = this.Model.Table("wechat").Where("appid","=",appId).Save(data)
		if err != nil {
			tools.LogError("updateAccessToken Error:",err.Error())
			return "",err
		}

		return data["access_token"].(string),err
	}
	return "",err

}

/**
获取三方平台component_access_token
 */
func (this *GetToken) GetComponentAccessToken(appId,appSecret,ticket string)(string,error){
	tools.LogInfo("---------------GetComponentAccessToken----------------")
	url := "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
	postData := make(map[string]interface{})
	postData["component_appid"] = appId
	postData["component_appsecret"] = appSecret
	postData["component_verify_ticket"] = ticket
	postDataJson,err := json.Marshal(postData)
	if err != nil {
		tools.LogError("Json Marshal Error:",err.Error())
	}
	res,err := tools.HttpPost(url,string(postDataJson))
	if err != nil {
		tools.LogError("GetComponentAccessToken Error:" ,err.Error())
		return "",err
	}
	return res,err
}

/**
获取三方平台预授权码pre_auth_code
 */
func (this *GetToken) GetPreAuthCode(componentAccessToken,componentAppId,componentAppSecret,componentVerifyTicket string)(string,error){
	tools.LogInfo("---------------GetPreAuthCode----------------")
	url := "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=" + componentAccessToken
	postData := make(map[string]interface{})
	postData["component_appid"] = componentAppId
	postDataJson,err := json.Marshal(postData)
	if err != nil {
		tools.LogError("Json Marshal Error:",err.Error())
	}
	res,err := tools.HttpPost(url,string(postDataJson))
	if err != nil {
		tools.LogError("GetPreAuthCode Error:" ,err.Error())
		return "",err
	}

	//请求数据转换成map
	resObj,err := tools.Obj2mapObj(res)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return "",err
	}

	if resObj["errcode"] != nil && (resObj["errcode"].(float64) == this.ACCESS_TOKEN_INVALID || resObj["errcode"].(float64) == this.ACCESS_TOKEN_EXPIRES_HINT || resObj["errcode"].(float64) == this.COMPONENTA_ACCESS_TOKEN_INVALID) && componentAppId != "" && componentAppSecret != "" {
		componentAccessToken,err = this.updateComponentAccessToken(componentAppId,componentAppSecret,componentVerifyTicket)
		if err != nil {
			tools.LogError("updateComponentAccessToken Error:",err.Error())
			return "",err
		}
		url = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=" + componentAccessToken
		res,err = tools.HttpPost(url,string(postDataJson))
		if err != nil {
			tools.LogError("GetApiTicket Error:" ,err.Error())
			return "",err
		}
	}
	return res,err
}

/**
更新第三方component_access_token
 */
func (this *GetToken) updateComponentAccessToken(componentAppId,componentAppSecret,componentVerifyTicket string)(string,error) {
	componentAccessToken,err := this.GetComponentAccessToken(componentAppId,componentAppSecret,componentVerifyTicket)
	if err != nil {
		tools.LogError("updateComponentAccessToken Error:",err.Error())
		return "",err
	}

	componentAccessTokenObj,err := tools.Obj2mapObj(componentAccessToken)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return "",err
	}

	if componentAccessTokenObj["component_access_token"] != nil {
		data := make(map[string]interface{})
		data["component_access_token"] = componentAccessTokenObj["component_access_token"]
		data["component_access_token_expires_in"] = componentAccessTokenObj["expires_in"]
		data["component_access_token_expires_time"] = time.Now().Format("2006-01-02 15:04:05")
		_,err = this.Model.Table("wechat_component").Where("component_appid","=",componentAppId).Save(data)
		if err != nil {
			tools.LogError("updateComponentAccessToken Error:",err.Error())
			return "",err
		}

		return data["component_access_token"].(string),err
	}

	return "",err

}

