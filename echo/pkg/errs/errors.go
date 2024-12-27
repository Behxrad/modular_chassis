package errs

import (
	"fmt"
	"runtime"
)

type ResponseCode int32

type ServiceError struct {
	Code         ResponseCode
	Cause        error
	TraceMessage string
	Origin       string
}

var (
	GeneralError  ResponseCode = 10000
	InternalError ResponseCode = 20000
	UnknownError  ResponseCode = 30000
	BadRequest    ResponseCode = 40000
)

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%d: \n%s", e.Code, e.TraceMessage)
}

func NewServiceErrorCode(code ResponseCode) *ServiceError {
	caller, _, _, _ := runtime.Caller(1)
	f, l := runtime.FuncForPC(caller).FileLine(caller)
	return &ServiceError{Code: code, Origin: fmt.Sprintf("at %s:line %d", f, l)}
}

func NewServiceErrorCodeAndCause(code ResponseCode, err error) *ServiceError {
	caller, _, _, _ := runtime.Caller(1)
	f, l := runtime.FuncForPC(caller).FileLine(caller)
	return &ServiceError{Code: code, Origin: fmt.Sprintf("at %s:line %d", f, l), Cause: err}
}
