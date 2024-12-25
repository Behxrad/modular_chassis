package carpay

import (
	"context"
	"modular_chassis/echo/pkg/services"
)

type ListVehiclesRequest struct {
	Mobile string `json:"username" binding:"required,min=5,max=20" swaggo:"required,min=5,max=20"`
}

type ListVehiclesResponse struct {
	Vehicles []Vehicle `json:"vehicles" swaggo:"required"`
}

type Vehicle struct {
	PlateNumber string `json:"plateNumber" binding:"required"`
}

type TestRequest struct {
	Input string `json:"input" binding:"required"`
}

type TestResponse struct {
	Output string `json:"output" binding:"required"`
}

type (
	Service3 interface {
		Test(ctx context.Context, request services.ServiceRequest[TestRequest]) (services.ServiceResponse[TestResponse], error)
		BuyVehicle(ctx context.Context, request services.ServiceRequest[TestRequest]) (services.ServiceResponse[TestResponse], error)
	}
	Service4 interface {
		ListVehicles(ctx context.Context, req services.ServiceRequest[ListVehiclesRequest]) (services.ServiceResponse[ListVehiclesResponse], error)
		Test(ctx context.Context, request services.ServiceRequest[TestRequest]) (services.ServiceResponse[TestResponse], error)
	}
)
