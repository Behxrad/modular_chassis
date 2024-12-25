package api

import (
	"context"
	"modular_chassis/echo/pkg/services"
	"modular_chassis/mediator/internal/service"
)

func Route[P, R any](ctx context.Context, serviceType, method string, request services.ServiceRequest[P]) (services.ServiceResponse[R], error) {
	response, err := service.HandleRequest[P, R](ctx, serviceType, method, request)
	if err != nil {
		return services.ServiceResponse[R]{}, err
	}
	return response, nil
}

func RegisterServiceFunc(serviceImpl interface{}) {
	service.GetRegistry().RegisterService(serviceImpl)
}
