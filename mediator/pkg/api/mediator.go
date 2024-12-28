package api

import (
	"context"
	"modular_chassis/echo/pkg/services"
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

func (m *mediatorAPI) GetServiceRequestModel(serviceType, method string) (any, error) {
	requestModel, err := service.GetRegistry().GetServiceRequestModel(serviceType, method)
	if err != nil {
		return nil, err
	}
	return reflect.New(requestModel).Interface(), nil
}

func (m *mediatorAPI) RegisterServiceFunc(serviceImpl interface{}) {
	service.GetRegistry().RegisterService(serviceImpl)
}

func (m *mediatorAPI) GetBaseReqFromModel(request any) *services.BaseReq {
	if reflect.ValueOf(request).Elem().FieldByName("BaseReq").IsValid() {
		v := reflect.ValueOf(request).Elem().FieldByName("BaseReq").Addr().Interface().(*services.BaseReq)
		return v
	}
	return nil
}
