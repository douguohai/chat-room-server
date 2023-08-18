package base

import (
	"encoding/json"
)

type Result struct {
	Code    interface{} `json:"errCode"`
	Message interface{} `json:"Message"`
	Success bool        `json:"success"`
}

func (r Result) ToJSONStr() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func BuildFail() Result {
	return Result{
		Code:    500,
		Message: "服务器异常",
		Success: false,
	}
}

func BuildFailWithMsg(code int, msg string) Result {
	return Result{
		Code:    code,
		Message: msg,
		Success: false,
	}
}

// CurrentUser 用户基本信息
type CurrentUser struct {
	Token    string `json:"token"`
	UserId   int    `json:"userId"`
	Username string `json:"username"`
}

// MarshalBinary 序列化
func (m CurrentUser) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// UnmarshalBinary 反序列化
func (m CurrentUser) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
