package api

import (
	"context"
	"modular_chassis/echo/pkg/utils/utils"
)

func SimpleRoute[T, R any](ctx context.Context, serviceType, method string, request T) (R, error) {
	response, err := GetMediatorAPI().Route(ctx, utils.ToSnakeCase(serviceType), utils.ToSnakeCase(method), &request)
	r := response.(R)
	if err != nil {
		return r, err
	}
	return r, nil
}
