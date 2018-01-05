package response


type BaseJsonBean struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status 	bool		`json:"status"`
}

func NewBaseJsonBean() *BaseJsonBean {
	return &BaseJsonBean{}
}

func (this *BaseJsonBean) GetCodeMessage(code int) string {
	var message string
	switch code{
	case 200: message = "success"
	case 400: message = "未知错误"
	case 401: message = "无此权限"
	case 500: message = "服务器异常"
	case 1001: message = "远程操作错误"
	case 1002: message = "缺少参数"
	case 1003: message = "签名错误"
	case 1004: message = "方法不存在"
	case 1005: message = "其他错误"
	case 1006: message = "参数验证错误"
	case 1007: message = "mysql错误"
	default:
		message = "未知code"

	}
	return message
}


