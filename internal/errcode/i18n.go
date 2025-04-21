package errcode

var DefaultLang = "zh"

var i18n = map[string]map[ErrCode]string{
	"zh": {
		ServerError:       "内部服务器错误",
		InvalidInput:      "输入不合法",
		Unauthorized:      "未授权",
		AuthorizedExpires: "授权过期",
		PermissionDenied:  "权限不足",
		NotFound:          "资源不存在",
		Unavailable:       "不可用",
		UnknownClient:     "非法设备",
		InvalidAccount:    "无效账号",
	},
}
