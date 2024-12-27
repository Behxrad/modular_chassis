package utils

import (
	"errors"
	"modular_chassis/echo/pkg/errs"
	"modular_chassis/echo/pkg/services"
	"modular_chassis/echo/pkg/utils/dictionary"
)

func ConvertErrToBaseResponse(err error) services.BaseResp {
	var serviceError *errs.ServiceError
	switch {
	case errors.As(err, &serviceError):
		return services.BaseResp{
			Code:    int32(serviceError.Code),
			Message: dictionary.GetCodeTranslator().TranslateResponseCode(dictionary.Farsi, int(serviceError.Code)),
		}
	default:
		return services.BaseResp{
			Code:    int32(errs.InternalError),
			Message: dictionary.GetCodeTranslator().TranslateResponseCode(dictionary.Farsi, int(errs.InternalError)),
		}
	}
}
