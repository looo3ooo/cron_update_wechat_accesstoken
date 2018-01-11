package crontab

import (
	"updatetoken/tools"
	"fmt"
	"time"
	"strconv"
	"reflect"
	"goini"
	"updatetoken/mysql"
)

type AutoUpdateToken struct {
	SIGNKEY string
	COMPONENT_ATTRIBUTE		[]map[string]string
	getToken *GetToken
	Model *gomysql.SqlModel
}

func (this *AutoUpdateToken) AttrInit() *AutoUpdateToken{
	ConfigCentor := goini.SetConfig("./config/config.ini")
	signkey := ConfigCentor.GetValue("autoupdatetoken", "signkey")
	component_appid1 := ConfigCentor.GetValue("autoupdatetoken", "component_appid1")
	component_app_secret1 := ConfigCentor.GetValue("autoupdatetoken", "component_app_secret1")

	this.SIGNKEY = signkey

	COMPONENT_ATTRIBUTE1 := make(map[string]string)
	COMPONENT_ATTRIBUTE1["component_appid"] = component_appid1
	COMPONENT_ATTRIBUTE1["component_app_secret"] = component_app_secret1

	this.COMPONENT_ATTRIBUTE = append(this.COMPONENT_ATTRIBUTE,COMPONENT_ATTRIBUTE1)

	this.Model = pool.NewModel()

	this.getToken = new(GetToken)
	this.getToken = this.getToken.AttrInit()

	 return this
}

func (this *AutoUpdateToken) Index(){

	wechat,err := this.Model.Table("wechat").FindAll()
	if err != nil {
		tools.LogError("mysql error:",err.Error())
	}
	fmt.Println("定时任务:updatetoken")
	tools.LogInfo("获取所有商户微信配置信息：",wechat)

	var wechatAttribute map[string]interface{}
	for _,v := range wechat {
		wechatAttribute,err = this.Model.Table("wechat_attribute").Where("psm_id","=",v["psm_id"]).FindOne()
		tools.LogInfo(wechatAttribute)

		//过期或者为空获取新数据，否则将原有数据填入
		timeStamp := time.Now().Unix()

		//access_token
		if len(wechatAttribute) != 0  && this.timeToUnix(this.toString(v["access_token_expires_time"])) + this.toInt64(v["access_token_expires_in"]) - 30*60 < timeStamp {
			if v["appid"] != nil && v["appid"].(string) != ""{
				newAccessToken,err := this.updateAccessToken(v,wechatAttribute)
				if err != nil {
					tools.LogError("updateAccessToken Error:",err.Error())
				} else {
					v["access_token"] = newAccessToken
				}
			}
		}

		// jsapi_ticket
		if len(wechatAttribute) != 0  && this.timeToUnix(this.toString(v["jsapi_ticket_expires_time"])) + this.toInt64(v["jsapi_ticket_expires_in"]) - 30*60 < timeStamp {
			if v["appid"] != nil && v["appid"].(string) != "" {
				this.updateJsapiTicket(v,wechatAttribute)
			}
		}


		// api_ticket
		if len(wechatAttribute) != 0  && this.timeToUnix(this.toString(v["api_ticket_expires_time"])) + this.toInt64(v["api_ticket_expires_in"]) - 30*60 < timeStamp {
			if v["appid"] != nil && v["appid"].(string) != "" {
				this.updateApiTicket(v,wechatAttribute)
			}
		}

		// app_access_token
		if len(wechatAttribute) != 0  && this.timeToUnix(this.toString(v["app_access_token_expires_time"])) + this.toInt64(v["app_access_token_expires_in"]) - 30*60 < timeStamp {
			if v["app_appid"] != nil && v["app_appid"].(string) != "" {
				this.updateAppAccessToken(v,wechatAttribute)
			}
		}
	}

	for _,ATTR := range this.COMPONENT_ATTRIBUTE {
		wechatComponent,err := this.Model.Table("wechat_component").Where("component_appid","=",ATTR["component_appid"]).FindOne()
		if err != nil {
			tools.LogInfo("Update component error:",err.Error())
		}

		timeStamp := time.Now().Unix()
		//component_access_token
		if len(wechatComponent) != 0  && this.timeToUnix(this.toString(wechatComponent["component_access_token_expires_time"])) + this.toInt64(wechatComponent["component_access_token_expires_in"]) - 30*60 < timeStamp {
			this.updateComponentAccessToken(ATTR["component_appid"],ATTR["component_app_secret"],wechatComponent["component_verify_ticket"])
		}

		//component_pre_auth_code
		if len(wechatComponent) != 0 && this.timeToUnix(this.toString(wechatComponent["component_pre_auth_code_expires_time"])) + this.toInt64(wechatComponent["component_pre_auth_code_expires_in"]) - 5*60 < timeStamp {
			this.updateComponentPreAuthCode(wechatComponent["component_access_token"],ATTR["component_appid"],ATTR["component_app_secret"],wechatComponent["component_verify_ticket"])
		}
	}

	tools.LogInfo("--------执行完毕---------")

}

