package service

import (
	"context"
	"modular_chassis/echo/pkg/services/telecom"
)

func (Base) Test(ctx context.Context, request telecom.TestRequest) (telecom.TestResponse, error) {
	return telecom.TestResponse{Value: request.Value}, nil
}
