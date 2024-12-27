package api

import (
	"context"
	"modular_chassis/mediator/internal/service"
	"reflect"
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

func (m *mediatorAPI) Route(ctx context.Context, serviceType, method string, request any) (any, error) {
	response, err := service.GetRouter().HandleRequest(ctx, serviceType, method, request)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (m *mediatorAPI) GetServiceModels(serviceType, method string) (request, response reflect.Type) {
	return service.GetRegistry().GetServiceModels(serviceType, method)
}

func (m *mediatorAPI) RegisterServiceFunc(serviceImpl interface{}) {
	service.GetRegistry().RegisterService(serviceImpl)
}
