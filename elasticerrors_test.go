package elasticerrors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddFind(t *testing.T) {
	ast := assert.New(t)
	var esError = &ElasticError{
		Code:           ErrorCode(1),
		Message:        "",
		HTTPStatusCode: http.StatusOK,
	}
	defaultGroup.Add(esError)

	err := defaultGroup.FindError(esError.Code)
	ast.NotNil(err)
	ast.Equal(esError, err)
}

func Test_CopyOnce(t *testing.T) {
	ast := assert.New(t)
	var esError = &ElasticError{
		Code:           ErrorCode(1),
		Message:        "",
		HTTPStatusCode: http.StatusOK,
	}
	defaultGroup.Add(esError)

	err := defaultGroup.FindError(esError.Code)
	err1 := err.SetCode(ErrorCode(2))
	err2 := err1.SetMessage("ss")
	err3 := err2.SetHTTPStatusCode(http.StatusAccepted)
	ast.NotEqual(err, err1)
	ast.Equal(err1, err2)
	ast.Equal(err1, err3)
}

func Test_UndefiendeCode(t *testing.T) {
	ast := assert.New(t)
	var esError = &ElasticError{
		Code:           ErrorCode(1),
		Message:        "",
		HTTPStatusCode: http.StatusOK,
	}
	defaultGroup.Add(esError)

	var f assert.PanicTestFunc = func() {
		defaultGroup.FindError(ErrorCode(2))
		return
	}
	ast.Panics(f)
}

func Test_CodeConflict(t *testing.T) {
	ast := assert.New(t)
	var esError = &ElasticError{
		Code:           ErrorCode(1),
		Message:        "",
		HTTPStatusCode: http.StatusOK,
	}
	defaultGroup.Add(esError)
	var f assert.PanicTestFunc = func() {
		defaultGroup.Add(esError)
		return
	}
	ast.Panics(f)
}
