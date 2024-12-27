package adaptor

import (
	"modular_chassis/echo/pkg/errs"
	"modular_chassis/echo/pkg/services"
	"modular_chassis/echo/pkg/utils/utils"
	"sync"
)

type errHTTPAdaptor struct {
}

var (
	errHTTPAdaptorOnce sync.Once
	errHTTPAdaptorIns  *errHTTPAdaptor
)

func GetErrHTTPAdaptor() *errHTTPAdaptor {
	errHTTPAdaptorOnce.Do(func() {
		if errHTTPAdaptorIns == nil {
			errHTTPAdaptorIns = &errHTTPAdaptor{}
		}
	})
	return errHTTPAdaptorIns
}

func (errHTTPAdaptor) MakeErroneousResponse(err error) (services.BaseResp, int) {
	response := utils.ConvertErrToBaseResponse(err)
	return response, GetCodeAdaptor().convertCodeToHttpStatus(errs.ResponseCode(response.Code))
}
