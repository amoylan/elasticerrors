package elasticerrors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码
type ErrorCode int

// ElasticError 弹性错误
type ElasticError struct {
	code           ErrorCode
	message        string
	httpStatusCode int

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
		code:           code,
		message:        message,
		httpStatusCode: httpCode,
	}
}

func (e *ElasticError) copy() *ElasticError {
	if e.isCopy {
		return e
	}
	n := &ElasticError{
		code:           e.code,
		message:        e.message,
		httpStatusCode: e.httpStatusCode,
		isCopy:         true,
	}
	return n
}

// SetMessage 修改错误描述
func (e *ElasticError) SetMessage(msg string) *ElasticError {
	e = e.copy()
	e.message = msg
	return e
}

// AppendMessage 添加错误描述
func (e *ElasticError) AppendMessage(msg string) *ElasticError {
	e = e.copy()
	e.message = fmt.Sprintf("%s %s", e.message, msg)
	return e
}

// Error 实现 error接口
func (e *ElasticError) Error() string {
	return e.message
}

// SetCode 修改错误码
func (e *ElasticError) SetCode(code ErrorCode) *ElasticError {
	e = e.copy()
	e.code = code
	return e
}

// Code 返回错误码
func (e *ElasticError) Code() int {
	return int(e.code)
}

// SetHTTPStatusCode 修改HTTP响应码
func (e *ElasticError) SetHTTPStatusCode(code int) *ElasticError {
	e = e.copy()
	e.httpStatusCode = code
	return e
}

// HTTPStatusCode 返回HTTP错误码
func (e *ElasticError) HTTPStatusCode() int {
	return e.httpStatusCode
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
	if tmpError, exist := g.errorMap[esError.code]; exist {
		var msg = fmt.Sprintf("ErrorCode %d has been added to the group, details: %+v", int(esError.code), tmpError)
		panic(msg)
	}
	if esError.httpStatusCode == 0 {
		esError.httpStatusCode = http.StatusOK
	}
	g.errorMap[esError.code] = esError
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
