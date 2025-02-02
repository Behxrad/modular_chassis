package api

import (
	"context"
	"modular_chassis/echo/pkg/services/telecom"
	"modular_chassis/mediator/pkg/api"
	"modular_chassis/telecom/internal/service"
	"sync"
)

type TelecomService struct{}

var (
	once     sync.Once
	implsIns *TelecomService
)

func init() {
	api.GetMediatorAPI().RegisterServiceFunc(GetTelecomService())
}

func GetTelecomService() *TelecomService {
	once.Do(func() {
		if implsIns == nil {
			implsIns = &TelecomService{}
		}
	})
	return implsIns
}

func (TelecomService) GetPackages(ctx context.Context, request telecom.GetPackagesRequest) (telecom.GetPackagesResponse, error) {
	return service.GetBaseService().GetPackages(ctx, request)
}

func (TelecomService) Test(ctx context.Context, request telecom.TestRequest) (telecom.TestResponse, error) {
	return service.GetBaseService().Test(ctx, request)
}
