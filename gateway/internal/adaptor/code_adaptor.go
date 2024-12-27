package adaptor

import (
	"modular_chassis/echo/pkg/errs"
	"net/http"
	"sync"
)

type codeAdaptor struct {
}

var (
	codeAdaptorOnce sync.Once
	codeAdaptorIns  *codeAdaptor
)

func GetCodeAdaptor() *codeAdaptor {
	codeAdaptorOnce.Do(func() {
		if codeAdaptorIns == nil {
			codeAdaptorIns = &codeAdaptor{}
		}
	})
	return codeAdaptorIns
}

func (codeAdaptor) convertCodeToHttpStatus(code errs.ResponseCode) int {
	switch {
	case code >= errs.GeneralError && code < errs.InternalError:
		return http.StatusServiceUnavailable
	case code >= errs.InternalError && code < errs.UnknownError:
		return http.StatusInternalServerError
	case code >= errs.UnknownError && code < errs.BadRequest:
		return 520
	case code >= errs.BadRequest:
		return http.StatusBadRequest
	}
	return 500
}
