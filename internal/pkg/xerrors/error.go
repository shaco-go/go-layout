package xerrors

import "net/http"

// BizError 业务错误
type BizError struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	SubMsg string `json:"sub_msg"` // 具体的业务错误信息
}

func (e *BizError) Error() string {
	return e.Msg
}

// NewBiz 业务错误
func NewBiz(msg string, code ...int) *BizError {
	err := &BizError{
		Code: http.StatusBadRequest,
		Msg:  msg,
	}
	if len(code) > 0 {
		err.Code = code[0]
	}
	return err
}

func NewBizWithSub(msg, subMsg string, code ...int) *BizError {
	err := &BizError{
		Code:   http.StatusBadRequest,
		Msg:    msg,
		SubMsg: subMsg,
	}
	if len(code) > 0 {
		err.Code = code[0]
	}
	return err
}
