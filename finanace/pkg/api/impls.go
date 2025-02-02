package api

import (
	"context"
	"modular_chassis/echo/pkg/services/finance"
	"modular_chassis/finanace/internal/service"
	"modular_chassis/mediator/pkg/api"
	"sync"
)

type FinanceService struct{}

var (
	once     sync.Once
	implsIns *FinanceService
)

func init() {
	api.GetMediatorAPI().RegisterServiceFunc(GetBalanceService())
}

func GetBalanceService() *FinanceService {
	once.Do(func() {
		if implsIns == nil {
			implsIns = &FinanceService{}
		}
	})
	return implsIns
}

func (FinanceService) GetBalance(ctx context.Context, request finance.GetBalanceRequest) (finance.GetBalanceResponse, error) {
	return service.GetBaseService().GetBalance(ctx, request)
}

func (FinanceService) Test(ctx context.Context, request finance.TestRequest) (finance.TestResponse, error) {
	return service.GetBaseService().Test(ctx, request)
}
