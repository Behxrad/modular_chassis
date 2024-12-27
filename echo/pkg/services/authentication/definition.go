package authentication

import (
	"context"
	"modular_chassis/echo/pkg/services"
)

type UserRequest struct {
	services.BaseReq
	ID int `json:"id" binding:"required,min=5,max=20"`
}

type UserResponse struct {
	services.BaseResp
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"gte=18"`
}

type TestRequest struct {
	services.BaseReq
	Input string `json:"input" binding:"required"`
}

type TestResponse struct {
	services.BaseResp
	Output string `json:"output" binding:"required"`
}

type (
	Service interface {
		Test(ctx context.Context, request TestRequest) (TestResponse, error)
		FetchUserID(ctx context.Context, request UserRequest) (UserResponse, error)
	}
	Service2 interface {
		Test(ctx context.Context, request TestRequest) (TestResponse, error)
	}
)
