package service

import (
	"context"
	"fmt"
	"modular_chassis/echo/pkg/services/finance"
	"modular_chassis/echo/pkg/services/telecom"
	"modular_chassis/mediator/pkg/api"
)

func (Base) Test(ctx context.Context, request finance.TestRequest) (finance.TestResponse, error) {
	response, err := api.SimpleRoute[telecom.TestRequest, telecom.TestResponse](ctx, "telecom", "Test", telecom.TestRequest{Value: request.Value})
	if err != nil {
		return finance.TestResponse{}, err
	}
	fmt.Println(response)
	return finance.TestResponse{Value: response.Value}, nil
}
