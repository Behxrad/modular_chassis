package service

import (
	"context"
	"modular_chassis/echo/pkg/services/finance"
	"modular_chassis/echo/pkg/services/telecom"
	"modular_chassis/mediator/pkg/api"
)

func (Base) GetPackages(ctx context.Context, request telecom.GetPackagesRequest) (telecom.GetPackagesResponse, error) {
	response, err := api.SimpleRoute[finance.GetBalanceRequest, finance.GetBalanceResponse](ctx, "finance", "GetBalance", finance.GetBalanceRequest{})
	if err != nil {
		return telecom.GetPackagesResponse{}, err
	}
	if response.Amount <= 5_000_000 {
		return telecom.GetPackagesResponse{Packages: []telecom.Package{{Name: "10GB"}}}, nil
	} else {
		return telecom.GetPackagesResponse{Packages: []telecom.Package{{Name: "200GB"}}}, nil
	}
}
