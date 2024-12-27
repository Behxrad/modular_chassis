package api

import (
	"context"
	"modular_chassis/echo/pkg/errs"
	"modular_chassis/echo/pkg/services/authentication"
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
	mediatorAPIs.GetMediatorAPI().RegisterServiceFunc(getService())
}

func (i *impls) Test(ctx context.Context, request authentication.TestRequest) (authentication.TestResponse, error) {
	return authentication.TestResponse{}, errs.NewServiceErrorCode(authentication.UserNotFound)
}
func (i *impls) FetchUserID(ctx context.Context, request authentication.UserRequest) (authentication.UserResponse, error) {
	return authentication.UserResponse{}, nil
}
