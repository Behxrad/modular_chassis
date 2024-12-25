package api

import (
	"context"
	"fmt"
	"modular_chassis/echo/pkg/services"
	"modular_chassis/echo/pkg/services/carpay"
	mediatorAPIs "modular_chassis/mediator/pkg/api"
	"sync"
)

type impls struct{}

var (
	once     sync.Once
	implsIns *impls
)

func getService() *impls {
	once.Do(func() {
		if implsIns == nil {
			implsIns = &impls{}
		}
	})
	return implsIns
}

func init() {
	mediatorAPIs.RegisterServiceFunc(getService())
}

func (impls) ListVehicles(ctx context.Context, req services.ServiceRequest[carpay.ListVehiclesRequest]) (services.ServiceResponse[carpay.ListVehiclesResponse], error) {
	fmt.Println(">>> ListVehicles called")
	return services.ServiceResponse[carpay.ListVehiclesResponse]{}, nil
}

func (impls) Test(ctx context.Context, request services.ServiceRequest[carpay.TestRequest]) (services.ServiceResponse[carpay.TestResponse], error) {
	fmt.Println(">>> Test called")
	return services.ServiceResponse[carpay.TestResponse]{}, nil
}
