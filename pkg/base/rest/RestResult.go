package rest

// RestResult 通用API响应结构
// 用于所有云函数的HTTP响应输出
type RestResult struct {
	Code    int         `json:"code"`    // 状态码：1=成功，0=失败
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// NewSuccess 创建成功响应
func NewSuccess(data interface{}) *RestResult {
	return &RestResult{
		Code:    1,
		Message: "OK",
		Data:    data,
	}
}

// NewError 创建错误响应
func NewError(message string) *RestResult {
	return &RestResult{
		Code:    0,
		Message: message,
		Data:    nil,
	}
}
