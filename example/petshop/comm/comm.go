package comm

// HttpCode http 接口错误代码
type HttpCode struct {
	ErrCode int    `json:"errcode"`          // 错误代码
	ErrMsg  string `json:"errmsg,omitempty"` // 错误说明
}
