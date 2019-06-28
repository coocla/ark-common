package constants

const (
	// EN 英语
	EN string = "en"
	// CN 简体中文
	CN string = "zh-CN"
)

// AllLang 支持的所有语言
var AllLang = []string{
	EN,
	CN,
}

const (
	// Success 成功
	Success = 0
	// ServerError 服务器内部错误
	ServerError = 500001
	// InvalidParam 错误的请求参数
	InvalidParam = 400001
	// NotSupportCloudAction 不支持的云操作
	NotSupportCloudAction = 400002
	// InvalidResourceID 非法的资源ID
	InvalidResourceID = 400003
	// InvalidCloudAccountID 非法的云商账户ID
	InvalidCloudAccountID = 400004
)

// CodeMessage code和文本对应关系
var (
	CodeMessage = map[int]map[string]string{
		Success: {
			EN: "success",
			CN: "成功",
		},
		ServerError: {
			EN: "server internal error",
			CN: "服务器内部错误",
		},
		InvalidParam: {
			EN: "invalid params",
			CN: "请求参数不合法",
		},
		NotSupportCloudAction: {
			EN: "not support the action",
			CN: "不支持此操作",
		},
		InvalidResourceID: {
			EN: "invalid resource id",
			CN: "非法的资源ID",
		},
		InvalidCloudAccountID: {
			EN: "invalid cloudaccount id",
			CN: "非法的云商账户ID",
		},
	}
)
