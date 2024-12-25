package api

import (
	"context"
	"modular_chassis/echo/pkg/services"
	"modular_chassis/mediator/internal/service"
	"sync"
)

var (
	once           sync.Once
	mediatorAPIIns *mediatorAPI
)

type mediatorAPI struct {
}

func GetMediatorAPI() *mediatorAPI {
	once.Do(func() {
		if mediatorAPIIns == nil {
			mediatorAPIIns = &mediatorAPI{}
		}
	})
	return mediatorAPIIns
}

func (m *mediatorAPI) Route(ctx context.Context, serviceType, method string, request services.ServiceRequest[string]) (services.ServiceResponse[string], error) {
	response, err := service.HandleRequest(ctx, serviceType, method, request)
	if err != nil {
		return services.ServiceResponse[string]{}, err
	}
	return response, nil
}

func (m *mediatorAPI) RegisterServiceFunc(serviceImpl interface{}) {
	service.GetRegistry().RegisterService(serviceImpl)
}