//更新access_token
func (this *AutoUpdateToken) updateAccessToken(wechat map[string]interface{},wechatAttribute map[string]interface{}) (string,error){
	tools.LogInfo("超时重新获取accesstoken")
	res,err := this.getToken.GetAccessToken(this.toString(wechat["appid"]),this.toString(wechatAttribute["secret"]))
	if err != nil {
		tools.LogError("GetAccessToken Error:",err.Error())
		return "",err
	}

	//请求数据转换成map
	resObj,err := tools.Obj2mapObj(res)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return "",err
	}

	if resObj["errcode"] != nil && resObj["errcode"].(float64) != 0 {
		tools.LogError("accesstoken远程调用失败:",resObj["errmsg"],":",resObj["errcode"])
		return "",err
	}

	//更新数据库
	data := make(map[string]interface{})
	data["access_token"] = resObj["access_token"].(string)
	data["access_token_expires_in"] = resObj["expires_in"].(float64)
	data["access_token_expires_time"] = time.Now().Format("2006-01-02 15:04:05")
	num,err := this.Model.Table("wechat").Where("appid","=",wechat["appid"]).Save(data)

	if err != nil {
		tools.LogError("更新accesstoken失败：",err.Error())
		return "",err
	}

	tools.LogInfo("更新",num,"条数据")
	return data["access_token"].(string),err
}

/**
更新app_access_token
 */
func (this *AutoUpdateToken) updateAppAccessToken(wechat map[string]interface{},wechatAttribute map[string]interface{}){
	tools.LogInfo("超时重新获取app_access_token")
	res,err := this.getToken.GetAppAccessToken(this.toString(wechat["app_appid"]),this.toString(wechatAttribute["app_secret"]))
	if err != nil {
		tools.LogError("GetAppAccessToken Error:",err.Error())
		return
	}

	//请求数据转换成map
	resObj,err := tools.Obj2mapObj(res)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return
	}

	if resObj["errcode"] != nil && resObj["errcode"].(float64) != 0 {
		tools.LogError("app_accesstoken远程调用失败:",resObj["errmsg"],":",resObj["errcode"])
		return
	}

	//更新数据库
	data := make(map[string]interface{})
	data["app_access_token"] = resObj["access_token"].(string)
	data["app_access_token_expires_in"] = resObj["expires_in"].(float64)
	data["app_access_token_expires_time"] = time.Now().Format("2006-01-02 15:04:05")
	num,err := this.Model.Table("wechat").Where("app_appid","=",wechat["app_appid"]).Save(data)

	if err != nil {
		tools.LogError("更新app_accesstoken失败：",err.Error())
		return
	}

	tools.LogInfo("更新",num,"条数据")
}

// 更新jsapi_ticket
func (this *AutoUpdateToken) updateJsapiTicket(wechat map[string]interface{},wechatAttribute map[string]interface{}){
	tools.LogInfo("超时重新获取jsapi_ticket")

	res,err := this.getToken.GetJsApiTicket(this.toString(wechat["access_token"]),this.toString(wechat["appid"]),this.toString(wechatAttribute["secret"]))
	if err != nil {
		tools.LogError("GetJsApiTicket Error:",err.Error())
		return
	}

	//请求数据转换成map
	resObj,err := tools.Obj2mapObj(res)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return
	}

	if resObj["errcode"] != nil && resObj["errcode"].(float64) != 0 {
		tools.LogError("jsapi_ticket远程调用失败:",resObj["errmsg"],":",resObj["errcode"])
		return
	}

	//更新数据库
	data := make(map[string]interface{})
	data["jsapi_ticket"] = resObj["ticket"].(string)
	data["jsapi_ticket_expires_in"] = resObj["expires_in"].(float64)
	data["jsapi_ticket_expires_time"] = time.Now().Format("2006-01-02 15:04:05")
	num,err := this.Model.Table("wechat").Where("appid","=",wechat["appid"]).Save(data)

	if err != nil {
		tools.LogError("更新jsapi_ticket失败：",err.Error())
		return
	}
	tools.LogInfo("更新",num,"条数据")
}

// 更新api_ticket
func (this *AutoUpdateToken) updateApiTicket(wechat map[string]interface{},wechatAttribute map[string]interface{}) {
	tools.LogInfo("超时重新获取api_ticket")

	res,err := this.getToken.GetApiTicket(this.toString(wechat["access_token"]),this.toString(wechat["appid"]),this.toString(wechatAttribute["secret"]))
	if err != nil {
		tools.LogError("GetApiTicket Error:",err.Error())
		return
	}

	//请求数据转换成map
	resObj,err := tools.Obj2mapObj(res)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return
	}
	if resObj["errcode"] != nil && resObj["errcode"].(float64) != 0 {
		tools.LogError("api_ticket远程调用失败:",resObj["errmsg"],":",resObj["errcode"])
		return
	}

	//更新数据库
	data := make(map[string]interface{})
	data["api_ticket"] = resObj["ticket"].(string)
	data["api_ticket_expires_in"] = resObj["expires_in"].(float64)
	data["api_ticket_expires_time"] = time.Now().Format("2006-01-02 15:04:05")
	num,err := this.Model.Table("wechat").Where("appid","=",wechat["appid"]).Save(data)

	if err != nil {
		tools.LogError("更新api_ticket失败：",err.Error())
		return
	}
	tools.LogInfo("更新",num,"条数据")

}

