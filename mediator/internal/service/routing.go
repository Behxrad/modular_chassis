package service

import (
	"context"
	"errors"
	"modular_chassis/echo/pkg/services"
)

var routerIns *router

type router struct{}

func getRouter() *router {
	once.Do(func() {
		if routerIns == nil {
			routerIns = &router{}
		}
	})
	return routerIns
}

type ImplFunc[P, R any] func(ctx context.Context, request services.ServiceRequest[P]) (services.ServiceResponse[R], error)

func HandleRequest[P, R any](ctx context.Context, serviceType, method string, request services.ServiceRequest[P]) (services.ServiceResponse[R], error) {
	service, exists := GetRegistry().services[serviceType+"/"+method]
	if !exists {
		return services.ServiceResponse[R]{}, errors.New("service not found")
	}
	return service.(ImplFunc[P, R])(ctx, request)
}
