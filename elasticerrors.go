package elasticerrors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码
type ErrorCode int

// ElasticError 弹性错误
type ElasticError struct {
	Code           ErrorCode
	Message        string
	HTTPStatusCode int

	// 是否为拷贝
	isCopy bool
}

// NewElasticError 新建错误
func NewElasticError(code ErrorCode, message string, httpStatusCode ...int) *ElasticError {
	var httpCode int
	if len(httpStatusCode) == 0 {
		httpCode = http.StatusOK
	} else {
		httpCode = httpStatusCode[0]
	}
	return &ElasticError{
		Code:           code,
		Message:        message,
		HTTPStatusCode: httpCode,
	}
}

func (e *ElasticError) copy() *ElasticError {
	if e.isCopy {
		return e
	}
	n := &ElasticError{
		Code:           e.Code,
		Message:        e.Message,
		HTTPStatusCode: e.HTTPStatusCode,
		isCopy:         true,
	}
	return n
}

// SetMessage 修改错误描述
func (e *ElasticError) SetMessage(msg string) *ElasticError {
	e = e.copy()
	e.Message = msg
	return e
}

// AppendMessage 添加错误描述
func (e *ElasticError) AppendMessage(msg string) *ElasticError {
	e = e.copy()
	e.Message = fmt.Sprintf("%s %s", e.Message, msg)
	return e
}

// SetCode 修改错误码
func (e *ElasticError) SetCode(code ErrorCode) *ElasticError {
	e = e.copy()
	e.Code = code
	return e
}

// SetHTTPStatusCode 修改HTTP响应码
func (e *ElasticError) SetHTTPStatusCode(code int) *ElasticError {
	e = e.copy()
	e.HTTPStatusCode = code
	return e
}

// ErrorGroup 错误组
type ErrorGroup struct {
	errorMap map[ErrorCode]*ElasticError
}

// NewErrorGroup 新建错误组
func NewErrorGroup() *ErrorGroup {
	group := &ErrorGroup{}
	group.errorMap = make(map[ErrorCode]*ElasticError)
	return group
}

// Add 添加错误
func (g *ErrorGroup) Add(esError *ElasticError) {
	if tmpError, exist := g.errorMap[esError.Code]; exist {
		var msg = fmt.Sprintf("ErrorCode %d has been added to the group, details: %+v", int(esError.Code), tmpError)
		panic(msg)
	}
	if esError.HTTPStatusCode == 0 {
		esError.HTTPStatusCode = http.StatusOK
	}
	g.errorMap[esError.Code] = esError
	return
}

// FindError 查找错误
func (g *ErrorGroup) FindError(code ErrorCode) *ElasticError {
	esError, ok := g.errorMap[code]
	if !ok {
		panic(fmt.Sprintf("错误码未定义，ErrorCode:%d", code))
	}
	return esError
}

var defaultGroup = NewErrorGroup()

// Add 添加错误
func Add(esError *ElasticError) {
	defaultGroup.Add(esError)
	return
}

// FindError 查找错误
func FindError(code ErrorCode) *ElasticError {
	return defaultGroup.FindError(code)
}