// 更新三方平台component_access_token
func (this *AutoUpdateToken) updateComponentAccessToken(componentAppId,componentAppSecret,componentVerifyTicket interface{}){
	tools.LogInfo("超时重新获取component_access_token")

	res,err := this.getToken.GetComponentAccessToken(this.toString(componentAppId),this.toString(componentAppSecret),this.toString(componentVerifyTicket))
	if err != nil {
		tools.LogError("GetComponentAccessToken Error:",err.Error())
		return
	}

	//请求数据转换成map
	resObj,err := tools.Obj2mapObj(res)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return
	}
	if resObj["errcode"] != nil && resObj["errcode"].(float64) != 0 {
		tools.LogError("component_access_token远程调用失败:",resObj["errmsg"],":",resObj["errcode"])
		return
	}

	//更新数据库
	data := make(map[string]interface{})
	data["component_access_token"] = resObj["component_access_token"].(string)
	data["component_access_token_expires_in"] = resObj["expires_in"].(float64)
	data["component_access_token_expires_time"] = time.Now().Format("2006-01-02 15:04:05")
	num,err := this.Model.Table("wechat_component").Where("component_appid","=",componentAppId).Save(data)

	if err != nil {
		tools.LogError("更新component_access_token失败：",err.Error())
		return
	}
	tools.LogInfo("更新",num,"条数据")
}

// 更新component_pre_auth_code
func (this *AutoUpdateToken) updateComponentPreAuthCode(componentAccessToken,componentAppId,componentAppSecret,componentVerifyTicket interface{}){
	tools.LogInfo("超时重新获取component_pre_auth_code")

	res,err := this.getToken.GetPreAuthCode(this.toString(componentAccessToken),this.toString(componentAppId),this.toString(componentAppSecret),this.toString(componentVerifyTicket))
	if err != nil {
		tools.LogError("GetPreAuthCode Error:",err.Error())
		return
	}

	//请求数据转换成map
	resObj,err := tools.Obj2mapObj(res)
	if err != nil {
		tools.LogError("Obj2mapObj Error:",err.Error())
		return
	}

	if resObj["errcode"] != nil && resObj["errcode"].(float64) != 0 {
		tools.LogError("component_pre_auth_code远程调用失败:",resObj["errmsg"],":",resObj["errcode"])
		return
	}

	//更新数据库
	data := make(map[string]interface{})
	data["component_pre_auth_code"] = resObj["pre_auth_code"].(string)
	data["component_pre_auth_code_expires_in"] = resObj["expires_in"].(float64)
	data["component_pre_auth_code_expires_time"] = time.Now().Format("2006-01-02 15:04:05")
	num,err := this.Model.Table("wechat_component").Where("component_appid","=",componentAppId).Save(data)

	if err != nil {
		tools.LogError("更新component_pre_auth_code失败：",err.Error())
		return
	}
	tools.LogInfo("更新",num,"条数据")
}

func (this *AutoUpdateToken) toString(str interface{}) string{
	var strString string
	if str == nil {
		strString = ""
	}
	if reflect.ValueOf(str).Kind() == reflect.String {
		strString = str.(string)
	} else if reflect.ValueOf(str).Kind() == reflect.Int64 {
		strString = strconv.FormatInt(str.(int64),10)
	} else if reflect.ValueOf(str).Kind() == reflect.Int {
		strString = strconv.Itoa(str.(int))
	}
	return strString
}

func (this * AutoUpdateToken) toInt64(num interface{}) int64  {
	var numInt64 int64
	if num == nil {
		numInt64 = 0
	}
	if reflect.ValueOf(num).Kind() == reflect.String {
		numInt64,_ = strconv.ParseInt(num.(string),10,64)
	} else if reflect.ValueOf(num).Kind() == reflect.Int64 {
		numInt64 = num.(int64)
	} else if reflect.ValueOf(num).Kind() == reflect.Int {
		numInt64,_ = strconv.ParseInt(strconv.Itoa(num.(int)),10,64)
	}
	return numInt64
}


func (this *AutoUpdateToken) timeToUnix(timestamp interface{}) int64 {
	var timeUnix int64
	if reflect.ValueOf(timestamp).Kind() == reflect.String {
		timeParse,_ := time.Parse("2006-01-02 15:04:05",timestamp.(string))
		timeUnix = timeParse.Unix()
	} else if reflect.ValueOf(timestamp).Kind() == reflect.Int64 {
		timeParse,_ := time.Parse("2006-01-02 15:04:05",strconv.FormatInt(timestamp.(int64),10))
		timeUnix = timeParse.Unix()
	}
	return timeUnix
}
