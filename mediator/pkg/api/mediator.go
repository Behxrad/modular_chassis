package api

import (
	"context"
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

func (m *mediatorAPI) Route(ctx context.Context, serviceType, method string, request string) (string, error) {
	response, err := service.GetRouter().HandleRequest(ctx, serviceType, method, request)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (m *mediatorAPI) RegisterServiceFunc(serviceImpl interface{}) {
	service.GetRegistry().RegisterService(serviceImpl)
}
