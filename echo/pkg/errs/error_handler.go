package errs

import (
	"errors"
	"modular_chassis/echo/pkg/utils/dictionary"
	"sync"
)

type errorHandler struct {
}

var (
	errorHandlerOnce sync.Once
	errorHandlerIns  *errorHandler
)

func GetErrorHandler() *errorHandler {
	errorHandlerOnce.Do(func() {
		if errorHandlerIns == nil {
			errorHandlerIns = &errorHandler{}
		}
	})
	return errorHandlerIns
}

func (errorHandler) ConvertErrToBaseResponse(err error) (int32, string) {
	var serviceError *ServiceError
	switch {
	case errors.As(err, &serviceError):
		return int32(serviceError.Code), dictionary.GetCodeTranslator().TranslateResponseCode(dictionary.Farsi, int(serviceError.Code))
	default:
		return int32(InternalError), dictionary.GetCodeTranslator().TranslateResponseCode(dictionary.Farsi, int(InternalError))
	}
}
