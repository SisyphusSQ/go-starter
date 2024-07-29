package vars

var MsgFlags = map[int]string{
	SUCCESS:               "success",
	UpdatePasswordSuccess: "修改密码成功",
	NotExistInentifier:    "该第三方账号未绑定",
	InternalERROR:         "failed",
	InvalidParams:         "请求参数错误",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[InternalERROR]
}
