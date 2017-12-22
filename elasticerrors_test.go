package elasticerrors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddFind(t *testing.T) {
	ast := assert.New(t)
	var esError = &ElasticError{
		code:           ErrorCode(1),
		message:        "",
		httpStatusCode: http.StatusOK,
	}
	group := NewErrorGroup()
	group.Add(esError)

	err := group.FindError(esError.code)
	ast.NotNil(err)
	ast.Equal(esError, err)
}

func Test_CopyOnce(t *testing.T) {
	ast := assert.New(t)
	var esError = &ElasticError{
		code:           ErrorCode(1),
		message:        "",
		httpStatusCode: http.StatusOK,
	}
	group := NewErrorGroup()
	group.Add(esError)

	err := group.FindError(esError.code)
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
		code:           ErrorCode(1),
		message:        "",
		httpStatusCode: http.StatusOK,
	}
	group := NewErrorGroup()
	group.Add(esError)

	var f assert.PanicTestFunc = func() {
		group.FindError(ErrorCode(2))
		return
	}
	ast.Panics(f)
}

func Test_CodeConflict(t *testing.T) {
	ast := assert.New(t)
	var esError = &ElasticError{
		code:           ErrorCode(1),
		message:        "",
		httpStatusCode: http.StatusOK,
	}
	group := NewErrorGroup()
	group.Add(esError)
	var f assert.PanicTestFunc = func() {
		group.Add(esError)
		return
	}
	ast.Panics(f)
}
