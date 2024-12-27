package api

import (
	"context"
	"encoding/json"
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

func (m *mediatorAPI) JSONRoute(ctx context.Context, serviceType, method string, raw string) (string, error) {
	reqModel, _ := service.GetRegistry().GetServiceModels(serviceType, method)

	request := reflect.New(reqModel).Interface()
	err := json.Unmarshal([]byte(raw), request)
	if err != nil {
		return "", err
	}

	response, err := service.GetRouter().HandleRequest(ctx, serviceType, method, request)
	if err != nil {
		return "", err
	}

	marshal, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func (m *mediatorAPI) RegisterServiceFunc(serviceImpl interface{}) {
	service.GetRegistry().RegisterService(serviceImpl)
}
